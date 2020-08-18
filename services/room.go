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

	PlayerMsg struct {
		Clover string `json:"clover"`
		Life string `json:"life"`
		Coin string `json:"coin"`
		Ticket string `json:"ticket"`
		Level string `json:"level"`
		Check_point string `json:"check_point"`
		PlayerId string `json:"playerId"`
		PlayerName string `json:"playerName"`
		Current_game_server_id string `json:"current_game_server_id"`
		Current_room_id string `json:"current_room_id"`
		Current_game_id string `json:"current_game_id"`
		Current_field_id string `json:"current_field_id"`
		Is_admin string `json:"is_admin"`
		PetId string `json:"petId"`
		Current_gate_id string `json:"current_gate_id"`
		Current_match_status string `json:"current_match_status"`
		Play_num string `json:"play_num"`
		Win_num string `json:"win_num"`
		Item_ddz string `json:"item_ddz"`
		Played_player string `json:"played_player"`
	}

	EnterGame struct {
		Content string `json:"content"`
		Player []PlayerMsg `json:"player"`
		GToken string `json:"player"`
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
