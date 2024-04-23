package minio

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

// Helper provides the interface to the S3 compatible backend.
type Helper struct {
	Client      *minio.Client
	BucketName  string
	Prefix      string
	maxAttempts int
}

// Config is used to create a minio client.
type Config struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	Secure    bool   `json:"secure"`
	Region    string `json:"region"`
}

// NewMinioHelperWithBucket creates a new minio.Helper object.
func NewMinioHelperWithBucket(config Config, bucketName string, prefix string, maxAttempts int) (Helper, error) {
	// create client
	client, err := NewMinioClient(config)
	if err != nil {
		return Helper{}, err
	}
	return Helper{Client: client, BucketName: bucketName, Prefix: prefix, maxAttempts: maxAttempts}, nil
}

// NewMinioClient creates a new minio client. If accesskey and secretkey are
// provided, they are used to create the client. Otherwise, AWS EC2 IAM role
// credentials will be used instead.
func NewMinioClient(config Config) (*minio.Client, error) {
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
		Key:          info.Key,
		Size:         info.Size,
		ETag:         info.ETag,
		LastModified: info.LastModified,
		ContentType:  info.ContentType}
}

// TestConnection test connection to backend
func (h *Helper) TestConnection() (string, error) {
	// test connection
	if h.BucketName != "" {
		exist, err := h.Client.BucketExists(h.BucketName)
		return fmt.Sprintf("bucket %s: %t", h.BucketName, exist), err
	}

	buckets, err := h.Client.ListBuckets()
	return fmt.Sprintf("# buckets found: %d", len(buckets)), err
}

// GetBucketNameAndPrefix infers the bucket name and prefix from the url.
func (h *Helper) GetBucketNameAndPrefix(url string) (bucketname string, prefix string) {
	// remove absolute reference
	if len(url) > 0 && url[0] == '/' {
		url = url[1:]
	}
	// infer bucketname
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
func (h *Helper) GetObject(url string) (Resource, error) {
	bucketName, prefix := h.GetBucketNameAndPrefix(url)
	if bucketName == "" {
		return Resource{Msg: fmt.Sprintf("GET[%s]: Bucket name not known", url)}, errors.New("Bucket name not known")
	}
	// add user provided prefix if any
	prefix = h.Prefix + prefix
	// get obj info
	info, err := h.Client.StatObject(bucketName, prefix, minio.StatObjectOptions{})
	if err != nil {
		return Resource{Msg: fmt.Sprintf("GET[%s] -> GetObject[%s/%s]: %v", url, bucketName, prefix, err)}, err
	}
	// get obj
	obj, err := h.Client.GetObject(bucketName, prefix, minio.GetObjectOptions{})
	if err != nil {
		return Resource{}, err
	}
	return Resource{
		Data: obj,
		Info: minioObjectInfoToResourceInfo(info),
		Msg:  fmt.Sprintf("GET[%s] -> GetObject[%s/%s] ok", url, bucketName, prefix)}, nil
}

// StatObject retrieves the metadata (only) from a S3 compatible backend.
func (h *Helper) StatObject(url string) (Resource, error) {
	bucketName, prefix := h.GetBucketNameAndPrefix(url)
	if bucketName == "" {
		return Resource{Msg: fmt.Sprintf("StatObject[%s]: Bucket name not known", url)}, errors.New("Bucket name not known")
	}
	// add user provided prefix if any
	prefix = h.Prefix + prefix
	// get obj info
	info, err := h.Client.StatObject(bucketName, prefix, minio.StatObjectOptions{})
	if err != nil {
		return Resource{Msg: fmt.Sprintf("StatObject[%s/%s]: %v", bucketName, prefix, err)}, err
	}
	return Resource{
		Data: io.Reader(nil),
		Info: minioObjectInfoToResourceInfo(info),
		Msg:  fmt.Sprintf("StatObject[%s/%s] ok", bucketName, prefix)}, nil
}
