package main

import (
	"archive/tar"
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type fileType struct {
	Name string
	Body []byte
}

var (
	method      string
	archivePath string
)

func init() {
	flag.StringVar(&method, "method", "tar", "Compression method: tar or zip")
	flag.StringVar(&archivePath, "path", "archive.tar", "Path of the output archived file")
	flag.Parse()
	if len(flag.Args()) == 0 {
		printUsage()
	}
	fmt.Println("method:", method)
	fmt.Println("archivePath is:", archivePath)
	fmt.Println("flag args are:", flag.Args())
}

func main() {
	//processCommandLine()
}

// func processCommandLine() {
// 	if len(os.Args) < 4 {
// 		printUsage()
// 	} else {
// 		cmd := os.Args[1]
// 		switch cmd {
// 		case "tar":
// 			tarFiles()
// 		case "zip":
// 			zipFiles()
// 		default:
// 			printUsage()
// 		}
// 	}
// }

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("archive tar path file1 file2 file3")
	fmt.Println("archive zip path file1 file2 file3")
}

func tarFiles() {
	//Create a file to write to.
	tarFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	//Create a new tar archive.
	tarWriter := tar.NewWriter(tarFile)

	//Read the files from disk and loop over them to create the tar
	files := readFilesFromDisk()
	for _, file := range files {
		header := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			log.Fatalln(err)
		}
		if _, err := tarWriter.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}

	// Make sure to check the error on Close.
	if err := tarWriter.Close(); err != nil {
		log.Fatalln(err)
	}
}

func zipFiles() {
	//Create a file to write to.
	zipFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	//Create a new zip archive.
	zipWriter := zip.NewWriter(zipFile)

	//Read the files from disk and loop over them to create the zip
	files := readFilesFromDisk()
	for _, file := range files {
		f, err := zipWriter.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	//Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func readFilesFromDisk() []fileType {
	//Get the arguments after the program name and top level command.
	filePaths := os.Args[3:]

	//Make a channel and run a go routines to read the files.
	c := make(chan []byte)
	for _, file := range filePaths {
		go func(file string) {
			buffer, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			c <- buffer
		}(file)
	}

	//Loop over our file paths again append them to our files list.
	files := make([]fileType, 0, len(filePaths))
	for _, filePath := range filePaths {
		_, fileName := path.Split(filePath)
		files = append(files, fileType{fileName, <-c})
	}

	return files
}
