package storage

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

type PublicStorage struct {
	client http.Client
}

func newPublicStorage() ClientManager {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return PublicStorage{
		client: client,
	}
}

func (d PublicStorage) Download(destFolder string, filepath string) {

	// Put content on file
	resp, err := d.client.Get(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Fatal(err)
	}

	// fileName := destFolder + "/" + utils.RandStringBytes(8)
	fileName := destFolder + "/" + path.Base(resp.Request.URL.String()) + "." + ext[0]
	// fmt.Println("fileName", fileName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)

}
