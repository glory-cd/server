module github.com/glory-cd/server

go 1.12

replace github.com/glory-cd/server => ./

require (
	github.com/glory-cd/utils v0.0.0-20190731013124-f69d4a28bc82
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	google.golang.org/grpc v1.19.0
)
