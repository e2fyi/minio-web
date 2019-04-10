// Package main starts a web server which serve resources in S3 compatible 
// backend with regular HTTP GET.
package main

import (
	"fmt"
	"github.com/e2fyi/minio-web/pkg"
	"github.com/minio/minio-go"
	"log"
	"net/http"
)

func main() {

	// load config
	config := pkg.LoadConfig()

	// Initialize minio client object.
	minioClient, err := minio.New(
		fmt.Sprintf("%s:%d", config.Endpoint, config.Port),
		config.AccessKeyID,
		config.SecretAccessKey,
		config.Ssl)
	if err != nil {
		log.Fatalln(err)
	}

	// minio helper
	helper := pkg.MinioHelper{Client: minioClient, BucketName: config.BucketName}
	log.Printf("serving bucket[%s] at %s:%d", config.BucketName, config.Endpoint, config.Port)

	// insert default index file if not provided (e.g. index.html)
	getObject := pkg.NewIndexHTML(config.DefaultHTML).GetIndexHTML(helper.GetObject)
	log.Printf("default index: %s", config.DefaultHTML)

	// provide default favicon if not provided
	favicon, err := pkg.NewFaviconFromFile(config.FavIcon)
	if err == nil {
		getObject = favicon.GetDefaultFavicon(getObject)
		log.Printf("default favicon: %s", config.FavIcon)
	}

	// return cache if available (1000 objects, max 10 Mb)
	getObject = pkg.NewCache(1000, 1024*1024*10).GetObjectCache(getObject)
	log.Printf("caching: enabled")

	// render markdown if needed
	serve := pkg.DefaultServe
	md, err := pkg.NewMarkdown(config.MarkdownTemplate)
	if err == nil {
		log.Printf("render markdown: enabled")
		serve = md.RenderMarkdown(serve)
	}

	// handler
	handler := pkg.Handlers{
		StatObject: helper.StatObject,
		GetObject:  getObject,
		Serve:      serve,
		SetHeaders: pkg.SetDefaultHeaders}

	http.HandleFunc("/", handler.Handler())
	log.Print("Listening to :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
