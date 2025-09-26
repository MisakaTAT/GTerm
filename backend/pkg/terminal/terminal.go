package terminal

import (
	"sync"

	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/gorilla/websocket"
)

type Handler interface {
	Input(quitSignal chan bool)
	Output(quitSignal chan bool)
	Wait(quitSignal chan bool)
}

type Payload struct {
	Type enums.TerminalType `json:"type"`
	Cmd  string             `json:"cmd"`
	Cols int                `json:"cols"`
	Rows int                `json:"rows"`
}

type Terminal struct {
	handler  Handler
	ws       *websocket.Conn
	stopFunc func(ws *websocket.Conn)
}

func NewTerminal(ws *websocket.Conn, handler Handler, stopFunc func(ws *websocket.Conn)) *Terminal {
	return &Terminal{
		ws:       ws,
		handler:  handler,
		stopFunc: stopFunc,
	}
}

func (t *Terminal) Start() {
	var wg sync.WaitGroup
	wg.Add(3)

	quitSignal := make(chan bool, 3)
	outputClosed := make(chan struct{})

	go func() {
		defer wg.Done()
		t.handler.Input(quitSignal)
	}()

	go func() {
		defer wg.Done()
		t.handler.Output(quitSignal)
		close(outputClosed)
	}()

	go func() {
		defer wg.Done()
		t.handler.Wait(quitSignal)
	}()

	go func() {
		// ensure flush write buffer before close
		// must output handler close before with websocket close
		<-outputClosed
		t.stopFunc(t.ws)
	}()

	wg.Wait()
}
