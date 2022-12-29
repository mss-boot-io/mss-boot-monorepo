/*
 * @Author: lwnmengjing
 * @Date: 2022/3/14 9:20
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/14 9:20
 */

package models

import (
	"context"
	"github.com/kamva/mgm/v3"
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/enum"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func init() {
	mongodb.AppendTable(&Role{})
}

type Role struct {
	mgm.DefaultModel `bson:",inline"`
	TenantID         string      `bson:"tenantID" json:"tenantID"`
	Name             string      `bson:"name" json:"name"`
	Status           enum.Status `bson:"status" json:"status"`
	Metadata         interface{} `bson:"metadata" json:"metadata"`
}

func (Role) TableName() string {
	return "role"
}

func (e *Role) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}

func (e *Role) Make() {
	ops := options.Index()
	ops.SetName("name")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{"name", bsonx.Int32(1)},
				{"tenantID", bsonx.Int32(1)},
			},
			Options: ops,
		})
	if err != nil {
		log.Fatal(err)
	}
}
