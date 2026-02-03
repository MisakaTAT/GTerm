package services

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ConnectionSrvSet,
	CredentialSrvSet,
	GroupSrvSet,
	MetadataSrvSet,
	PreferencesSrvSet,
	TerminalSrvSet,
	WebsocketSrvSet,
)
