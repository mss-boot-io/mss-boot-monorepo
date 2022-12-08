/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 14:59:13
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 14:59:13
 */

package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/s3-gateway/cfg"
)

//func init() {
//	e := &OOS{}
//	//e.Auth = false
//	response.AppendController(e)
//}

type OOS struct {
	//curd.DefaultController
}

func (e *OOS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	key = key[1:]
	if key == "" {
		key = "index.html"
	}
	client := cfg.Cfg.Provider.GetClient()

	switch r.Method {
	case http.MethodGet:
		// get object
		output, err := client.GetObject(r.Context(), &s3.GetObjectInput{
			Bucket: &cfg.Cfg.Provider.Bucket,
			Key:    &key,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("get object failed, err:%s", err.Error()))
			return
		}
		defer output.Body.Close()
		w.Header().Set("Content-Type", *output.ContentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", output.ContentLength))

		_, err = io.Copy(w, output.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("copy object failed, err:%s", err.Error()))
			return
		}
	case http.MethodPut, http.MethodPost:
		// put object
		defer r.Body.Close()
		h := make([]byte, 512)
		_, err := io.ReadAtLeast(r.Body, h, 512)
		if err != nil {
			if err == io.EOF {
				output, err := s3.NewPresignClient(client).PresignPutObject(r.Context(), &s3.PutObjectInput{
					Bucket: &cfg.Cfg.Provider.Bucket,
					Key:    &key,
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintln(w, fmt.Sprintf("get put object sign url failed, err:%s", err.Error()))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(output.URL))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("read body(min 512) failed, err:%s", err.Error()))
			return
		}
		contentType := http.DetectContentType(h)
		buffer := bytes.NewBuffer(h)
		_, err = buffer.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("read body failed, err:%s", err.Error()))
			return
		}
		_, err = client.PutObject(r.Context(), &s3.PutObjectInput{
			Bucket:      &cfg.Cfg.Provider.Bucket,
			Key:         &key,
			Body:        buffer,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("read body failed, err:%s", err.Error()))
			return
		}
	case http.MethodDelete:
		// delete object
		_, err := client.DeleteObject(r.Context(), &s3.DeleteObjectInput{
			Bucket: &cfg.Cfg.Provider.Bucket,
			Key:    &key,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("delete object failed, err:%s", err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
