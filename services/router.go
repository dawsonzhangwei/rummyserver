package services

import (
	"context"
	
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"

	"rummy/modules"
)

type (
	Router struct {
		component.Base

		msgBus *modules.MessageBus
	}

	RouterReply struct {
		Code int `json:"code"`
		Result string `json:"result"`
	}
)

func NewRouter() * Router {
	return &Router{}
}

// AfterInit was called after the component is initialized.
func (r *Router) AfterInit() {
}

func (r *Router) C2G(ctx context.Context, data []byte) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx)

	if r.msgBus == nil {
		m, err := pitaya.GetModule("messageBus")
		if err != nil {
			logger.Errorf("Router, get messageBus failed, err:%v", err)
		}
		r.msgBus, _ = m.(*modules.MessageBus)
	}

	if r.msgBus == nil {
		logger.Errorf("Router, C2G msg:%v, msgBus is nil", string(data))
		return
	}

	s := pitaya.GetSessionFromCtx(ctx)
	
	r.msgBus.SendMsgToGameX(s.UID(), data, "c2g")

	logger.Debugf("c2g:%v", string(data))
}