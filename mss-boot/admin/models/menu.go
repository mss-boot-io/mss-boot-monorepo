/*
 * @Author: lwnmengjing
 * @Date: 2022/3/14 9:24
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/14 9:24
 */

package models

import (
	"context"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/enum"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mongodb.AppendTable(&Menu{})
}

// Menu <no value>
type Menu struct {
	mgm.DefaultModel `bson:",inline"`
	TenantID         string      `bson:"tenantID" json:"tenantID"`
	Name             string      `bson:"name" json:"name"`
	Icon             string      `bson:"icon" json:"icon"`
	Path             string      `bson:"path" json:"path"`
	Access           string      `bson:"access" json:"access"`
	HideInMenu       bool        `bson:"hideInMenu" json:"hideInMenu"`
	Status           enum.Status `bson:"status" json:"status"`
	Routes           []Menu      `bson:"routes" json:"routes"`
	ParentKeys       []string    `bson:"parentKeys" json:"parentKeys"`
	Redirect         string      `bson:"redirect" json:"redirect"`
	Layout           bool        `bson:"layout" json:"layout"`
	Component        string      `bson:"component" json:"component"`
}

func (Menu) TableName() string {
	return "menu"
}

func (e *Menu) Make() {
	ops := options.Index()
	ops.SetName("path")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{"tenantID", bsonx.Int32(1)},
			{"path", bsonx.Int32(1)},
		},
		Options: ops,
	})
	if err != nil {
		log.Fatal(err)
	}
	//初始化数据
	count, err := e.C().CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		now := time.Now()
		e.C().InsertMany(context.TODO(), []interface{}{
			&Menu{
				Path:   "/user",
				Layout: false,
				Routes: []Menu{
					{
						Path: "/user",
						Routes: []Menu{
							{
								Name:      "login",
								Path:      "/user/login",
								Component: "./user/Login",
							},
						},
					},
					{
						Component: "./404",
					},
				},
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			&Menu{
				Name:   "welcome",
				Path:   "/welcome",
				Status: enum.Enabled,
				Layout: true,
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			&Menu{
				Name:   "menu",
				Path:   "/menu/list",
				Status: enum.Enabled,
				Layout: true,
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				Routes: []Menu{
					{
						Name:       "",
						Path:       "/menu/list",
						HideInMenu: true,
					},
					{
						Name:       "",
						Path:       "/menu/control/new",
						HideInMenu: true,
					},
					{
						Name:       "",
						Path:       "/menu/:id",
						HideInMenu: true,
					},
					{
						Name:       "",
						Path:       "/menu/control/:id",
						HideInMenu: true,
					},
				},
			},
			&Menu{
				Path:     "/",
				Redirect: "/welcome",
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			&Menu{
				Component: "./404",
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			&Menu{
				Name:   "generate",
				Path:   "/generate",
				Layout: true,
				DefaultModel: mgm.DefaultModel{
					DateFields: mgm.DateFields{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
		})
	}
}

func (e *Menu) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}
