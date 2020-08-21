package modules

import (
	"time"
	"github.com/astaxie/beego/cache"

	"github.com/topfreegames/pitaya/modules"
	"github.com/topfreegames/pitaya/logger"

	"rummy/base"
)

const (
	CacheKey_Player string = "player"
)

var Cache *DataCache

type DataCache struct {
	modules.Base

	storage cache.Cache
}

func NewDataCache() *DataCache {
	cache := &DataCache{
	}
	return cache
}


func (d *DataCache) Init() error {
	c, err := cache.NewCache("memory", `{}`)
	if err != nil {
		logger.Log.Errorf("DataCache storage init failed, err:%v", err)
	} else {
		d.storage = c
		Cache = d
	}

	return err
}

func DataKey(prefix string, key string) string {
	return prefix + key
}

func (d *DataCache) Get(prefix string, key string) interface{} {
	k := DataKey(prefix, key)
	return d.storage.Get(k)
}

func (d *DataCache) Put(prefix string, key string, value interface{}, timeout time.Duration) error {
	k := DataKey(prefix, key)
	return d.storage.Put(k, value, timeout)
}

func (d *DataCache) GetPlayer(uid string) *base.Player {
	data := d.Get(CacheKey_Player, uid)
	if data != nil {
		player := data.(*base.Player)
		return player
	}

	return nil
}

func (d *DataCache) AddPlayer(player *base.Player) error {
	return d.Put(CacheKey_Player, player.UID, player, 86400000*time.Second)
}
