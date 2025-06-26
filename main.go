package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// 加载前端模板
	r.LoadHTMLGlob("frontend/*.html")
	r.Static("/static", "./frontend")

	// 注册路由
	router.InitRoutes(r)

	// 启动服务
	r.Run(":8080")
}
