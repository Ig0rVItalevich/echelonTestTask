package thumbnail

import (
	"errors"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/http"
	"os"
)

func ThumbnailHandler(thumbnail *youtube.ThumbnailDetails) ([]byte, error) {
	var thumbnailLink string
	if thumbnail.High != nil {
		thumbnailLink = thumbnail.High.Url
	} else if thumbnail.Maxres != nil {
		thumbnailLink = thumbnail.Maxres.Url
	} else if thumbnail.Medium != nil {
		thumbnailLink = thumbnail.Medium.Url
	} else if thumbnail.Standard != nil {
		thumbnailLink = thumbnail.Standard.Url
	} else if thumbnail.Default != nil {
		thumbnailLink = thumbnail.Default.Url
	} else {
		return nil, errors.New("video does not have thumbnail")
	}

	return downloadThumbnail(thumbnailLink)
}

func downloadThumbnail(thumbnailLink string) ([]byte, error) {
	res, err := http.Get(thumbnailLink)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func SaveThumbnail(thumbnail []byte, videoId string) (string, error) {
	if err := checkImageDir(); err != nil {
		return "", err
	}
	pathFile := fmt.Sprintf("./images/%s.jpg", videoId)
	file, err := os.Create(pathFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(thumbnail)

	return pathFile, nil
}

func checkImageDir() error {
	info, err := os.Stat("./images")
	if os.IsNotExist(err) || !info.IsDir() {
		if err := os.Mkdir("images", 0777); err != nil {
			return err
		}
	}

	return nil
}
