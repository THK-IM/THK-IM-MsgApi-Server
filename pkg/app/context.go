package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-msgapi-server/pkg/loader"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

type Context struct {
	*server.Context
}

func (c *Context) SessionModel() model.SessionModel {
	return c.Context.ModelMap["session"].(model.SessionModel)
}

func (c *Context) SessionMessageModel() model.SessionMessageModel {
	return c.Context.ModelMap["session_message"].(model.SessionMessageModel)
}

func (c *Context) SessionUserModel() model.SessionUserModel {
	return c.Context.ModelMap["session_user"].(model.SessionUserModel)
}

func (c *Context) UserMessageModel() model.UserMessageModel {
	return c.Context.ModelMap["user_message"].(model.UserMessageModel)
}

func (c *Context) UserSessionModel() model.UserSessionModel {
	return c.Context.ModelMap["user_session"].(model.UserSessionModel)
}

func (c *Context) ObjectModel() model.ObjectModel {
	return c.Context.ModelMap["object"].(model.ObjectModel)
}

func (c *Context) SessionObjectModel() model.SessionObjectModel {
	return c.Context.ModelMap["session_object"].(model.SessionObjectModel)
}

func (c *Context) LoginApi() userSdk.LoginApi {
	return c.Context.SdkMap["login_api"].(userSdk.LoginApi)
}

func (c *Context) UserApi() userSdk.UserApi {
	return c.Context.SdkMap["user_api"].(userSdk.UserApi)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	c.Context.SdkMap = loader.LoadSdks(c.Config().Sdks, c.Logger())
	c.Context.ModelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
