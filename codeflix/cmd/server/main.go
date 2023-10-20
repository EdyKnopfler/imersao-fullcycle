package main

import (
	"os"

	"derso.com/imersao-fullcycle/codepix-go/application/grpc"
	"derso.com/imersao-fullcycle/codepix-go/infrastructure/db"
	"gorm.io/gorm"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}
