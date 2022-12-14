/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/13 23:16:08
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/13 23:16:08
 */

package cfg

import (
	"embed"
	"time"
)

var (
	//go:embed *.yml
	FS  embed.FS
	Cfg Config
)

type Config struct {
	Client Client `yaml:"client"`
}

type Client struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}
