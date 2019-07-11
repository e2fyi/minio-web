// Package main starts a web server which serve resources in S3 compatible
// backend with regular HTTP GET.
package main

import (
	app "github.com/e2fyi/minio-web/pkg/app"
	ext "github.com/e2fyi/minio-web/pkg/ext"
)

func main() {

	// create new app and load config
	app := app.NewApp().LoadConfig()
	// config backend
	app.ConfigMinioHelper(app.Config.Minio, app.Config.Ext.BucketName)
	// install default index file extension
	app.ApplyExtension(ext.DefaultIndexFileExtension(app.Config.Ext.DefaultHTMLs...))
	// install default favicon extension
	app.ApplyExtension(ext.DefaultFaviconExtension(app.Config.Ext.FavIcon))
	// return cache if available (1000 objects, max 10 Mb)
	app.ApplyExtension(ext.CacheRequestsExtension(1000, 1024*1024*10))
	// list folder if needed
	app.ApplyExtension(ext.ListFolderExtension(app.Helper, app.Config.Ext.ListFolder, app.Config.Ext.ListFolderObjects))
	// render markdown if needed
	app.ApplyExtension(ext.RenderMarkdownExtension(app.Config.Ext.MarkdownTemplate))
	// start server
	app.StartServer(app.Config.Server)
}
