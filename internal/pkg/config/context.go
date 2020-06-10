package config

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Context struct {
	Ctx context.Context
	Logger *logrus.Entry
	Config *ServiceConfig
}

func NewContext(parent context.Context, svcConfig *ServiceConfig) *Context {
	return &Context{
		Ctx: parent,
		Logger: svcConfig.Logger,
		Config: svcConfig,
	}
}

func (c *Context) Deadline() (time.Time, bool) {
	return c.Ctx.Deadline()
}
