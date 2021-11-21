package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
)

func main() {
	url := os.Args[1]

	file, err := downloadFile(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := base64.StdEncoding.EncodeToString(file)
	imageDecoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	img, err := jpeg.Decode(imageDecoder)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	colors, err := prominentcolor.Kmeans(img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, c := range colors {
		fmt.Println(c.AsString())
	}
}

func downloadFile(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Non 200 http response")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
