/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/26 00:15:32
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/26 00:15:32
 */

package cfg

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ProviderType string

const (
	S3   ProviderType = "s3"   //aws s3
	OSS  ProviderType = "oss"  //aliyun oss
	OOS  ProviderType = "oos"  //ctyun oos
	KODO ProviderType = "kodo" //qiniu kodo
	COS  ProviderType = "cos"  //tencent cos
)

var URLTemplate = map[ProviderType]string{
	OSS:  "https://%s.aliyuncs.com",
	OOS:  "https://oos-%s.ctyunapi.cn",
	KODO: "https://s3-%s.qiniucs.com",
	COS:  "https://cos.%s.myqcloud.com",
}

var endpointResolverFunc = func(urlTemplate, signingMethod string) s3.EndpointResolverFunc {
	return func(region string, _ s3.EndpointResolverOptions) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           fmt.Sprintf(urlTemplate, region),
			SigningRegion: region,
			SigningMethod: signingMethod,
		}, nil
	}
}

type ProviderConfig struct {
	Type            ProviderType `yaml:"type"`
	SigningMethod   string       `yaml:"signingMethod"`
	Region          string       `yaml:"region"`
	Bucket          string       `yaml:"bucket"`
	AccessKeyID     string       `yaml:"accessKeyID"`
	SecretAccessKey string       `yaml:"secretAccessKey"`
	client          *s3.Client
}

// Init init
func (o *ProviderConfig) Init() {
	var endpointResolver s3.EndpointResolver
	if o.Type != S3 {
		if urlTemplate, exist := URLTemplate[o.Type]; exist && urlTemplate != "" {
			endpointResolver = endpointResolverFunc(urlTemplate, o.SigningMethod)
		}
	}
	o.client = s3.New(s3.Options{
		Region: o.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     o.AccessKeyID,
				SecretAccessKey: o.SecretAccessKey,
			}, nil
		}),
		EndpointResolver: endpointResolver,
	})
}

// GetClient get client
func (o *ProviderConfig) GetClient() *s3.Client {
	return o.client
}
