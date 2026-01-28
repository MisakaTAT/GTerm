package adapter

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
)

type Serial struct {
	port     serial.Port
	ws       *websocket.Conn
	paused   bool
	pausedMu sync.Mutex
	buffer   *bytes.Buffer
	bufferMu sync.Mutex
	flushCh  chan struct{}
}

const (
	serialBufferSize = 128 * 1024 // 与前端 bytesThreshold 对齐
)

func NewSerial(ws *websocket.Conn) *Serial {
	return &Serial{
		ws:      ws,
		buffer:  &bytes.Buffer{},
		flushCh: make(chan struct{}, 1),
	}
}

func (s *Serial) Open(portName string) error {
	slog.Info("Opening serial port: %s", portName)
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	slog.Debug("Serial config: baudRate=%d, parity=%v, dataBits=%d, stopBits=%v",
		mode.BaudRate, mode.Parity, mode.DataBits, mode.StopBits)

	port, err := serial.Open(portName, mode)
	if err != nil {
		slog.Error("Failed to open serial port: %v", err)
		return err
	}
	s.port = port
	slog.Info("Serial port opened successfully: %s", portName)
	return nil
}

func (s *Serial) close() {
	slog.Info("Closing serial connection")
	if s.port != nil {
		_ = s.port.Close()
	}
}

func (s *Serial) Input(quitSignal chan bool) {
	slog.Info("Starting WebSocket input monitoring")
	defer s.setQuit(quitSignal)
	for {
		select {
		case <-quitSignal:
			slog.Debug("Received quit signal, stopping input handler")
			return
		default:
			_, data, err := s.ws.ReadMessage()
			if err != nil {
				slog.Error("Failed to read WebSocket message: %v", err)
				return
			}
			msg := &terminal.Payload{}
			_ = json.Unmarshal(data, &msg)
			switch msg.Type {
			case enums.TerminalTypeCMD:
				slog.Debug("Sending command to serial port: %s", msg.Cmd)
				if _, err = s.port.Write([]byte(msg.Cmd)); err != nil {
					slog.Error("Failed to write to serial port: %v", err)
				}
			case enums.TerminalTypeFlowControl:
				if msg.Pause != nil {
					s.pausedMu.Lock()
					s.paused = *msg.Pause
					s.pausedMu.Unlock()
					if *msg.Pause {
						slog.Debug("Flow control: paused")
					} else {
						slog.Debug("Flow control: resumed")
					}
				}
			}
		}
	}
}

func (s *Serial) flushBuffer() {
	s.pausedMu.Lock()
	paused := s.paused
	s.pausedMu.Unlock()

	if paused {
		s.bufferMu.Lock()
		s.buffer.Reset()
		s.bufferMu.Unlock()
		return
	}

	s.bufferMu.Lock()
	bufferLen := s.buffer.Len()
	if bufferLen == 0 {
		s.bufferMu.Unlock()
		return
	}

	const maxChunkSize = 64 * 1024 // 64KB per chunk
	var content string
	if bufferLen > maxChunkSize {
		data := make([]byte, maxChunkSize)
		n, _ := s.buffer.Read(data)
		content = string(data[:n])
	} else {
		content = s.buffer.String()
		s.buffer.Reset()
	}
	s.bufferMu.Unlock()

	if len(content) > 0 {
		if err := s.ws.WriteJSON(&types.Message{
			Type:    enums.TerminalTypeData,
			Content: content,
		}); err != nil {
			slog.Error("Failed to write WebSocket message: %v", err)
		}
	}
}

func (s *Serial) Output(quitSignal chan bool) {
	slog.Info("Starting serial port output reading")
	defer s.setQuit(quitSignal)
	if s.port == nil {
		slog.Error("Serial port not open")
		s.setQuit(quitSignal)
		return
	}

	go func() {
		buff := make([]byte, 4096)
		for {
			select {
			case <-quitSignal:
				return
			default:
				n, err := s.port.Read(buff)
				if err != nil {
					slog.Error("Failed to read data from serial port: %v", err)
					return
				}
				if n == 0 {
					continue
				}

				s.pausedMu.Lock()
				paused := s.paused
				s.pausedMu.Unlock()

				if paused {
					continue
				}

				s.bufferMu.Lock()
				currentLen := s.buffer.Len()
				if currentLen > serialBufferSize {
					s.buffer.Reset()
				}
				if currentLen+len(buff[:n]) > serialBufferSize*2 {
					keepLen := serialBufferSize - len(buff[:n])
					if keepLen > 0 {
						oldData := s.buffer.Bytes()
						s.buffer.Reset()
						s.buffer.Write(oldData[len(oldData)-keepLen:])
					} else {
						s.buffer.Reset()
					}
				}
				s.buffer.Write(buff[:n])
				s.bufferMu.Unlock()

				// 通知 flush（非阻塞）
				select {
				case s.flushCh <- struct{}{}:
				default:
				}
			}
		}
	}()

	fallbackTick := time.NewTicker(50 * time.Millisecond)
	defer fallbackTick.Stop()

	cleanupTick := time.NewTicker(time.Second)
	defer cleanupTick.Stop()

	for {
		select {
		case <-quitSignal:
			s.flushBuffer()
			return
		case <-s.flushCh:
			// 单次唤醒最多处理 1MB，避免长时间占用
			const maxFlushBytes = 1024 * 1024 // 1MB
			flushedBytes := 0
			for flushedBytes < maxFlushBytes {
				s.bufferMu.Lock()
				beforeLen := s.buffer.Len()
				s.bufferMu.Unlock()

				if beforeLen == 0 {
					break
				}

				s.flushBuffer()

				s.bufferMu.Lock()
				afterLen := s.buffer.Len()
				s.bufferMu.Unlock()

				flushedBytes += beforeLen - afterLen
				if afterLen == 0 {
					break
				}
			}
		case <-cleanupTick.C:
			// 暂停时释放 buffer 底层容量
			s.pausedMu.Lock()
			paused := s.paused
			s.pausedMu.Unlock()

			if paused {
				s.bufferMu.Lock()
				if s.buffer.Len() > 0 {
					s.buffer.Reset()
					s.buffer = &bytes.Buffer{}
				}
				s.bufferMu.Unlock()
			}
		case <-fallbackTick.C:
			s.flushBuffer()
		}
	}
}

func (s *Serial) Wait(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	<-quitSignal
	s.close()
}

func (s *Serial) setQuit(ch chan bool) {
	ch <- true
}
