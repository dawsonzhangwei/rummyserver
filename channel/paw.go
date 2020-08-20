package channel

import (
	"github.com/go-redis/redis"
	"github.com/topfreegames/pitaya/logger"

	"rummy/base"
	"rummy/gameConst"
)

type Paw struct {
	redisAddr string
	redisDbIndex int

	rc *redis.Client	
}

func (p *Paw) configure() error {
	p.redisAddr = mb.config.GetString("rummy.redis.addr")
	if mb.redisAddr == "" {
		return errors.New("paw channel redisAddr is nil")
	}

	p.redisDbIndex = mb.config.GetInt("rummy.redis.db")
	
	logger.Log.Infof("[Channel] paw configure, redisAddr:%v, db:%v", p.redisAddr, p.redisDbIndex)
	return nil
}

func (p *Paw) Init() error {
	var rdb = redis.NewClient(&redis.Options{
		Addr: p.RedisAddr,
		Password: "",
		DB: p.RedisDbIndex,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		logger.Log.Errorf("redis Ping failed, err:%v", err)
	} else {
		rc = rdb
		logger.Log.Infof("paw redis conn succeed.")
	}
}

func (p *Paw) Login(loginParam map[string] string) (player base.Player, error) {
	key := fmt.Sprintf("%s%v", gameConst.RedisKey_Player_Prefix, loginParam["uid"])
	res, err := p.rc.HGetAll(key).Result()
	if err != nil {
		logger.Log.Errorf("paw login error:%v, key:%v", err, key)
		return nil, error.New("paw redis get failed")
	}

	player = &base.Player{

	}

	return player, nil
}