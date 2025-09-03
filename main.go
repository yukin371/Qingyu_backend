package main

import (
	"log"

	"Qingyu_backend/core"
)

func main() {
	// 初始化MongoDB连接
	if err := core.InitMongoDB(); err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	// 运行服务器
	core.RunServer()
}
