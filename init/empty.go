// +build create_router

package init

import (
	"github.com/kataras/iris/v12"
)

func NewApp() *iris.Application {
	return iris.New()
}
