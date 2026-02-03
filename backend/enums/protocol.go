package enums

type ConnProtocol string

const (
	SSH    ConnProtocol = "SSH"
	Telnet ConnProtocol = "Telnet"
	RDP    ConnProtocol = "RDP"
	VNC    ConnProtocol = "VNC"
	Serial ConnProtocol = "Serial"
)
