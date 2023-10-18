package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ZipFiles/internal/storage"
	"ZipFiles/internal/utils"

	"github.com/gin-gonic/gin"
)

type filesInput struct {
	files string
}

func init() {
	log.Println("Hello from init func")
}

func main() {
	s := storage.GetStorage(os.Getenv("STORAGE_CLIENT"))

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		// TODO validate files input
		queryFiles := c.Query("files")

		log.Println("queryFiles", queryFiles)

		files := strings.Split(queryFiles, ",")

		destFolder := worker(s, files)
		zipFile := destFolder + ".zip"

		c.FileAttachment(zipFile, zipFile)

		go cleanFiles(destFolder)
	})

	log.Printf("running get and post...")

	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func worker(storageClient storage.ClientManager, files []string) string {
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
			fmt.Printf("Worker %d starting file %s \n", i, files[i])

			storageClient.Download(destFolder, files[i])

			fmt.Printf("Worker %d done\n", i)
		}()
	}

	wg.Wait()

	if err := zipSource(destFolder, destFolder+".zip"); err != nil {
		log.Fatal(err)
	}

	return destFolder
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

func cleanFiles(dest string) {
	log.Println("cleaning files", dest)

	err := os.RemoveAll(dest)
	if err != nil {
		log.Panic(err)
	}

	err = os.Remove(dest + ".zip")
	if err != nil {
		log.Panic(err)
	}
}
