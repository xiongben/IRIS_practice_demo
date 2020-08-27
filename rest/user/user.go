package user

import (
	"fmt"
	"github.com/kataras/iris"
)

type UserInfo struct {
	Email     string
	LogoutURL string
}

func GetUserMessage(ctx iris.Context) {
	fmt.Print(ctx.Params().Get("name"))

}
