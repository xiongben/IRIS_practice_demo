package user

import (
	"fmt"
	"github.com/kataras/iris"
)

func UsersRouter() {

}

type UserInfo struct {
	Email     string
	LogoutURL string
}

func GetUserMessage(ctx iris.Context) {
	fmt.Println(ctx.URLParams())
	fmt.Print(ctx.GetContentLength())
	fmt.Println("=======")
}
