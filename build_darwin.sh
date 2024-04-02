#!/bin/bash

#Go编译
#Golang 支持在一个平台下生成另一个平台可执行程序的交叉编译功能。
#Mac下编译Linux, Windows平台的64位可执行程序：
#$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/gredis ./main.go
#$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/gredis ./main.go
#Linux下编译Mac, Windows平台的64位可执行程序：
#$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/gredis ./main.go
#$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/gredis ./main.go
#Windows下编译Mac, Linux平台的64位可执行程序：
#$ SET CGO_ENABLED=0SET GOOS=darwin3 SET GOARCH=amd64 go build -o ./bin/gredis ./main.go
#$ SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build -o ./bin/gredis ./main.go
#注：如果编译web等工程项目，直接cd到工程目录下直接执行以上命令

go build -o ./bin/gredis ./main.go