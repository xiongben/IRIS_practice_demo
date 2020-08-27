package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"iris_demo/rest/user"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	//（可选）添加两个内置处理程序
	//可以从任何与http相关的panics中恢复
	//并将请求记录到终端。
	app.Use(recover.New())
	app.Use(logger.New())
	// 请求方法: GET
	// 资源标识: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome,xiongben!</h1>")
	})
	// 等同于 app.Handle("GET", "/ping", [...])
	// 请求方法: GET
	// 资源标识: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	// 请求方法: GET
	// 资源标识: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	userPart := app.Party("/users", myAuthMiddlewareHandler)

	userPart.Get("/message", user.GetUserMessage)

	app.OnErrorCode(iris.StatusNotFound, func(context iris.Context) {
		context.HTML("<h1>404!not found page!try again!</h1>")
	})
	app.Run(iris.Addr(":8086"), iris.WithoutServerError(iris.ErrServerClosed))
}

func myAuthMiddlewareHandler(ctx iris.Context) {
	fmt.Println("9999")
	ctx.Next()
}
