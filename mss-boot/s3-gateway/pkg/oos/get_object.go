/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/25 23:32:54
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/25 23:32:54
 */

package oos

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var endpointResolverFunc = func(region string, _ s3.EndpointResolverOptions) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           fmt.Sprintf("https://%s.aliyuncs.com", region),
		SigningRegion: region,
		SigningMethod: "v4",
	}, nil
}

// GetObject gets an object from the bucket.
func GetObject(ctx context.Context, region, accessKeyID, secretAccessKey, bucket, key string) (*s3.GetObjectOutput, error) {
	client := s3.New(s3.Options{
		Region: region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
			}, nil
		}),
		EndpointResolver: s3.EndpointResolverFunc(endpointResolverFunc),
	})
	return client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}
