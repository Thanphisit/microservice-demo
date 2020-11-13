package cmd

import (
	"github.com/Joker666/microservice-demo/apiService/server"
	pb "github.com/Joker666/microservice-demo/protos/api"
	"github.com/Joker666/microservice-demo/protos/user"

	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// srvCmd is the serve sub command to start the api server
var srvCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves the api server",
	RunE:  serve,
}

func init() {
	//
}

func serve(cmd *cobra.Command, args []string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := ":" + os.Getenv("PORT")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userConn, err := grpc.Dial(os.Getenv("USER_ADDRESS"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer userConn.Close()
	userSvcClient := user.NewUserSvcClient(userConn)

	s := grpc.NewServer()
	pb.RegisterAPIServer(s, server.New(userSvcClient))

	log.Println("Starting GRPC server at: " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
