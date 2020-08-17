package services

import (
	"fmt"
	"time"
	"context"
	"strconv"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/logger"
	"github.com/topfreegames/pitaya/timer"
	"github.com/topfreegames/pitaya/groups"
	"github.com/topfreegames/pitaya/config"
)
const (
	ServerType = "connector"
)

type (
	Room struct {
		component.Base
		timer *timer.Timer
	}

	JoinResponse struct {
		Code int `json:"code"`
		Result string `json:"result"`	
	}

	UserMessage struct {
		Name string `json:"name"`
		Content string `json:"content"`	
	}

	NewUser struct {
		Content string `json:"content"`
	}

	AllMembers struct {
		Members []string `json:"members"` 	
	}
)

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) Init() {
	gsi := groups.NewMemoryGroupService(config.NewConfig())
	pitaya.InitGroups(gsi)
	pitaya.GroupCreate(context.Background(), "room")
}

func (r *Room) AfterInit() {
	r.timer = pitaya.NewTimer(time.Minute, func() {
		count, err := pitaya.GroupCountMembers(context.Background(), "room")
		logger.Log.Debugf("%v UserCount:%v, error:%v", time.Now().String(), count, err)
	})
}

func (r *Room) Join(ctx context.Context, msg []byte) (*JoinResponse, error) {
	s := pitaya.GetSessionFromCtx(ctx)
	FakeUID := s.ID()
	err := s.Bind(ctx, strconv.Itoa(int(FakeUID)))
	if err != nil {
		return nil, pitaya.Error(err, "RH-000", map[string] string{"failed": "bind"})
	}

	uids, err := pitaya.GroupMembers(ctx, "room")
	if err != nil {
		return nil, err
	}
	s.Push("onMembers", &AllMembers{Members: uids})
	pitaya.GroupBroadcast(ctx, ServerType, "room", "onNewUser", &NewUser{Content: fmt.Sprintf("New user: %v", s.UID())})
	pitaya.GroupAddMember(ctx, "room", s.UID())

	s.OnClose(func() {
		pitaya.GroupRemoveMember(ctx, "room", s.UID())
	})

	return &JoinResponse{Result: "success"}, nil
}

func (r *Room) Message(ctx context.Context, msg *UserMessage) {
	err := pitaya.GroupBroadcast(ctx, ServerType, "room", "onMessage", msg)
	if err != nil {
		fmt.Printf("GroupBroadcast err:%v", err)	
	}
}
