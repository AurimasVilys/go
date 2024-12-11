package cmd

import (
	"fmt"
	"github.com/stephenafamo/bob"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	pb "scootin/internal/grpc"
	"scootin/internal/grpc/server"
	"scootin/internal/repository"
	"sync"
)

func InitializeGrpcApi(wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, _ := openDB()
	scooterEventServer, _ := InitializeScooterEventServer(db)

	s := grpc.NewServer()
	pb.RegisterScooterEventServiceServer(s, scooterEventServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitializeScooterEventServer(db *bob.DB) (*server.ScooterEventServer, error) {
	scooterRepository := repository.NewScooterRepository(db)
	scooterEventRepository := repository.NewScooterEventRepository(db)
	scooterEventServer := server.NewScooterEventServer(scooterRepository, scooterEventRepository)

	return scooterEventServer, nil
}
