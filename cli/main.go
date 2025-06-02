package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const url = "http://localhost:8080/tasks/%s/code"

func main() {
	if len(os.Args) != 3 {
		panic("need 3 arguments")
	}

	if os.Args[1] != "download" {
		panic("unknown argument")
	}

	exerciseNo := os.Args[2]

	err := os.Mkdir(exerciseNo, os.ModePerm)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(fmt.Sprintf(url, exerciseNo))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create the file
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		fmt.Println("Reading file:", zipFile.Name)
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			log.Println(err)
			continue
		}

		out, err := os.Create(exerciseNo + "/" + zipFile.Name)
		if err != nil {
			panic(err)
		}
		_, err = out.Write(unzippedFileBytes)
		if err != nil {
			panic(err)
		}

		out.Close()
	}

}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// func downloadFile(fullPath string) {
// 	out, err := os.Create(exerciseNo)
// 	if err != nil  {
// 		panic(err)
// 	}
// 	defer out.Close()
//
// 	resp, err := http.Get(fullPath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		panic("status: " + resp.Status)
// 	}
//
// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil  {
// 		panic(err)
// 	}
//
// }
