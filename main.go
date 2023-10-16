package main

import (
	"ZipFiles/internal/utils"
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	ginHttpServer()
}

func ginHttpServer() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		// TODO validate files input
		queryFiles := c.Query("files")

		files := strings.Split(queryFiles, ",")

		zipFile := start(files)

		c.FileAttachment(zipFile, zipFile)

		// TODO clean files
	})

	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func start(files []string) string {
	destFolder := utils.RandStringBytes(10)

	if err := os.Mkdir(destFolder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		i := i

		go func() {
			defer wg.Done()
			worker(i, destFolder, files[i])
		}()
	}

	wg.Wait()

	if err := zipSource(destFolder, destFolder+".zip"); err != nil {
		log.Fatal(err)
	}

	return destFolder + ".zip"
}

func worker(id int, destFolder string, fileUrl string) {
	fmt.Printf("Worker %d starting\n", id)

	downloadFile(destFolder, fileUrl)

	fmt.Printf("Worker %d done\n", id)
}

func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func downloadFile(destFolder string, fullURLFile string) {

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Extension:", ext)

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
