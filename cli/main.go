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

	"github.com/kirsle/configdir"
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
		Aufgabe string `arg:"" name:"aufgabe" help:"Nummer der Aufgabe die gedownloadet werden soll."`
	} `cmd:"" name:"download" help:"Downloade eine Aufgabe"`

	Upload struct {
		Aufgabe string `arg:"" name:"aufgabe" help:"Nummer der Aufgabe für die eine Lösung hochgeladen wird."`
		Ordner  string `arg:"" name:"ordner" help:"Name des Ordners der hochgeladen werden soll" type:"existingdir"`
	} `cmd:"" name:"upload" help:"Upload deiner Lösung für die Aufgabe"`

	Review struct {
		Aufgabe  string `help:"Nummer der Aufgabe für die eine Lösung hochgeladen wird."`
		Benutzer string `help:"Nutzername der Person die Lösung hochgeladen hat"`
	} `cmd:"" name:"review" help:"Downloade eine Aufgabe"`
}

func main() {
	ctx := kong.Parse(&CLI, kong.Name("judge-cli"),
		kong.Description("Zum downloaden, uploaden und anschauen von Übungen"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "configure <token>":
		configure()
	case "download":
		download(CLI.Download.Aufgabe)
	case "upload <aufgabe> <ordner>":
		upload(CLI.Upload.Aufgabe, CLI.Upload.Ordner)
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

func upload(aufgabe string, ordner string) {
	token := loadConfig()

	if _, err := os.Stat(ordner); os.IsNotExist(err) {
		panic("folder does not exist: " + ordner) // TODO: NO panics
	}

	// 1. Zip the folder
	zipBytes, err := zipDirectory(ordner)
	if err != nil {
		panic(err)
	}

	// 2. Base64-encode zip
	b64 := base64.StdEncoding.EncodeToString(zipBytes)

	// 4. POST it
	req, err := http.NewRequest(http.MethodPost, urlSubmit, bytes.NewReader([]byte(b64)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("token", token)
	req.Header.Set("exercise", aufgabe)
	req.Header.Set("Content-Type", "text/plain")

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

func download(aufgabe string) {
	//token := loadConfig()

	err := os.Mkdir(aufgabe, os.ModePerm)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(fmt.Sprintf(url, aufgabe))
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

		out, err := os.Create(aufgabe + "/" + zipFile.Name)
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
