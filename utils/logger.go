package utils

import (
	"fmt"
	"os"
	"time"
)

// Loginlog 写入登录日志
func LoginLog(username, ip string) {
	f, err := os.OpenFile("login.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return
	}
	defer f.Close()

	log := fmt.Sprintf("[%s] 用户 %s 登录，IP：%s\n", time.Now().Format(time.RFC3339), username, ip)
	f.WriteString(log)
}
