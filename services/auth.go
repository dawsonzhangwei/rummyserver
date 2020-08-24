package services

import (
	"context"
	
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"

	"rummy/channel"
	"rummy/msg"
	"rummy/modules"
	"rummy/gameConst"
	"rummy/base"
)

type Auth struct {
	component.Base
} 

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(ctx context.Context, req *msg.LoginReq)(*msg.LoginRsp, error) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx) // The default logger contains a requestId, the route being executed and the sessionId
	s := pitaya.GetSessionFromCtx(ctx)
	s.Bind(ctx, req.Uid)

	logger.Infof("User login, msg:%v", req)

	loginParam := make(map[string]string)
	loginParam["uid"] = req.GetUid()
	loginParam["channel"] = req.GetChannel()

	/*
	channel := a.getChannel(req.GetChannel())
	if channel == nil {
		return &msg.LoginRsp {
			Code : 500,
			Result : fmt.Sprintf("channel:%v is not exist", req.GetChannel()),
		}, nil
	}

	player, err := channel.Login(loginParam)
	if err != nil {
		return &msg.LoginRsp {
			Code : 500,
			Result: err.Error(),
		}, nil
	}
	*/

	player := &base.Player{}
	player.UID = req.GetUid()

	modules.DB.LoadPlayer(player)
	modules.Cache.AddPlayer(player)

	return &msg.LoginRsp {Code: 200, Result: "success"}, nil
}

func (a *Auth) getChannel(strChannel string) (c channel.ThirdChannel) {
	switch strChannel {
	case gameConst.PAW:
		c = &channel.Paw{}
	default:
		c = nil
	}

	return c
}