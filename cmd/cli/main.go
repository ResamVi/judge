package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ResamVi/judge/cli"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cli.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

// -----------------------------------------------------------------------------
func upload(aufgabe string, ordner string) {
	//token := loadConfig()
	//
	//if _, err := os.Stat(ordner); os.IsNotExist(err) {
	//	panic("folder does not exist: " + ordner) // TODO: NO panics
	//}
	//
	//// 1. Zip the folder
	//zipBytes, err := zipDirectory(ordner)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 2. Base64-encode zip
	//b64 := base64.StdEncoding.EncodeToString(zipBytes)
	//
	//// 4. POST it
	//req, err := http.NewRequest(http.MethodPost, urlSubmit, bytes.NewReader([]byte(b64)))
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("token", token)
	//req.Header.Set("exercise", aufgabe)
	//req.Header.Set("Content-Type", "text/plain")
	//
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//_, err = io.Copy(os.Stdout, resp.Body)
	//if err != nil {
	//	panic(err)
	//}
}

func download(aufgabe string) {
	//token := loadConfig()

	//err := os.Mkdir(aufgabe, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//
	//resp, err := http.Get(fmt.Sprintf(url, aufgabe))
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Create the file
	//zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Read all the files from zip archive
	//for _, zipFile := range zipReader.File {
	//	fmt.Println("Reading file:", zipFile.Name)
	//	unzippedFileBytes, err := readZipFile(zipFile)
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//
	//	out, err := os.Create(aufgabe + "/" + zipFile.Name)
	//	if err != nil {
	//		panic(err)
	//	}
	//	_, err = out.Write(unzippedFileBytes)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	out.Close()
	//}
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// zipDirectory walks dirPath and writes a .zip (in-memory) with full structure.
// Returns the raw bytes of the zip file.
func zipDirectory(dirPath string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Build the path inside the zip (relative to dirPath)
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Skip the root itself ("."), but still include its children
		if relPath == "." {
			return nil
		}

		// Directories in zip need a trailing slash
		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// It's a file: create header, copy contents
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileHeader.Name = relPath
		fileHeader.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		return err
	})

	if err != nil {
		zipWriter.Close()
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
