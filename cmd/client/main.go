package main

import (
	"context"
	"flag"
	"fmt"
	api "github.com/Ig0rVItalevich/echelon/pkg/api/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net/url"
)

func main() {
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	flag.Parse()
	if flag.NArg() < 1 {
		logrus.Fatal("not enough arguments")
	}

	videoUrl := flag.Arg(0)
	link := getLink(videoUrl)

	address := fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error while dial server: %s", err.Error())
	}

	client := api.NewThumbnailsClient(conn)
	response, err := client.Get(context.Background(), &api.GetRequest{Link: link})
	if err != nil {
		logrus.Fatalf("error while requesting GRPC-server: %s", err.Error())
	}

	fmt.Printf("%s\n", response.Thumbnail)
}

func getLink(videoUrl string) string {
	urlParse, err := url.Parse(videoUrl)
	if err != nil {
		logrus.Fatalf("error while parsing url: %s", err.Error())
	}
	urlParams, err := url.ParseQuery(urlParse.RawQuery)
	if err != nil {
		logrus.Fatalf("error while parsing url params: %s", err.Error())
	}
	link := urlParams.Get("v")

	return link
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
