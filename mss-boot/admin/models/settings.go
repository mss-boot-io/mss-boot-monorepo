/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/29 15:21:40
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/29 15:21:40
 */

package models

import (
	"github.com/kamva/mgm/v3"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mongodb.AppendTable(&Settings{})
}

type Settings struct {
	mgm.DefaultModel `bson:",inline"`
	TenantID         string         `bson:"tenantID" json:"tenantID"`
	UserID           string         `bson:"userID" json:"userID"`
	Metadata         map[string]any `bson:"metadata" json:"metadata"`
}

func (*Settings) TableName() string {
	return "settings"
}

func (e *Settings) Make() {}

func (e *Settings) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}
