package services

import (
	"log/slog"

	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"github.com/MisakaTAT/GTerm/backend/pkg/metadata"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var MetadataSrvSet = wire.NewSet(wire.Struct(new(MetadataSrv), "*"))

type MetadataSrv struct {
	Query *query.Query
}

func (s *MetadataSrv) UpdateByConnection(conn *model.Connection) {
	t := s.Query.Metadata

	config := &commonssh.Config{
		Host:             conn.Host,
		Port:             conn.Port,
		User:             conn.Credential.Username,
		AuthMethod:       conn.Credential.AuthMethod,
		Password:         conn.Credential.Password,
		PrivateKey:       conn.Credential.PrivateKey,
		Passphrase:       conn.Credential.Passphrase,
		TrustUnknownHost: true,
	}
	client, err := exec.NewExec(config)
	if err != nil {
		slog.Error("failed to create ssh client", zap.Error(err))
		return
	}
	defer func() {
		_ = client.Close()
	}()

	meta, err := t.Where(t.ConnectionID.Eq(conn.ID)).FirstOrInit()
	if err != nil {
		slog.Error("failed to get metadata", zap.Error(err))
		return
	}

	metaInfo := metadata.NewMetadata(client).Parser()
	if metaInfo != nil {
		meta.Vendor = metaInfo.Vendor
		meta.Type = metaInfo.Type
	}

	if err = t.Save(meta); err != nil {
		slog.Error("failed to update metadata", zap.Error(err))
	}
}
