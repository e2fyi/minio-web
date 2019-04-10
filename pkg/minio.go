// Package pkg provides utils to interact with S3 compatible backend with 
// minio.
package pkg

import (
	"errors"
	"github.com/minio/minio-go"
	"io"
	"strings"
)

// MinioHelper provides the interface to the S3 compatible backend.
type MinioHelper struct {
	Client     *minio.Client
	BucketName string
}

// minioObjectInfoToResourceInfo converts a minio ObjectInfo to ResourceInfo.
func minioObjectInfoToResourceInfo(info minio.ObjectInfo) ResourceInfo {
	return ResourceInfo{
		Size:         info.Size,
		ETag:         info.ETag,
		LastModified: info.LastModified,
		ContentType:  info.ContentType}
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
