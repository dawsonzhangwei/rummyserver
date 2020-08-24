package modules

import (
	"fmt"
	"errors"

	"github.com/go-redis/redis"
	"github.com/topfreegames/pitaya/logger"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/modules"

	"rummy/base"
	"rummy/gameConst"
)

var DB *Db

type Db struct {
	modules.Base

	config *config.Config
	redisAddr string
	redisDbIndex int

	rc *redis.Client	
}

func NewDb(config *config.Config) (*Db, error) {
	db := &Db{
		config: config,
	}
	if err := db.configure(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Db) configure() error {
	db.redisAddr = db.config.GetString("rummy.redis.addr")
	if db.redisAddr == "" {
		return errors.New("paw channel redisAddr is nil")
	}

	db.redisDbIndex = db.config.GetInt("rummy.redis.db")
	
	logger.Log.Infof("[Channel] paw configure, redisAddr:%v, db:%v", db.redisAddr, db.redisDbIndex)
	return nil
}

func (db *Db) Init() error {
	var rdb = redis.NewClient(&redis.Options{
		Addr: db.redisAddr,
		Password: "",
		DB: db.redisDbIndex,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		logger.Log.Errorf("redis Ping failed, err:%v", err)
	} else {
		db.rc = rdb
		DB = db
		logger.Log.Infof("paw redis conn succeed.")
	}

	return err
}

func (db *Db) LoadPlayer(player *base.Player) error {
	key := fmt.Sprintf("%s%v", gameConst.RedisKey_Player_Prefix, player.UID)
	res, err := db.rc.HGetAll(key).Result()
	if err != nil {
		logger.Log.Errorf("paw login error:%v, key:%v", err, key)
		return errors.New("paw redis get failed")
	}

	player.SetAttr(res)

	return nil
}