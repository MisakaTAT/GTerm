package services

import (
	"fmt"
	"runtime"
	"time"

	"github.com/MisakaTAT/GTerm/backend/pkg/base"
	"github.com/google/wire"
)

var PreferencesSrvSet = wire.NewSet(wire.Struct(new(PreferencesSrv), "*"))

type PreferencesSrv struct {
}

func (s *PreferencesSrv) Version() string {
	return base.Version
}

func (s *PreferencesSrv) VersionURL() string {
	return base.VersionURL
}

func (s *PreferencesSrv) Copyright() string {
	return fmt.Sprintf("Copyright © 2024-%d MisakaTAT.\nAll rights reserved.", time.Now().Year())
}

func (s *PreferencesSrv) GOOS() string {
	return runtime.GOOS
}

func (s *PreferencesSrv) IsDarwin() bool {
	return s.GOOS() == "darwin"
}
