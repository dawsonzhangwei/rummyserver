package modules

import (
	"time"
	"fmt"
	"errors"
	"sync/atomic"
	"encoding/json"

	"github.com/go-redis/redis"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/logger"
	"github.com/topfreegames/pitaya/modules"

	"rummy/gameConst"
)

type RouteMessage struct {
	Sid int `json:"sid"`
	Uid string `json:"uid"`
	Dir string `json:"router"`
	Data []byte `json:"data"`
	GType int `json:"gtype,omitempty"`
	Gid int `json:"gid,omitempty"`
}

type MessageBus struct {
	modules.Base

	config *config.Config

	sid int
	gamesid int

	redisAddr string
	redisDbIndex int

	pubsubRedis *redis.Client

	writeChan chan []byte
	readChan chan *RouteMessage

	msgCount uint32

	stopChan chan bool
}

func NewMessageBus(
	config *config.Config,
) (*MessageBus, error) {
	mb := &MessageBus{
		config: config,
		stopChan: make(chan bool),
		writeChan: make(chan []byte, 1000),
		readChan: make(chan *RouteMessage, 1000),
	}
	if err := mb.configure(); err != nil {
		return nil, err
	}

	return mb, nil
}

func (mb *MessageBus) configure() error {
	mb.redisAddr = mb.config.GetString("rummy.redis.addr")
	if mb.redisAddr == "" {
		return errors.New("redisAddr is nil")
	}

    mb.redisDbIndex = mb.config.GetInt("rummy.redis.db")
	return nil
}

func (mb *MessageBus) sendMsgToGamex() error {
	return nil
}

func (mb *MessageBus) Init() error {
	err := mb.createRedisClient()
	if err != nil {
		go mb.writeData()
		go mb.readData()
		go mb.handleMessage()
	}
	return err
}

func (mb *MessageBus) createRedisClient() error {
	logger.Log.Infof("[Pubsub] start connect redis:%v, db:%v", mb.redisAddr, mb.redisDbIndex)

	rdb := redis.NewClient(&redis.Options{
		Addr: mb.redisAddr,
		Password: "",
		DB: mb.redisDbIndex,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		logger.Log.Errorf("[Pubsub] redis Ping failed, err:%v", err)
	} else {
		mb.pubsubRedis = rdb
		logger.Log.Infof("redis conn succeed.")
	}

	return err;
}

func (mb *MessageBus) msgPublish() {
    t := time.NewTicker(time.Millisecond * 10)
	defer t.Stop()

	channel := fmt.Sprintf("MQ_CHANNEL_%v", mb.gamesid)
	logger.Log.Infof("%d publish start, channel:[%v]", mb.sid, channel)

	for {
		select {
		case <-mb.stopChan:
			logger.Log.Infof("%d goroutine publish exit, channel:[%v]", mb.sid, channel)
			return
		case <-t.C:
			count := atomic.SwapUint32(&mb.msgCount, 0)
			if count > 0 {
				msg := fmt.Sprintf("%v|%v|%v", mb.gamesid, mb.sid, count)
				_, err := mb.pubsubRedis.Publish(channel, msg).Result()
				if err != nil {
					logger.Log.Errorf("[%v] Publish failed, err:%v, channel:%v, msg:%v", mb.sid, err, channel, msg)
				} else {
					logger.Log.Debugf("[%v] Publish channel:%v, msg:%v", mb.sid, channel, msg)
				}
			}
		}
	}
}

func (mb *MessageBus) writeData() {
	defer(func(){
		close(mb.writeChan)
	})()

	qout := fmt.Sprintf("mq_%v", mb.gamesid)
	logger.Log.Infof("%d goroutine write start, qout:[%v]", mb.sid, qout)
	
	for str := range mb.writeChan {
		if len(str) == 0 {
			logger.Log.Debugf("[%v] recv EOF, ready to exit", qout)
			break
		}

		if err := mb.pubsubRedis.RPush(qout, str).Err(); err != nil {
			logger.Log.Errorf("[%v] rc.RPush err:%v, value:%v", qout, err, str)
		} else {
			total := atomic.AddUint32(&mb.msgCount, 1)
			logger.Log.Debugf("[%v] RPush:%v, total:%v", qout, str, total)
		}
	}

	logger.Log.Infof("%d goroutine write exit, qout:[%v]", mb.sid, qout)

	go mb.msgPublish()
}

func (mb *MessageBus) readData() {
	qin := fmt.Sprintf("mq_%v", mb.sid)

	for {
		val, err := mb.pubsubRedis.BLPop(time.Second * time.Duration(1), qin).Result()
		if err != nil {
			if err == redis.Nil {
				continue;
			}
		}

		var msg RouteMessage
		err = json.Unmarshal([]byte(val[1]), &msg)
		if err != nil {
			logger.Log.Errorf("json.Unmarshal err:%v, val:%v", err, val)
		} else {
			mb.readChan <- &msg
		}
	}
}

func (mb *MessageBus) handleMessage() {
	defer(func(){
		close(mb.readChan)
	})()

	for {
		select {
		case msg := <- mb.readChan:
			if msg.Dir == "g2c" {
				mb.SendMsgToClient(gameConst.PushCmd_DDZ, msg.Uid, msg.Data)
			}

		case <- mb.stopChan:
			return
		}
	}
}

func (mb *MessageBus) SendMsgToClient(route string, uid string, data []byte) {
	uids := [] string { uid }
	pitaya.SendPushToUsers(route, data, uids, "connector")
}

func (mb *MessageBus) SendMsgToGameX(uid string, data []byte, dir string) error {
	msg := RouteMessage{
		Sid: mb.sid,
		Uid: uid,
		Dir: dir,
		Data: data,
		GType: 13,
		Gid:30011,
	}

	data, err := json.Marshal(&msg)
	if err != nil {
		logger.Log.Errorf("json.Marshal failed, msg:%v, err:%v", msg, err)
		return err
	}

	mb.writeChan <- data

	return err
}
