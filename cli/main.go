package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"encoding/json"
	"github.com/alecthomas/kong"
	
)

const (
	url       = "http://localhost:8080/tasks/%s/code"
	urlSubmit = "http://localhost:8080/submission"
)

var CLI struct {
	Configure struct {
		Token string `arg:"" name:"token" help:"Der Token ist auf der Homepage unter 'Wichtigsten Befehle'"`
	} `cmd:"" name:"configure" help:"Kopiere dein Token um Code hochladen und herunterladen zu können"`

	Download struct {
		Aufgabe bool `arg:"" name:"aufgabe" help:"Nummer der Aufgabe die gedownloadet werden soll."`
	} `cmd:"" name:"download" help:"Downloade eine Aufgabe"`

	Upload struct {
		Aufgabe bool   `arg:"" name:"aufgabe" help:"Nummer der Aufgabe für die eine Lösung hochgeladen wird."`
		Ordner  string `arg:"" name:"ordner" help:"Name des Ordners der hochgeladen werden soll" type:"existingdir"`
	} `cmd:"" name:"upload" help:"Upload deiner Lösung für die Aufgabe"`

	Review struct {
		Aufgabe  bool   `help:"Nummer der Aufgabe für die eine Lösung hochgeladen wird."`
		Benutzer string `help:"Nutzername der Person die Lösung hochgeladen hat"`
	} `cmd:"" name:"review" help:"Downloade eine Aufgabe"`
}

func main() {
	ctx := kong.Parse(&CLI)

	switch ctx.Command() {
	case "configure <token>":
		configure()
	case "download":
		download()
	case "upload":
		upload()
	case "review":
		// Todo: downloading other submissions
	default:
		panic(ctx.Command())
	}
}

func configure() {
	configPath := configdir.LocalConfig("judge")
	err := configdir.MakePath(configPath)
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "settings.json")
	type Config struct {
		Token string `json:"token"`
	}

	config := Config{Token: os.Args[2]}
	fh, err := os.Create(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	err = json.NewEncoder(fh).Encode(&config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created a config at " + configPath)
}

func upload() {
	token := loadConfig()

	exerciseNo := os.Args[2]

	if _, err := os.Stat(exerciseNo); os.IsNotExist(err) {
		panic("folder does not exist: " + exerciseNo) // TODO: NO panics
	}

	var buf bytes.Buffer
	w := zip.NewWriter(base64.NewEncoder(base64.StdEncoding, &buf))

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err := filepath.Walk(exerciseNo, walker)
	if err != nil {
		panic(err)
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, urlSubmit, &buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("token", token)
	req.Header.Set("exercise", exerciseNo)
	req.Header.Set("Content-Type", "application/zip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		panic(err)
	}
}

func download() {
	//token := loadConfig()

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

type Config struct {
	Token string `json:"token"`
}

func loadConfig() string {
	configPath := configdir.LocalConfig("judge")
	configFile := filepath.Join(configPath, "settings.json")

	var config Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic("no token found. Please run 'configure' first")
	}

	fh, err := os.Open(configFile)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer fh.Close()

	err = json.NewDecoder(fh).Decode(&config)
	if err != nil {
		log.Println(err)
		return ""
	}
	return config.Token
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
