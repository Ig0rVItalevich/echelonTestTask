package thumbnail

import (
	"errors"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/http"
	"os"
)

func ThumbnailHandler(thumbnail *youtube.ThumbnailDetails, videoLink string) (string, error) {
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
		thumbnailLink = ""
		return thumbnailLink, errors.New("video does not have thumbnail")
	}

	return downloadThumbnail(thumbnailLink, videoLink)
}

func downloadThumbnail(thumbnailLink string, videoLink string) (string, error) {
	res, err := http.Get(thumbnailLink)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	if err := checkImageDir(); err != nil {
		return "", err
	}
	pathFile := fmt.Sprintf("./images/%s.jpg", videoLink)
	file, err := os.Create(pathFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(body)

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
