package main

import (
	"github.com/spf13/viper"

	"github.com/topfreegames/pitaya"
)

func main() {
	defer pitaya.Shutdown()

	conf := configApp()
	pitaya.Configure(true, "rummy", pitaya.Cluster, map[string] string{}, conf)
	pitaya.Start()
}

func configApp() *viper.Viper {
	conf := viper.New()
	conf.SetEnvPrefix("rummy")
	conf.SetDefault("pitaya.buff.handler.localprocess", 15)
	conf.SetDefault("pitaya.heartbeat.interval", "15s")
	//conf.SetDefault("pitaya.buff.agent.messagesBufferSize")

	return conf
}