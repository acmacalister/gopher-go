package main

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	processCommandLine()
}

//Process Command Line Args.
func processCommandLine() {
	if len(os.Args) < 3 {
		printUsage()
	} else {
		file, err := os.Open(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		cmd := os.Args[1]
		switch cmd {
		case "gzip":
			gzipFiles(file)
		case "bzip":
			bzipFiles(file)
		case "zlibFiles":
			zlibFiles(file)
		default:
			printUsage()
		}
	}
}

// Create new Reader for gzip files.
func gzipFiles(f *os.File) {
	in, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	decompress(in, ".gz")
}

// Create new Reader for bzip2 files.
func bzipFiles(f *os.File) {
	in := bzip2.NewReader(f)
	decompress(in, ".bz2")
}

// Create new Reader for zlib files.
func zlibFiles(f *os.File) {
	in, err := zlib.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	decompress(in, ".zlib")

}

// Create a new file and write the contents of the
// compressed file to it.
func decompress(r io.Reader, ext string) {
	trimmed := strings.TrimSuffix(os.Args[2], ext)
	out, err := os.Create(trimmed)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(out, r)
	if err != nil {
		log.Fatal(err)
	}
	out.Close()
}

// Prints a simple help menu.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("decompress gzip file.gz")
	fmt.Println("decompress bzip file.bz2")
	fmt.Println("decompress zlib file.zlib")
}
