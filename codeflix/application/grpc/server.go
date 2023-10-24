package grpc

import (
	"fmt"
	"log"
	"net"

	"derso.com/imersao-fullcycle/codepix-go/application/grpc/pb"
	"derso.com/imersao-fullcycle/codepix-go/application/usecase"
	"derso.com/imersao-fullcycle/codepix-go/infrastructure/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()

	// Para debug com o Evans
	reflection.Register(grpcServer)

	// TODO refatoração das dependências :P (ver pasta factory para exemplo)
	pixRepository := repository.PixKeyRepositoryDb{DB: database}
	pixUseCase := usecase.PixUseCase{PixKeyRepository: pixRepository}
	pixGrpcService := NewPixGrpcService(pixUseCase)

	// Registro de serviços (ver pix.go, ... neste pacote)
	pb.RegisterPixServiceServer(grpcServer, pixGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}

	log.Printf("gRPC server has been started on port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}
}
