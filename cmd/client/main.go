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
	"os"
	"sync"
)

func main() {
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	address := fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error while dial server: %s", err.Error())
	}

	client := api.NewThumbnailsClient(conn)

	asyncFlag := flag.Bool("async", false, "implements asynchronous execution")
	flag.Parse()

	if *asyncFlag {
		var wg sync.WaitGroup
		args := os.Args[2:]
		for _, arg := range args {
			wg.Add(1)
			go requestServer(arg, client, &wg)
		}

		wg.Wait()
	} else {
		args := os.Args[1:]
		for _, arg := range args {
			requestServer(arg, client, nil)
		}
	}
}

func requestServer(videoUrl string, client api.ThumbnailsClient, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	link := getLink(videoUrl)
	response, err := client.Get(context.Background(), &api.GetRequest{Link: link})
	if err != nil {
		logrus.Fatalf("error while requesting gRPC-server: %s", err.Error())
	}

	fmt.Printf("Thumbnail-file of video %s: %s\n", videoUrl, response.Thumbnail)
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

//go run cmd/client/main.go --async "https://www.youtube.com/watch?v=qu3Vpdnndi4&ab_channel=Rozetked" "https://www.youtube.com/watch?v=2maIPAWo-UM&ab_channel=Rozetked" "https://www.youtube.com/watch?v=KnINsmXT9_c&ab_channel=Droider"
