package cmd

import (
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/services"
	"github.com/google/wire"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var AppSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	HTTPListenerPort *initialize.HTTPListenerPort
	TerminalSrv      *services.TerminalSrv
	PreferencesSrv   *services.PreferencesSrv
	GroupSrv         *services.GroupSrv
	ConnectionSrv    *services.ConnectionSrv
	MetadataSrv      *services.MetadataSrv
	CredentialSrv    *services.CredentialSrv
	WebsocketSrv     *services.WebsocketSrv
}

func (a *App) Services() (servers []application.Service) {
	servers = append(servers, application.NewService(a.TerminalSrv))
	servers = append(servers, application.NewService(a.PreferencesSrv))
	servers = append(servers, application.NewService(a.GroupSrv))
	servers = append(servers, application.NewService(a.ConnectionSrv))
	servers = append(servers, application.NewService(a.MetadataSrv))
	servers = append(servers, application.NewService(a.CredentialSrv))
	servers = append(servers, application.NewService(a.WebsocketSrv))
	return
}
