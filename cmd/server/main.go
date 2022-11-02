package main

import (
	"context"
	"fmt"
	api "github.com/Ig0rVItalevich/echelon/pkg/api/proto"
	"github.com/Ig0rVItalevich/echelon/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
	"net"
)

func main() {
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(viper.GetString("api_key")))
	if err != nil {
		logrus.Fatalf("error creating youtubeService: %s", err.Error())
	}

	s := grpc.NewServer()
	srv := server.NewServer(youtubeService)
	api.RegisterThumbnailsServer(s, srv)

	address := fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Fatalf("error creating listener: %s", err.Error())
	}

	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("server listener service error: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
