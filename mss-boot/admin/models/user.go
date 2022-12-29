/*
 * @Author: lwnmengjing
 * @Date: 2022/3/14 13:25
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/14 13:25
 */

package models

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kamva/mgm/v3"
	"github.com/mss-boot-io/mss-boot/pkg/enum"
	"github.com/mss-boot-io/mss-boot/pkg/middlewares"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/security"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Subject          string      `bson:"subject" json:"subject"`
	TenantID         string      `bson:"tenantID" json:"tenantID"`
	Username         string      `bson:"username" json:"username"`
	Nickname         string      `bson:"nickname" json:"nickname"`
	Avatar           string      `bson:"avatar" json:"avatar"`
	Email            string      `bson:"email" json:"email"`
	Phone            string      `bson:"phone" json:"phone"`
	Status           enum.Status `bson:"status" json:"status"`
	PWD              UserPwd     `bson:"pwd" json:"pwd"`
	Metadata         interface{} `bson:"metadata" json:"metadata"`
}

type UserPwd struct {
	Salt string `bson:"salt" json:"salt"`
	Hash string `bson:"hash" json:"hash"`
}

func (User) TableName() string {
	return "user"
}

func (e *User) Make() {
	ops := options.Index()
	ops.SetName("name")
	ops.SetUnique(true)
	_, err := e.C().Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{"username", bsonx.Int32(1)},
			{"tenantID", bsonx.Int32(1)},
		},
		Options: ops,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (e *User) C() *mongo.Collection {
	return mongodb.DB.Collection(e.TableName())
}

func (e *User) Encrypt(pwd string) (err error) {
	if pwd == "" {
		return errors.New("password can't be empty")
	}
	e.PWD.Salt = security.GenerateRandomKey16()
	e.PWD.Hash, err = security.SetPassword(pwd, e.PWD.Salt)
	return
}

func (e *User) VerifyPassword(pwd string) bool {
	verify, err := security.SetPassword(pwd, e.PWD.Salt)
	if err != nil {
		return false
	}
	return verify == e.PWD.Hash
}

// CreateOrUpdateUser create or update user
func CreateOrUpdateUser(ctx context.Context, domain string, idToken *oidc.IDToken) error {
	tenant := &Tenant{}
	err := tenant.GetTenantByDomain(ctx, domain)
	if err != nil {
		return err
	}
	claims := &middlewares.User{}
	err = idToken.Claims(claims)
	if err != nil {
		return nil
	}
	//todo set usernamePostfix
	usernamePostfix := ""
	user := &User{
		TenantID: tenant.ID.Hex(),
		Subject:  claims.Subject,
		Username: claims.Name + usernamePostfix,
		Nickname: claims.Name + usernamePostfix,
		Avatar:   claims.Audience,
		Email:    claims.Email,
		Status:   enum.Enabled,
		Metadata: *claims,
		DefaultModel: mgm.DefaultModel{
			DateFields: mgm.DateFields{
				CreatedAt: time.Now(),
			},
		},
	}
	var count int64
	count, err = user.C().CountDocuments(ctx, bson.M{"username": user.Username, "tenantID": user.TenantID})
	if err != nil {
		return err
	}
	if count > 0 {
		//user exist
		//user.UpdatedAt = time.Now()
		//// todo update item
		//_, err = user.C().UpdateOne(ctx, bson.M{"username": user.Username, "tenantID": user.TenantID}, bson.M{"$set": user})
		//if err != nil {
		//	return err
		//}
		return nil
	}
	user.UpdatedAt = user.CreatedAt
	user.ID = primitive.NewObjectID()
	_, err = user.C().InsertOne(ctx, user)
	return err
}
