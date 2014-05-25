package main

import (
	//"archive/tar"
	//"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var files []struct {
	Name string
	Body []byte
}

func main() {
	processCommandLine()
}

func processCommandLine() {
	if len(os.Args) < 3 {
		printUsage()
	} else {
		cmd := os.Args[1]
		switch cmd {
		case "tar":
			tarFiles()
		case "untar":
			untarFiles()
		case "zip":
			zipFiles()
		case "unzip":
			unzipFiles()
		default:
			printUsage()
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("archive tar file1 file2 file3")
	fmt.Println("archive untar file1")
	fmt.Println("archive zip file1 file2 file3")
	fmt.Println("archive unzip file1")
}

func tarFiles() {
	readFilesFromDisk()
}

func untarFiles() {

}

func zipFiles() {

}

func unzipFiles() {

}

func readFilesFromDisk() {
	files := os.Args[2:]
	c := make(chan []byte)
	for i, file := range files {
		go func(file string) {
			bag, _ := ioutil.ReadFile(file)
			c <- bag
		}(file)
	}

	for i := range files {
		fmt.Println(i)
		fmt.Println(string(<-c))
	}
}
