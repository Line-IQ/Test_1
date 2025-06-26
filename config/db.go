package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Mongo_Client 对象, 全局变量
var Mongo_Client *mongo.Client

// DB 用于保存指定数据库的引用, 例如: game_dev
var DB *mongo.Database

// InitMongoDB 初始化 MongoDB 数据库连接
func InitMongoDB() {
	// 设置 MongoDB 连接 URI
	uri := "mongodb://localhost:27017" // 可改成从环境变量加载(不懂)

	// 创建连接选项 (可 设置超时.身份认证等)
	clientOptions := options.Client().ApplyURI(uri)

	// 使用上下文控制连接超时, 最多 10 秒
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB 连接失败: ", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Ping 失败: ", err)
	}

	fmt.Println("MongoDB 连接成功")
	Mongo_Client = client

	// 指定使用哪个数据库 (自动创建)
	DB = client.Database("game_dev")
}
