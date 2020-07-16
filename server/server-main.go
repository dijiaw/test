package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "./userproto"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type server struct{}

var c *cache.Cache

func (*server) GetUser(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	username := r.GetUsername()

	userVal, found := c.Get(username)
	if found {
		u := userVal.(userproto.User)
		return &user.UserResponse{
			Result: &u,
		}, nil
	}

	return nil, fmt.Errorf("could not find user with id: %v", id)
}

func main() {
	fmt.Println("Starting gRPC micro-service...")
	c = cache.New(60*time.Minute, 70*time.Minute)
	c.Set("f", user.User{
		Id:       "1",
		Username: "f",
		Nickname: "f",
		Password: "f",
	}, cache.DefaultExpiration)

	l, e := net.Listen("tcp", ":50051")
	if e != nil {
		log.Fatalf("Failed to start listener %v", e)
	}

	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{})

	if e := s.Serve(l); e != nil {
		log.Fatalf("failed to serve %v", e)
	}
}
