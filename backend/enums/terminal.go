package enums

import "strings"

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

var TerminalTypeEnums = []TerminalType{TerminalTypeError, TerminalTypeData, TerminalTypeConnected, TerminalTypeFingerprintConfirm, TerminalTypeResize, TerminalTypeCMD, TerminalTypeFlowControl}

func (a TerminalType) TSName() string {
	return strings.ToUpper(string(a))
}
