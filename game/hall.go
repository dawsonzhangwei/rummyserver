package game

import (
	"fmt"
	"strconv"	
)

var (
	roomList map[string] *GameRoom
)

func init() {
	roomIdList := [] int {
		10001, 10002, 10003,
	}

	roomList = make(map[string] *GameRoom)
	for id := range roomIdList {
		roomList[strconv.Itoa(id)] = new(GameRoom)
	}
}

func GetRoom(roomId string) (*GameRoom, error) {
	if r, ok := roomList[roomId]; ok {
		return r, nil
	}

	return nil, fmt.Errorf("cannout find room by id:%v", roomId)
}
