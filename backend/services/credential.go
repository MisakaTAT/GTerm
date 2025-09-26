package services

import (
	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
)

var CredentialSrvSet = wire.NewSet(wire.Struct(new(CredentialSrv), "*"))

type CredentialSrv struct {
	Query *query.Query
}

func (s *CredentialSrv) CreateCredential(cred *model.Credential) *resp.Resp {
	t := s.Query.Credential
	cred.IsCommonCredential = true
	if err := t.Create(cred); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.CreateSuccess)
}

func (s *CredentialSrv) UpdateCredential(cred *model.Credential) *resp.Resp {
	t := s.Query.Credential
	if _, err := t.Where(t.ID.Eq(cred.ID)).Updates(cred); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.UpdateSuccess)
}

func (s *CredentialSrv) ListCredential() *resp.Resp {
	t := s.Query.Credential
	credentials, err := t.Where(t.IsCommonCredential.Is(true)).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(credentials)
}

func (s *CredentialSrv) DeleteCredential(id uint) *resp.Resp {
	t := s.Query.Credential
	_, err := t.Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.DeleteSuccess)
}

func (s *CredentialSrv) FindCredentialByID(id uint) *resp.Resp {
	conn, err := s.FindByID(id)
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(conn)
}

func (s *CredentialSrv) FindByID(id uint) (*model.Credential, error) {
	t := s.Query.Credential
	cred, err := t.Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	if err = cred.Decrypt(); err != nil {
		return nil, err
	}
	return cred, nil
}
