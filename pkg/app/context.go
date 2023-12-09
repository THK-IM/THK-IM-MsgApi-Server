package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/loader"
	"github.com/thk-im/thk-im-base-server/server"
)

type Context struct {
	*server.Context
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
