package controller

import (
	"GameDev/config"
	"GameDev/model"
	"GameDev/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Register 处理用户注册请求
func Register(c *gin.Context) {
	var req model.User

	// 前端传来的 JSON 绑定给 req
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"参数绑定失败: ": err.Error()})
		return
	}

	// 查询用户名是否存在
	// 定义 数据库table 变量
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"uid": req.Username}).Decode(&user)
	if err == nil {
		// 存在同名用户, 返回冲突
		c.JSON(http.StatusConflict, gin.H{
			"name":  user.Username,
			"error": "用户名已存在",
		})
		return
	} else if err != mongo.ErrNoDocuments {
		// 除 不存在用户名 外的查询失败
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 最后, 插入用户
	now := time.Now()
	req.Password = string(hashedPwd)
	req.UID = utils.GenerateUID()
	req.CreatedAt = now
	req.UpdatedAt = now

	_, err = collection.InsertOne(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户创建失败",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":    req.Username,
		"uid":     req.UID,
		"message": "注册成功",
	})
}

func Login(c *gin.Context) {
	// 登录请求所需的结构体
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户数据
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 判断 用户是否存在
	var user model.User
	// FindOne 查询 是否存在 username 为 req.Username; Decode 将其 解码到 user
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "数据库查询失败",
			"error":   err.Error()})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "密码错误",
			"error":   err.Error(),
		})
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(c, user.UID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "generate Token Failed",
			"error":   err.Error(),
		})
		return
	}

	// 写入登录日志
	utils.LoginLog(user.Username, c.ClientIP())

	// 登录成功
	c.JSON(http.StatusOK, gin.H{
		"message": "登陆成功",
		"token":   token,
	})
}
