/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 15:13:44
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 15:13:44
 */

package form

type HelloworldCallReq struct {
	//your name
	Name string `json:"name" binding:"required"`
}

type HelloworldCallResp struct {
	//send message
	Message string `json:"message"`
}
