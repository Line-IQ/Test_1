package config

import (
	"log"
	"os"
)

var LogFile *os.File

func InitLogger() {
	file, err := os.OpenFile("../logs/server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	LogFile = file
	log.SetOutput(file) // 将默认日志输出写入文件
}
