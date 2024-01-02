package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/handler"
)

func main() {
	configPath := "etc/msg_api_server.yaml"
	config := &conf.Config{}
	if err := conf.LoadConfig(configPath, config); err != nil {
		panic(err)
	}

	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterMsgApiHandlers(appCtx)

	appCtx.StartServe()
}
