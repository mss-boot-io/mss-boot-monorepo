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

var endpointResolverFunc = func(region string, _ s3.EndpointResolverOptions) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           fmt.Sprintf("https://oos-%s.ctyunapi.cn", region),
		SigningRegion: region,
		SigningMethod: "v4",
	}, nil
}

type OOS struct {
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	client          *s3.Client
}

// Init init
func (o *OOS) Init() {
	o.client = s3.New(s3.Options{
		Region: o.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     o.AccessKeyID,
				SecretAccessKey: o.SecretAccessKey,
			}, nil
		}),
		EndpointResolver: s3.EndpointResolverFunc(endpointResolverFunc),
	})
}

// GetClient get client
func (o *OOS) GetClient() *s3.Client {
	return o.client
}
