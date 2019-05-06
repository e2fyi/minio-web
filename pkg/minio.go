// Package pkg provides utils to interact with S3 compatible backend with 
// minio.
package pkg

import (
	"errors"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	"io"
	"strings"
	"fmt"
	"time"
)

// MinioHelper provides the interface to the S3 compatible backend.
type MinioHelper struct {
	Client     	*minio.Client
	BucketName 	string
	maxAttempts int
}

// ConfigMinioHelper creates a internal helper to interact with the S3
// compatible backend.
func (app *App) ConfigMinioHelper(config MinioConfig, bucketName string) *App {
	// create client
	client, err := NewMinioClient(config)
	if err != nil {
		app.sugar.Fatal(err)
	}

	helper := MinioHelper{Client: client, BucketName: bucketName, maxAttempts: 5}

	// try connection for 5 times before throwing error
	app.sugar.Infof("connecting to %s", config.Endpoint)
	for i := 1;  i<=5; i++ {
		msg, err := helper.testConnection()
		if err == nil {
			app.sugar.Info(msg)
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		app.sugar.Fatal(err)
	}

	// create helper
	app.helper = helper
	app.handler.GetObject = app.helper.GetObject
	app.handler.StatObject = app.helper.StatObject
	return app
}

// NewMinioClient creates a new minio client. If accesskey and secretkey are 
// provided, they are used to create the client. Otherwise, AWS EC2 IAM role
// credentials will be used instead. 
func NewMinioClient(config MinioConfig) (*minio.Client, error) {
	if config.AccessKey == "" && config.SecretKey == "" {
		return minio.NewWithCredentials(config.Endpoint, 
										credentials.NewIAM(""), 
										config.Secure, 
										config.Region)
	}
	if config.Region == "" {
		return minio.NewWithRegion(config.Endpoint, 
								   config.AccessKey, 
								   config.SecretKey, 
								   config.Secure, 
								   config.Region)
	}
	return minio.New(config.Endpoint, 
					 config.AccessKey, 
				 	 config.SecretKey, 
				 	 config.Secure)
}

// minioObjectInfoToResourceInfo converts a minio ObjectInfo to ResourceInfo.
func minioObjectInfoToResourceInfo(info minio.ObjectInfo) ResourceInfo {
	return ResourceInfo{
		Size:         info.Size,
		ETag:         info.ETag,
		LastModified: info.LastModified,
		ContentType:  info.ContentType}
}

// test connection to backend
func (h *MinioHelper) testConnection() (string, error) {
	// test connection
	if h.BucketName != "" {
		exist, err := h.Client.BucketExists(h.BucketName)
		return fmt.Sprintf("bucket %s: %t", h.BucketName, exist), err
	} 

	buckets, err := h.Client.ListBuckets()
	return fmt.Sprintf("# buckets found: %d", len(buckets)), err
}

// GetBucketNameAndPrefix infers the bucket name and prefix from the url.
func (h *MinioHelper) GetBucketNameAndPrefix(url string) (bucketname string, prefix string) {
	if url[0] == '/' {
		url = url[1:]
	}
	if h.BucketName == "" {
		parts := strings.Split(url, "/")
		switch numParts := len(parts); {
		case numParts > 1:
			return parts[0], strings.Join(parts[1:], "/")
		case numParts == 1:
			return parts[0], ""
		default:
			return "", ""
		}
	}
	return h.BucketName, url
}

// GetObject retrieves the metadata and data from the S3 compatible backend.
func (h *MinioHelper) GetObject(url string) (Resource, error) {
	bucketName, prefix := h.GetBucketNameAndPrefix(url)
	if bucketName == "" {
		return Resource{}, errors.New("No bucket provided")
	}
	info, err := h.Client.StatObject(bucketName, prefix, minio.StatObjectOptions{})
	if err != nil {
		return Resource{}, err
	}
	obj, err := h.Client.GetObject(bucketName, prefix, minio.GetObjectOptions{})
	if err != nil {
		return Resource{}, err
	}
	return Resource{Data: obj, Info: minioObjectInfoToResourceInfo(info)}, nil
}

// StatObject retrieves the metadata (only) from a S3 compatible backend.
func (h *MinioHelper) StatObject(url string) (Resource, error) {
	bucketName, prefix := h.GetBucketNameAndPrefix(url)
	if bucketName == "" {
		return Resource{}, errors.New("No bucket provided")
	}
	info, err := h.Client.StatObject(bucketName, prefix, minio.StatObjectOptions{})
	if err != nil {
		return Resource{}, err
	}
	return Resource{Data: io.Reader(nil), Info: minioObjectInfoToResourceInfo(info)}, nil
}
