/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/19 15:29:18
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/19 15:29:18
 */

package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func init() {
	mongodb.AppendTable(&Github{})
}

type Github struct {
	ID        string    `bson:"_id" json:"id"`
	TenantID  string    `bson:"tenantID" json:"tenantID"`
	Email     string    `bson:"email" json:"email"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (Github) TableName() string {
	return "github"
}

func (e *Github) Make() {
	ops := options.Index()
	ops.SetName("user")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{"tenantID", bsonx.Int32(1)},
			{"userID", bsonx.Int32(1)},
		},
		Options: ops,
	})
	if err != nil {
		log.Fatalf("create index error: %s", err.Error())
	}
}

func (e *Github) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}

// GetMyGithubConfig 获取当前用户的github配置
func GetMyGithubConfig(c context.Context, tenantID, email string) (*Github, error) {
	g := &Github{}
	err := g.C().FindOne(c, bson.M{"tenantID": tenantID, "email": email}).Decode(g)
	if err != nil {
		return nil, err
	}
	return g, nil
}
