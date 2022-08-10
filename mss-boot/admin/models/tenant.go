/*
 * @Author: lwnmengjing
 * @Date: 2022/3/14 9:24
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/14 9:24
 */

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mss-boot-io/mss-boot/client"
	pb "github.com/mss-boot-io/mss-boot/proto/store/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/cfg"
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/enum"
	"github.com/mss-boot-io/mss-boot/pkg/store"
)

func init() {
	store.DefaultOAuth2Store = &Tenant{}
	mongodb.AppendTable(&Tenant{})
}

// Tenant 租户
type Tenant struct {
	ID          string         `json:"id" bson:"_id"`
	Name        string         `json:"name" bson:"name"`
	Contact     string         `json:"contact" bson:"contact"`
	System      bool           `json:"system" bson:"system"`
	Status      enum.Status    `json:"status" bson:"status"`
	Description string         `json:"description" bson:"description"`
	Domains     []string       `json:"domains" bson:"domains"`
	Client      config.OAuth2  `json:"client" bson:"client"`
	Metadata    interface{}    `json:"metadata" bson:"metadata"`
	ExpiredAt   time.Time      `json:"expiredAt" bson:"expiredAt" binding:"required"`
	CreatedAt   time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt" bson:"updatedAt"`
	storeClient pb.StoreClient `bson:"-"`
}
type OnlyClient struct {
	Client config.OAuth2 `json:"client" bson:"client"`
}

func NewTenant(storeclient pb.StoreClient) *Tenant {
	return &Tenant{storeClient: storeclient}
}

func (Tenant) TableName() string {
	return "tenant"
}

func (e *Tenant) Make() {
	ops := options.Index()
	ops.SetName("name")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{{"name", bsonx.Int32(1)}},
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
		e.C().InsertOne(context.TODO(), Tenant{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "mss-boot-io",
			Contact:     "mss-boot-io",
			System:      true,
			Status:      enum.Enabled,
			Description: "mss-boot-io",
			Domains:     []string{"localhost:8080", "localhost:8081"},
			Client: config.OAuth2{
				Issuer:       cfg.Cfg.OAuth2.Issuer,
				ClientID:     cfg.Cfg.OAuth2.ClientID,
				ClientSecret: cfg.Cfg.OAuth2.ClientSecret,
				Scopes:       cfg.Cfg.OAuth2.Scopes,
				RedirectURL:  cfg.Cfg.OAuth2.RedirectURL,
			},
			ExpiredAt: now.Add(time.Hour * 24 * 365 * 100),
			CreatedAt: now,
			UpdatedAt: now,
		})
		return
	}
}

func (e *Tenant) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}

// GetClientByDomain 获取租户的client
func (e *Tenant) GetClientByDomain(c context.Context, domain string) (store.OAuth2Configure, error) {
	fmt.Println(domain)
	if e.storeClient == nil {
		e.storeClient = client.Store().GetClient()
	}
	resp, err := e.storeClient.Get(c, &pb.GetReq{Key: fmt.Sprintf("tenant_%s", domain)})
	if err != nil {
		return nil, err
	}
	data := &OnlyClient{}
	fmt.Println(resp.Value)
	err = json.Unmarshal([]byte(resp.Value), data)
	if err != nil {
		return nil, err
	}
	return &data.Client, nil
}

func (e *Tenant) InitStore() {
	//todo 初始化所有租户到store
	c := client.Store().GetClient()

	cur, err := e.C().Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var tenant Tenant
		err := cur.Decode(&tenant)
		if err != nil {
			log.Fatal(err)
		}
		rb, err := json.Marshal(tenant)
		if err != nil {
			log.Fatal(err)
		}
		for _, domain := range tenant.Domains {
			_, err = c.Set(context.TODO(), &pb.SetReq{
				Key:   fmt.Sprintf("tenant_%s", domain),
				Value: string(rb),
			})
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
