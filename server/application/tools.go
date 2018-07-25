package application

import (
	"github.com/boomfunc/base/conf"
)

func New(router *conf.Router) Interface {
	return &Application{
		router: router,
		packer: new(JSON),
	}
}
