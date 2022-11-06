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
	defer func() {
		err := conn.Close()
		if err != nil {
			logrus.Fatalf("error while connection closing: %s", err.Error())
		}
	}()

	client := api.NewThumbnailsClient(conn)

	asyncFlag := flag.Bool("async", false, "implements asynchronous execution")
	flag.Parse()

	argsHandler(os.Args, *asyncFlag, client)
}

func argsHandler(args []string, async bool, client api.ThumbnailsClient) {
	if async {
		args = args[2:]
	} else {
		args = args[1:]
	}

	var wg sync.WaitGroup
	for _, arg := range args {
		if async {
			wg.Add(1)
			go requestServer(arg, client, &wg)
		} else {
			requestServer(arg, client, nil)
		}
	}

	if async {
		wg.Wait()
	}
}

func requestServer(videoUrl string, client api.ThumbnailsClient, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	videoId := getVideoId(videoUrl)
	response, err := client.Get(context.Background(), &api.GetRequest{VideoId: videoId})
	if err != nil {
		logrus.Fatalf("error while requesting gRPC-server: %s", err.Error())
	}

	fmt.Printf("Thumbnail-file of video %s: %s\n", videoUrl, response.Thumbnail)
}

func getVideoId(videoUrl string) string {
	urlParse, err := url.Parse(videoUrl)
	if err != nil {
		logrus.Fatalf("error while parsing url: %s", err.Error())
	}
	urlParams, err := url.ParseQuery(urlParse.RawQuery)
	if err != nil {
		logrus.Fatalf("error while parsing url params: %s", err.Error())
	}
	videoId := urlParams.Get("v")

	return videoId
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
