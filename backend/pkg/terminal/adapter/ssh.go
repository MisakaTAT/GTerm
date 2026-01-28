package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/MisakaTAT/GTerm/backend/enums"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	conf      *commonssh.Config
	ws        *websocket.Conn
	session   *ssh.Session
	stdinPipe io.WriteCloser
	writer    *writer
	flushCh   chan struct{}
	paused    bool
	pausedMu  sync.Mutex
}

type writer struct {
	buffer   bytes.Buffer
	mu       sync.Mutex
	maxSize  int // buffer 上限
	flushCh  chan struct{}
	paused   bool
	pausedMu sync.Mutex
}

const (
	maxBufferSize = 128 * 1024 // 与前端 bytesThreshold 对齐
)

func (w *writer) Write(p []byte) (int, error) {
	w.pausedMu.Lock()
	paused := w.paused
	w.pausedMu.Unlock()

	// 暂停时丢弃输出，避免内存累积
	if paused {
		return len(p), nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// 限制 buffer 最大长度，超过则截断
	currentLen := w.buffer.Len()
	if currentLen > w.maxSize {
		w.buffer.Reset()
		currentLen = 0
	}

	if currentLen+len(p) > w.maxSize {
		writeLen := w.maxSize - currentLen
		if writeLen > 0 {
			w.buffer.Write(p[:writeLen])
		}
		return len(p), nil
	}

	n, err := w.buffer.Write(p)
	if err != nil {
		return n, err
	}

	// 通知输出协程 flush（非阻塞）
	if w.flushCh != nil {
		select {
		case w.flushCh <- struct{}{}:
		default:
		}
	}
	return n, nil
}

func (w *writer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}

func (w *writer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.String()
}

func (w *writer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
	// 释放底层容量
	w.buffer = bytes.Buffer{}
}

func NewSSH(conf *commonssh.Config, ws *websocket.Conn) *SSH {
	flushCh := make(chan struct{}, 1)
	return &SSH{
		conf:    conf,
		ws:      ws,
		flushCh: flushCh,
		writer: &writer{
			maxSize: maxBufferSize,
			flushCh: flushCh,
		},
	}
}

func (s *SSH) Connect() (*SSH, error) {
	host := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	slog.Info("Attempting to connect SSH, host: %s, port: %d", s.conf.Host, s.conf.Port)
	client, err := commonssh.NewSSHClient(s.conf)
	if err != nil {
		var fingerprintErr *types.FingerprintError
		if errors.As(err, &fingerprintErr) {
			return s, &types.FingerprintError{
				Host:        fingerprintErr.Host,
				Fingerprint: fingerprintErr.Fingerprint,
			}
		}
		slog.Error("SSH connection failed: %v", err)
		return s, err
	}
	slog.Info("SSH connection successful, %s@%s", s.conf.User, host)

	slog.Info("Creating SSH session")
	session, err := client.NewSession()
	if err != nil {
		slog.Error("Failed to create SSH session: %v", err)
		return s, err
	}
	s.session = session

	slog.Debug("Getting session stdin pipe")
	s.stdinPipe, err = s.session.StdinPipe()
	if err != nil {
		slog.Error("Failed to get stdin pipe: %v", err)
		return nil, err
	}

	s.session.Stdout = s.writer
	s.session.Stderr = s.writer
	slog.Debug("Stdout and stderr configured")

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// TODO: 支持自定义终端类型
	slog.Debug("Requesting PTY terminal, type: xterm")
	if err = s.session.RequestPty("xterm", 0, 0, modes); err != nil {
		slog.Error("Failed to request PTY terminal: %v", err)
		return nil, err
	}

	slog.Debug("Starting shell")
	if err = s.session.Shell(); err != nil {
		slog.Error("Failed to start shell: %v", err)
		return nil, err
	}

	slog.Info("SSH session ready")
	return s, nil
}

func (s *SSH) flushWriter() {
	s.pausedMu.Lock()
	paused := s.paused
	s.pausedMu.Unlock()

	if paused {
		s.writer.Reset()
		return
	}

	s.writer.mu.Lock()
	bufferLen := s.writer.buffer.Len()
	s.writer.mu.Unlock()

	if bufferLen == 0 {
		return
	}

	const maxChunkSize = 64 * 1024 // 64KB per chunk
	if bufferLen > maxChunkSize {
		s.writer.mu.Lock()
		data := make([]byte, maxChunkSize)
		n, _ := s.writer.buffer.Read(data)
		s.writer.mu.Unlock()

		if n > 0 {
			if err := s.ws.WriteJSON(&types.Message{
				Type:    enums.TerminalTypeData,
				Content: string(data[:n]),
			}); err != nil {
				slog.Error("failed write data to websocket: %v", err)
			}
		}
	} else {
		s.writer.mu.Lock()
		content := s.writer.buffer.String()
		s.writer.buffer.Reset()
		s.writer.mu.Unlock()

		if len(content) > 0 {
			if err := s.ws.WriteJSON(&types.Message{
				Type:    enums.TerminalTypeData,
				Content: content,
			}); err != nil {
				slog.Error("failed write data to websocket: %v", err)
			}
		}
	}
}

func (s *SSH) Input(quitSignal chan bool) {
	slog.Info("Starting WebSocket input monitoring")
	defer s.setQuit(quitSignal)

	for {
		select {
		case <-quitSignal:
			return
		default:
			_, data, err := s.ws.ReadMessage()
			if err != nil {
				return
			}
			msg := &terminal.Payload{}
			_ = json.Unmarshal(data, &msg)

			switch msg.Type {
			case enums.TerminalTypeResize:
				if msg.Cols > 0 && msg.Rows > 0 {
					if err = s.session.WindowChange(msg.Rows, msg.Cols); err != nil {
						slog.Error("failed change ssh pty window size: %v", err)
					}
				}
			case enums.TerminalTypeCMD:
				if _, err = s.stdinPipe.Write([]byte(msg.Cmd)); err != nil {
					slog.Error("failed write command to stdin pipe: %v", err)
				}
			case enums.TerminalTypeFlowControl:
				if msg.Pause != nil {
					s.pausedMu.Lock()
					s.paused = *msg.Pause
					s.pausedMu.Unlock()

					s.writer.pausedMu.Lock()
					s.writer.paused = *msg.Pause
					s.writer.pausedMu.Unlock()

					if *msg.Pause {
						slog.Debug("Flow control: paused")
						s.writer.Reset()
					} else {
						slog.Debug("Flow control: resumed")
					}
				}
			}
		}
	}
}

func (s *SSH) Output(quitSignal chan bool) {
	slog.Info("Starting WebSocket output")
	defer s.setQuit(quitSignal)

	fallbackTick := time.NewTicker(50 * time.Millisecond)
	defer fallbackTick.Stop()

	cleanupTick := time.NewTicker(time.Second)
	defer cleanupTick.Stop()

	for {
		select {
		case <-quitSignal:
			s.flushWriter()
			return
		case <-s.flushCh:
			// 单次唤醒最多处理 1MB，避免长时间占用
			const maxFlushBytes = 1024 * 1024 // 1MB
			flushedBytes := 0
			for flushedBytes < maxFlushBytes {
				s.writer.mu.Lock()
				beforeLen := s.writer.buffer.Len()
				s.writer.mu.Unlock()

				if beforeLen == 0 {
					break
				}

				s.flushWriter()

				s.writer.mu.Lock()
				afterLen := s.writer.buffer.Len()
				s.writer.mu.Unlock()

				flushedBytes += beforeLen - afterLen

				if afterLen == 0 {
					break
				}
			}
		case <-cleanupTick.C:
			// 暂停时定期释放 buffer 底层容量
			s.pausedMu.Lock()
			paused := s.paused
			s.pausedMu.Unlock()

			if paused {
				s.writer.mu.Lock()
				if s.writer.buffer.Len() > 0 {
					s.writer.buffer.Reset()
					oldCap := s.writer.buffer.Cap()
					if oldCap > maxBufferSize {
						s.writer.buffer = bytes.Buffer{}
					}
				}
				s.writer.mu.Unlock()
			}
		case <-fallbackTick.C:

			s.flushWriter()
		}
	}
}

func (s *SSH) close() {
	slog.Info("Closing SSH session")
	if s.session != nil {
		_ = s.session.Close()
	}
}

func (s *SSH) Wait(quitSignal chan bool) {
	defer s.close()
	defer s.setQuit(quitSignal)
	_ = s.session.Wait()
}

func (s *SSH) setQuit(ch chan bool) {
	ch <- true
}
