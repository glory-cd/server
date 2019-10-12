module github.com/glory-cd/server

go 1.12

require (
	github.com/glory-cd/utils v0.0.0-20190918013744-736ed1a53da9
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/robfig/cron/v3 v3.0.0
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/tredoe/osutil v0.0.0-20161130133508-7d3ee1afa71c
	google.golang.org/grpc v1.19.0
)

replace (
	github.com/glory-cd/server => ./
	github.com/glory-cd/utils => E:\GoProject\cdp\utils
)
