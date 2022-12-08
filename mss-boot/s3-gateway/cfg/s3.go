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
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Provider string

const (
	S3   S3Provider = "s3"
	OSS  S3Provider = "oss"
	OOS  S3Provider = "oos"
	KODO S3Provider = "kodo"
)

var URLTemplate = map[S3Provider]string{
	S3:   "https://%s.amazonaws.com",   //aws s3
	OSS:  "https://%s.aliyuncs.com",    //aliyun oss
	OOS:  "https://oos-%s.ctyunapi.cn", //ctyun oos
	KODO: "https://%s.qiniucs.com",     //qiniu kodo
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

type S3Config struct {
	Provider        S3Provider `yaml:"provider"`
	SigningMethod   string     `yaml:"signingMethod"`
	Region          string     `yaml:"region"`
	Bucket          string     `yaml:"bucket"`
	AccessKeyID     string     `yaml:"accessKeyID"`
	SecretAccessKey string     `yaml:"secretAccessKey"`
	client          *s3.Client
}

// Init init
func (o *S3Config) Init() {
	urlTemplate, exist := URLTemplate[o.Provider]
	if o.SigningMethod == "" {
		o.SigningMethod = "v4"
	}
	if !exist || urlTemplate == "" {
		log.Fatalf("s3 provider %s not support", o.Provider)
	}
	o.client = s3.New(s3.Options{
		Region: o.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     o.AccessKeyID,
				SecretAccessKey: o.SecretAccessKey,
			}, nil
		}),
		EndpointResolver: endpointResolverFunc(urlTemplate, o.SigningMethod),
	})
}

// GetClient get client
func (o *S3Config) GetClient() *s3.Client {
	return o.client
}
