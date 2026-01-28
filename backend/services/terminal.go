package services

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"

	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal/adapter"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
)

var TerminalSrvSet = wire.NewSet(wire.Struct(new(TerminalSrv), "*"))

type TerminalSrv struct {
	ConnectionSrv    *ConnectionSrv
	MetadataSrv      *MetadataSrv
	HTTPListenerPort *initialize.HTTPListenerPort
}

func (s *TerminalSrv) SSH(ws *websocket.Conn, hostID uint) error {
	slog.Info("Starting SSH connection, hostID: %d", hostID)
	conn, err := s.ConnectionSrv.FindByID(hostID)
	if err != nil {
		slog.Error("Failed to find host information: %v, hostID: %d", err, hostID)
		return fmt.Errorf("failed to find host: %v", err)
	}

	slog.Info("Found host information, host: %s, port: %d", conn.Host, conn.Port)
	if conn.Metadata == nil {
		slog.Info("Host metadata is empty, starting metadata update")
		go s.MetadataSrv.UpdateByConnection(conn)
	}

	sshConf := &commonssh.Config{
		Host:       conn.Host,
		Port:       conn.Port,
		User:       conn.Credential.Username,
		AuthMethod: conn.Credential.AuthMethod,
		Password:   conn.Credential.Password,
		PrivateKey: conn.Credential.PrivateKey,
		Passphrase: conn.Credential.Passphrase,
	}

	// if len(conn.SSHCiphers) > 0 {
	// 	sshConf.Ciphers = conn.SSHCiphers
	// 	slog.Debug("Using custom ciphers: %v", conn.SSHCiphers)
	// }
	//
	// if len(conn.SSHKeyExchanges) > 0 {
	// 	sshConf.KeyExchanges = conn.SSHKeyExchanges
	// 	slog.Debug("Using custom key exchanges: %v", conn.SSHKeyExchanges)
	// }
	//
	// if len(conn.SSHMACs) > 0 {
	// 	sshConf.MACs = conn.SSHMACs
	// 	slog.Debug("Using custom MACs: %v", conn.SSHMACs)
	// }
	//
	// if len(conn.SSHPublicKeyAlgorithms) > 0 {
	// 	sshConf.PublicKeyAlgorithms = conn.SSHPublicKeyAlgorithms
	// 	slog.Debug("Using custom public key algorithms: %v", conn.SSHPublicKeyAlgorithms)
	// }
	//
	// if len(conn.SSHHostKeyAlgorithms) > 0 {
	// 	sshConf.HostKeyAlgorithms = conn.SSHHostKeyAlgorithms
	// 	slog.Debug("Using custom host key algorithms: %v", conn.SSHHostKeyAlgorithms)
	// }
	//
	// if conn.SSHCharset != "" {
	// 	sshConf.Charset = conn.SSHCharset
	// 	slog.Debug("Using charset: %s", conn.SSHCharset)
	// }

	slog.Info("SSH configuration ready, host: %s, user: %s, auth method: %s",
		conn.Host,
		conn.Credential.Username,
		conn.Credential.AuthMethod)

	slog.Info("Connecting to SSH server, host: %s, port: %d", conn.Host, conn.Port)
	ssh, err := adapter.NewSSH(sshConf, ws).Connect()
	if err != nil {
		slog.Error("SSH connection failed: %v, host: %s, port: %d", err, conn.Host, conn.Port)
		return err
	}
	slog.Info("SSH connection successful, host: %s, port: %d", conn.Host, conn.Port)

	// 发送连接成功消息
	if err = ws.WriteJSON(&types.Message{Type: enums.TerminalTypeConnected}); err != nil {
		slog.Error("Failed to send connection success message: %v", err)
		return err
	}
	slog.Info("Connection success message sent")

	term := terminal.NewTerminal(ws, ssh, s.SessionEnded)
	slog.Info("Starting terminal session, host: %s, port: %d", conn.Host, conn.Port)
	term.Start()

	return nil
}

func (s *TerminalSrv) AddFingerprint(hostID uint, host string, fingerprint string) error {
	slog.Info("Adding host fingerprint, hostID: %d, host: %s, fingerprint: %s", hostID, host, fingerprint)
	conn, err := s.ConnectionSrv.FindByID(hostID)
	if err != nil {
		slog.Error("Failed to find host information: %v, hostID: %d", err, hostID)
		return fmt.Errorf("failed to find host: %v", err)
	}
	sshConf := &commonssh.Config{
		Host: conn.Host,
		Port: conn.Port,
		User: conn.Credential.Username,
	}
	if err = commonssh.AddFingerprint(sshConf, host, fingerprint); err != nil {
		slog.Error("Failed to add host fingerprint: %v, host: %s", err, host)
		return fmt.Errorf("failed to add host fingerprint: %v", err)
	}
	slog.Info("Successfully added host fingerprint, host: %s", host)
	return nil
}

// func (s *TerminalSrv) Serial(ws *websocket.Conn) error {
// 	serial := adapter.NewSerial(ws, s.Logger)
//
// 	// test code
// 	serialPort := "/dev/cu.usbserial-2130"
//
// 	err := serial.Open(serialPort)
// 	if err != nil {
// 		return fmt.Errorf("failed to open serial port: %v", err)
// 	}
//
// 	term := terminal.NewTerminal(ws, serial, s.closeWsWrapper)
// 	term.Start()
//
// 	return nil
// }

func (s *TerminalSrv) SerialPorts() *resp.Resp {
	slog.Info("Getting available serial ports")
	ports, err := serial.GetPortsList()
	if err != nil {
		slog.Error("Failed to get serial port list: %v", err)
		return resp.FailWithMsg(err.Error())
	}
	slog.Info("Found %d available serial ports", len(ports))
	return resp.OkWithData(ports)
}

func (s *TerminalSrv) CloseSession(ws *websocket.Conn, reason string) {
	slog.Info("Closing session, reason: %s", reason)
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	err := ws.WriteControl(websocket.CloseMessage, data, time.Now().Add(consts.WebSocketWriteWait))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		slog.Error("Failed to close session: %v, forcibly closing connection", err)
		// If close message could not be sent, then close without the handshake.
		_ = ws.Close()
	}
}

func (s *TerminalSrv) SessionEnded(ws *websocket.Conn) {
	s.CloseSession(ws, messages.SessionEnded)
}

func (s *TerminalSrv) WebsocketPort() int {
	port := int(*s.HTTPListenerPort)
	slog.Debug("WebSocket service port: %d", port)
	return port
}
