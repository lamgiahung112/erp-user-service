package main

import (
	"context"
	"erp-user-service/data"
	"erp-user-service/grpc/authenticate"
	"erp-user-service/grpc/logger"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServer struct {
	authenticate.UnimplementedAuthenticateServiceServer

	Models *data.Models
}

func (*UserServer) Authenticate(ctx context.Context, req *authenticate.AuthenticateRequest) (*authenticate.AuthenticateResponse, error) {
	// input := req.GetToken()

	return &authenticate.AuthenticateResponse{
		User: &authenticate.JwtUser{
			UserId: "ok",
			Email:  "email@ok",
			Name:   "Hung",
		},
	}, nil
}

func (app *Config) startGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if err != nil {
		log.Fatalf("Failed to listen to grpc: %v", err)
	}

	s := grpc.NewServer()

	authenticate.RegisterAuthenticateServiceServer(s, &UserServer{
		Models: app.Models,
	})

	log.Printf("grpc server started on port %s", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen to grpc port %v", grpcPort)
	}
}

func (app *Config) LogViaGRPC(event string, details string) {
	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	defer conn.Close()

	c := logger.NewLoggerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	r := &logger.LogRequest{
		Event:         event,
		Details:       details,
		Timestamp:     time.Now().Unix(),
		CallerService: "user-service",
	}

	c.WriteLog(ctx, r)
}
