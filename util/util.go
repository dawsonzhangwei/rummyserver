package util

import (
	"fmt"
	"time"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/logger"

	"rummy/base"
	"rummy/redis"
	"rummy/gameConst"
)

func getPlayerGameCacheKey(uid string) string {
	return fmt.Sprintf("%s%s", gameConst.GAME_CACHE_PLAYER_CACHE, uid)
}

func SavePlayerInfoToGameCache(uid string, coin int, room_id string, token string) error {
	client, err := redis.GetRedis(gameConst.REDIS_GAME_CACHE, pitaya.GetConfig())
	if err != nil {
		return err	
	} else {
		playerData := map[string]interface{}{
			"ID": uid,
			"OpenID": uid,
			"MoMoID": uid,
			"Gold": coin,
			"CurrentRoomID": room_id,
			"gToken": token,
			"package": "",
			"is_robot": 0,
			"IsSeen": 0, //明牌开始
			"IsOwner": 0, //是否房主
			"WinOrLose": 0,
		}
		_, err := client.HMSet(getPlayerGameCacheKey(uid), playerData).Result()
		return err
	}
}

func getGameConfigField(gid int) (string, error) {
	conf := map[string] interface{}{}
	secMap := map[string] interface{} {
		"game_id": fmt.Sprintf("2-1-1-1-%v", gid),
	}
	conf["game_base"] = secMap
	
	byteConf, err := json.Marshal(conf)
	return string(byteConf), err
}

func getRoomGameCacheKey(gid int, roomId string) string {
	return fmt.Sprintf("%v%v_%v", gameConst.GAME_CACHE_ROOM_CACHE, gid, roomId)
}

func myMD5(params ...string) string {
	var buffer bytes.Buffer
	for _, v := range params {
		buffer.WriteString(v)
	}
	s := buffer.String()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	str := hex.EncodeToString(md5Ctx.Sum(nil))
	return str
}
	
func SaveRoomInfoToGameCache(gid int, roomId string, players []*base.Player, tokens []string) error {
	client, err := redis.GetRedis(gameConst.REDIS_GAME_CACHE, pitaya.GetConfig())
	if err != nil {
		return err	
	} else {
		strPlayers, err := json.Marshal(players)
		if err != nil {
			logger.Log.Errorf("saveRoomInfoToGameCache, json.Marshal err:%v, players:%v", err, players)
			return err
		}

		strTokens, err := json.Marshal(tokens)
		if err != nil {
			logger.Log.Errorf("saveRoomInfoToGameCache, json.Marshal err:%v, tokens:%v", err, tokens)
			return err
		}

		gameConf, err := getGameConfigField(gid)
		if err != nil {
			logger.Log.Errorf("saveRoomInfoToGameCache, getGameConfigField, err:%v, gid:%v", err, gid)
			return err
		}

		roomData := map[string]interface{}{
			"ID": roomId,
			"ShowID": roomId,
			"GameID": fmt.Sprintf("2-1-1-1-%v", gid),
			"GameType": 2,
			"GameSvrID": gid,
			"GameConfig": gameConf,
			"playerList": string(strPlayers),
			"GTokenList": string(strTokens),
			"HallSvrId": 0, 
			"DateTime": time.Now().Unix(), 
			"fieldId": gameConst.GAME_FIELD_ID_TWO,
			"uniqueId": myMD5(fmt.Sprintf("%v%v", time.Now().Unix(), roomId)),
			"AgainstType": 0,
		}

		roomKey := getRoomGameCacheKey(gid, roomId)
		_, err1 := client.HMSet(roomKey, roomData).Result()
		logger.Log.Infof("saveRoomInfoToGameCache:roomKey:%v, data%v", roomKey, roomData)
		return err1
	}
}

func InsertTokenInfo( token_id string, uid string, room_id string) error {
	client, err := redis.GetRedis(gameConst.REDIS_GAME_CACHE, pitaya.GetConfig())
	if err != nil {
		return err	
	} else {
		tokenData := map[string]interface{}{
			"TokenID": token_id,
			"Owner": uid,
			"Time": time.Now().Unix(),
			"RoomID": room_id,
			"Ip": "127.0.0.1",
			"ShowID": room_id,
			"Position": "0",
			"RoomAllocTime": time.Now().Unix(),
			"RedisIp": "127.0.0.1", 
			"RedisPort": "6379", 
		}

		tokenKey := fmt.Sprintf("%v%v", gameConst.GAME_CACHE_G_TOKEN, token_id)
		_, err := client.HMSet(tokenKey, tokenData).Result()
		logger.Log.Infof("insertTokenInfo:tokenKey:%v, data%v", tokenKey, tokenData)
		return err
	}
}

func storePlayerInfo(uid string, data map[string]interface{}) error {
	client, err := redis.GetRedis(gameConst.REDIS_AGC_CACHE, pitaya.GetConfig())
	if err != nil {
		return err	
	} else {
		playerKey := fmt.Sprintf("%v%v", gameConst.PLAYER_PLAYER_INFO, uid)
		_, err := client.HMSet(playerKey, data).Result()
		logger.Log.Infof("storePlayerInfo, playerKey:%v, data%v", playerKey, data)
		return err
	}
}

func GenerateToken(key string) string {
	return myMD5(key)
}