package main

import (
	"fmt"
	"flag"
	"strings"

	"github.com/spf13/viper"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/acceptor"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/serialize/json"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/logger"

	"rummy/services"
	"rummy/modules"
)

func configureBackend() {
	//
	// module
	//

	// MessageBus
	busConfig := viper.New()
	busConfig.SetDefault("rummy.redis.addr", "20.10.1.82:6379")
	busConfig.SetDefault("rummy.redis.db", 3)
	msgBus, err := modules.NewMessageBus(config.NewConfig(busConfig))
	if err != nil {
		logger.Log.Errorf("NewMessageBus failed, err:%v", err)
	} else {
		pitaya.RegisterModule(msgBus, "messageBus")
	}

	// DataCache
	cache := modules.NewDataCache()
	pitaya.RegisterModule(cache, "cache")

	// Db
	dbConfig := viper.New()
	dbConfig.SetDefault("rummy.redis.addr", "20.10.1.82:6379")
	dbConfig.SetDefault("rummy.redis.db", 3)
	db, err := modules.NewDb(config.NewConfig(busConfig))
	if err != nil {
		logger.Log.Errorf("NewDb failed, err:%v", err)
	} else {
		pitaya.RegisterModule(db, "db")
	}

	//
	// service 
	//

	auth := services.NewAuth()
	pitaya.Register(auth,
		component.WithName("auth"),
		component.WithNameFunc(strings.ToLower))

	room := services.NewRoom()
	pitaya.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower))

	router := services.NewRouter()
	pitaya.Register(router,
		component.WithName("router"),
		component.WithNameFunc(strings.ToLower))
}

func configureFrontend(port int) {
	t := acceptor.NewTCPAcceptor(fmt.Sprintf(":%v", port))
	pitaya.AddAcceptor(t)
}

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svrType := flag.String("type", "connector", "the server type")
	isFrontend := flag.Bool("frontend", true, "is server is frontend")

	flag.Parse()

	defer pitaya.Shutdown()

	s := json.NewSerializer()
	pitaya.SetSerializer(s)

	if *isFrontend {
		configureFrontend(*port)
	//} else {
		configureBackend()
	}

	conf := configApp()
	pitaya.Configure(*isFrontend, *svrType, pitaya.Cluster, map[string] string{}, conf)
	pitaya.Start()
}

func configApp() *viper.Viper {
	conf := viper.New()
	conf.SetEnvPrefix("rummy")
	conf.SetDefault("pitaya.buffer.handler.localprocess", 15)
	conf.SetDefault("pitaya.heartbeat.interval", "15s")
	conf.SetDefault("pitaya.buffer.agent.messages", 32)
	conf.SetDefault("pitaya.handler.messages.compression", false)

	conf.SetDefault("pitaya.cluster.rpc.client.nats.connect", "nats://20.10.1.63:4222")
	conf.SetDefault("pitaya.cluster.rpc.server.nats.connect", "nats://20.10.1.63:4222")
	conf.SetDefault("pitaya.cluster.sd.etcd.endpoints", "20.10.1.63:2379")

	defaultMap := map[string] interface{} {
		"custom.redis.pre.url":                "redis://localhost:9010",
		"custom.redis.pre.connectionTimeout":  10,
		"custom.redis.post.url":               "redis://localhost:9010",
		"custom.redis.post.connectionTimeout": 10,
	}

	for param := range defaultMap {
		conf.SetDefault(param, defaultMap[param])
	}
	
	return conf
}