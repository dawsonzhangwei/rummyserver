package services

import (
	"fmt"
	"time"
	"context"
	"errors"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/logger"
	"github.com/topfreegames/pitaya/timer"
	"github.com/topfreegames/pitaya/groups"
	"github.com/topfreegames/pitaya/config"

	"rummy/base"
	"rummy/util"
	"rummy/modules"
	"rummy/gameConst"
	"rummy/msg"
)
const (
	ServerType = "connector"
)

type (
	Room struct {
		component.Base
		timer *timer.Timer
	}

	JoinReq struct {
		Code int `json:"roomId"`
		IsSeen int `json:"isSeen"`
		ServerId int `json:"serverId"`	
	}

	JoinResponse struct {
		Code int `json:"code"`
		Msg string `json:"msg"`	
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

	EnterGame struct {
		Player []base.Player `json:"player"`
		Token string `json:"token"`
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

func (r *Room) Join(ctx context.Context, req *JoinReq) (*JoinResponse, error) {
	s := pitaya.GetSessionFromCtx(ctx)

	pitaya.GroupAddMember(ctx, "room", s.UID())

	uids, err := pitaya.GroupMembers(ctx, "room")
	if err != nil {
		return nil, err
	}

	if len(uids) > 1 {
		r.matchSuccess(uids);
		pitaya.GroupBroadcast(ctx, ServerType, "room", "onNewUser", &NewUser{Content: fmt.Sprintf("New user: %v", s.UID())})
	}
	

	//s.Push("onMembers", &AllMembers{Members: uids})
	

	s.OnClose(func() {
		pitaya.GroupRemoveMember(ctx, "room", s.UID())
	})

	return &JoinResponse{Msg: "success"}, nil
}

func (r *Room) Message(ctx context.Context, msg *UserMessage) {
	err := pitaya.GroupBroadcast(ctx, ServerType, "room", "onMessage", msg)
	if err != nil {
		fmt.Printf("GroupBroadcast err:%v", err)	
	}
}

func (r *Room) matchSuccess(uids []string) error {
	roomId := "6001"
	gid := 30011

	var players []*base.Player
	var attrs []*msg.PlayerAttr
	var tokens []string
	for _, uid := range uids {
		player := modules.Cache.GetPlayer(uid)
		if player == nil {
			logger.Log.Errorf("matchSuccess getPlayer failed by id:%v", uid)
			return errors.New("GetPlayer failed")
		}
		token := util.GenerateToken(uid)
		
		util.InsertTokenInfo(token, uid, roomId)
		util.SavePlayerInfoToGameCache(uid, player.GetCoin(), roomId, token)

		players = append(players, player)
		attrs = append(attrs, &player.PlayerAttr)
		tokens = append(tokens, token)
	}

	err := util.SaveRoomInfoToGameCache(gid, roomId, players, tokens)
	if err != nil {
		return err
	}

	for idx, uid := range uids {
		ntf := &msg.EnterGameNtf {
			Players : attrs,
			Token : tokens[idx],
		}
		pitaya.SendPushToUsers(gameConst.PushCmd_EnterGame, ntf, []string{uid}, "connector")
	}

	return nil
}