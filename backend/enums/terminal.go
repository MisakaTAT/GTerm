package enums

type TerminalType string

const (
	TerminalTypeError              TerminalType = "Error"
	TerminalTypeData               TerminalType = "Data"
	TerminalTypeConnected          TerminalType = "Connected"
	TerminalTypeFingerprintConfirm TerminalType = "FingerprintConfirm"
	TerminalTypeResize             TerminalType = "Resize"
	TerminalTypeCMD                TerminalType = "CMD"
	TerminalTypeFlowControl        TerminalType = "FlowControl"
)
