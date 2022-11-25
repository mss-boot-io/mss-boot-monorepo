/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/19 03:50:28
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/19 03:50:28
 */

package form

import "time"

type SettingsUpdateReq struct {
	ID        string    `bson:"_id" json:"-"`
	TenantID  string    `bson:"tenantID" json:"-"`
	UserID    string    `bson:"userID" json:"-"`
	Settings  []KV      `json:"settings"`
	CreatedAt time.Time `json:"-" bson:"createdAt"`
	UpdatedAt time.Time `json:"-" bson:"updatedAt"`
}
