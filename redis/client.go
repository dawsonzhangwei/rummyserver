package redis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
	"github.com/topfreegames/pitaya/config"
)

var (
	clients map[string] *redis.Client

	mu sync.Mutex
)

func GetRedis(
	prefix string,
	config *config.Config,
) (*redis.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if clients == nil {
		clients = make(map[string] *redis.Client)
	}
	if _, ok := clients[prefix]; !ok {
		addrKey := fmt.Sprintf("rummy.redis.addr.%v", prefix)
		dbKey := fmt.Sprintf("rummy.redis.db.%v", prefix)

		addr := config.GetString(addrKey)
		if addr == "" {
			return nil, fmt.Errorf("get redis address failed, addKey:%v, prefix:%v", addrKey, prefix)
		}

		db := config.GetInt(dbKey)
		
		var rdb = redis.NewClient(&redis.Options{
			Addr: addr,
			Password: "",
			DB: db,
		})

		_, err := rdb.Ping().Result()
		if err != nil {
			return nil, fmt.Errorf("redis ping failed, prefix:%v, err:%v", prefix, err)
		}

		clients[prefix] = rdb
	}

	return clients[prefix], nil
	
}
