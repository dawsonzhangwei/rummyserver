package services

import (
	"context"
	
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"

	"rummy/msg"
	"rummy/modules"
)

type Auth struct {
	component.Base
} 

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(ctx context.Context, msg *msg.LoginReq)(*msg.LoginRsp, error) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx) // The default logger contains a requestId, the route being executed and the sessionId
	s := pitaya.GetSessionFromCtx(ctx)
	s.Bind(ctx, msg.Uid)

	logger.Infof("User login, msg:%v", msg)

	loginParam := make(map[string]string)
	loginParam["uid"] = msg.GetUid()
	loginParam["channel"] = msg.GetChannel()

	channel := a.getChannel(msg.GetChannel())
	if channel == nil {
		return &msg.LoginRsp {
			Code : 500,
			Result : fmt.Sprintf("channel:%v is not exist", msg.GetChannel()),
		}, nil
	}

	player, err := channel.Login(loginParam)
	if err != nil {
		return &msg.LoginRsp {
			Code : 500,
			Result: err.Error(),
		}, nil
	}

	cache := getDataCache()

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