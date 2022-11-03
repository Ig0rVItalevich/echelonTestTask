package main

import (
	"context"
	"fmt"
	api "github.com/Ig0rVItalevich/echelon/pkg/api/proto"
	"github.com/Ig0rVItalevich/echelon/pkg/cache"
	"github.com/Ig0rVItalevich/echelon/pkg/server"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(viper.GetString("api_key")))
	if err != nil {
		logrus.Fatalf("error creating youtubeService: %s", err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: "",
		DB:       0,
	})
	cacheServer := cache.NewCache(redisClient)

	s := grpc.NewServer()
	srv := server.NewServer(youtubeService, cacheServer)
	api.RegisterThumbnailsServer(s, srv)

	address := fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Fatalf("error creating listener: %s", err.Error())
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			logrus.Fatalf("server listener service error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	s.Stop()
	if err := srv.Cache.DB.Close(); err != nil {
		logrus.Fatalf("error occured while closing cache db: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
