module github.com/e2fyi/minio-web

go 1.12

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.3

// see https://github.com/Azure/go-autorest/issues/449
replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8

// -20190108154635-e99af5d43a04
// e99af5d43a047907825b8231e080e39665aef867
// 47c0da630f72

require (
	github.com/bluele/gcache v0.0.0-20190301044115-79ae3b2d8680
	github.com/dustin/go-humanize v1.0.0
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/gorilla/pat v1.0.1
	github.com/markbates/goth v1.56.0
	github.com/micro/go-config v1.1.0
	github.com/micro/go-micro v1.7.0
	github.com/minio/minio-go v6.0.14+incompatible
	gitlab.com/golang-commonmark/html v0.0.0-20180917080848-cfaf75183c4a // indirect
	gitlab.com/golang-commonmark/linkify v0.0.0-20180917065525-c22b7bdb1179 // indirect
	gitlab.com/golang-commonmark/markdown v0.0.0-20181102083822-772775880e1f
	gitlab.com/golang-commonmark/mdurl v0.0.0-20180912090424-e5bce34c34f2 // indirect
	gitlab.com/golang-commonmark/puny v0.0.0-20180912090636-2cd490539afe // indirect
	gitlab.com/opennota/wd v0.0.0-20180912061657-c5d65f63c638 // indirect
	go.uber.org/zap v1.9.1
	gopkg.in/ini.v1 v1.42.0 // indirect
)
