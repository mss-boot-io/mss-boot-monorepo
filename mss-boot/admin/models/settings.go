/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/19 03:11:38
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/19 03:11:38
 */

package models

import (
	"context"
	"time"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Settings struct {
	ID        string    `bson:"_id" json:"id"`
	TenantID  string    `bson:"tenantID" json:"tenantID"`
	UserID    string    `bson:"userID" json:"userID"`
	Items     []KV      `bson:"items" json:"items"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (*Settings) TableName() string {
	return "settings"
}

func (e *Settings) Make() {
	ops := options.Index()
	ops.SetName("key")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{"key", bsonx.Int32(1)},
			{"userID", bsonx.Int32(1)},
			{"tenantID", bsonx.Int32(1)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Settings) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}
