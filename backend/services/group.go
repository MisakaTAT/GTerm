package services

import (
	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
)

var GroupSrvSet = wire.NewSet(wire.Struct(new(GroupSrv), "*"))

type GroupSrv struct {
	Query *query.Query
}

func (s *GroupSrv) CreateGroup(group *model.Group) *resp.Resp {
	t := s.Query.Group
	if err := t.Create(group); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.CreateSuccess)
}

func (s *GroupSrv) UpdateGroup(group *model.Group) *resp.Resp {
	t := s.Query.Group
	if _, err := t.Where(t.ID.Eq(group.ID)).Updates(group); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.UpdateSuccess)
}

func (s *GroupSrv) DeleteGroup(id uint) *resp.Resp {
	if err := s.Query.Transaction(func(tx *query.Query) error {
		_, err := tx.Connection.Where(tx.Connection.GroupID.Eq(id)).UpdateSimple(tx.Connection.GroupID.Null())
		if err != nil {
			return err
		}
		_, err = tx.Group.Where(tx.Group.ID.Eq(id)).Delete()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithCode(messages.DeleteSuccess)
}

func (s *GroupSrv) ListGroup() *resp.Resp {
	t := s.Query.Group
	groups, err := t.Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(groups)
}
