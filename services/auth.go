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

type LoginResponse struct {
	Code int `json:"code"`
	Result string `json:"result"`	
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(ctx context.Context, msg *MsgLogin)(*LoginResponse, error) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx) // The default logger contains a requestId, the route being executed and the sessionId
	s := pitaya.GetSessionFromCtx(ctx)
	s.Bind(ctx, msg.Uid)

	logger.Infof("User login, msg:%v", msg)

	return &LoginResponse {Result: "success"}, nil
}