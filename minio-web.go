// Package main starts a web server which serve resources in S3 compatible
// backend with regular HTTP GET.
package main

import (
	"github.com/e2fyi/minio-web/pkg"
)

func main() {

	// create new app and load config
	app := pkg.NewApp().LoadConfig()
	// config backend
	app.ConfigMinioHelper(app.Config.Minio, app.Config.Ext.BucketName)
	// install default index file extension
	app.InsertIndexFile(app.Config.Ext.DefaultHTML)
	// install default favicon extension
	app.SetDefaultFavicon(app.Config.Ext.FavIcon)
	// return cache if available (1000 objects, max 10 Mb)
	app.CacheRequests(1000, 1024*1024*10)
	// render markdown if needed
	app.RenderMarkdowns(app.Config.Ext.MarkdownTemplate)
	// start server
	app.StartServer(app.Config.Server)
}
