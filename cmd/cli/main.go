package main

import (
	"archive/zip"
	"fmt"
	"github.com/ResamVi/judge/cli"
	tea "github.com/charmbracelet/bubbletea"
	"io"
)

func main() {
	p := tea.NewProgram(cli.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

// -----------------------------------------------------------------------------

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
