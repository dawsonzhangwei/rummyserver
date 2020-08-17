package services

import (
	"context"
	
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
)

type Auth struct {
	component.Base	
} 

type MsgLogin struct {
	Uid string `json:"uid"`
	Nick string `json:"nick"`
}

func NewAuthen() *Auth {
	return &Auth{}
}

func (a *Auth) login(ctx context.Context, msg *MsgLogin) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx) // The default logger contains a requestId, the route being executed and the sessionId
	s := pitaya.GetSessionFromCtx(ctx)
	s.Bind(ctx, msg.Uid)

	logger.Infof("User login, msg:%v", msg)
}