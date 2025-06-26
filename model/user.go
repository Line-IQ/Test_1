package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`                      // Mongo 的主键
	UID       string             `bson:"uid" json:"uid"`                               // 自定义唯一标识符
	Username  string             `bson:"username" json:"username"`                     // 用户名
	Password  string             `bson:"password,omitempty" json:"password,omitempty"` // 加密密码
	Email     string             `bson:"email" json:"email"`                           // 邮箱
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`                 // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`                 // 更新时间
}
