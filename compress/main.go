package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	processCommandLine()
}

func processCommandLine() {
	if len(os.Args) < 3 {
		printUsage()
	} else {
		cmd := os.Args[1]
		switch cmd {
		case "gzip":
			gzipFiles()
		case "bzip":
			bzipFiles()
		case "zlibFiles":
			zlibFiles()
		default:
			printUsage()
		}
	}
}

func gzipFiles() {
	b := readFile()
	r, err := gzip.NewReader(b)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	filename := strings.TrimSuffix(os.Args[2], ".gz")
	writeFile(filename, r)
}

func bzipFiles() {
	b := readFile()
	r := bzip2.NewReader(b)
	filename := strings.TrimSuffix(os.Args[2], ".bz2")
	writeFile(filename, r)
}

func zlibFiles() {
	b := readFile()
	r, err := zlib.NewReader(b)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	filename := strings.TrimSuffix(os.Args[2], ".zlib")
	writeFile(filename, r)
}

func writeFile(filename string, r io.Reader) {
	fo, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	w := bufio.NewWriter(fo)
	w.ReadFrom(r)
	if err = w.Flush(); err != nil {
		log.Fatal(err)
	}
}

func readFile() *bytes.Reader {
	buffer, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewReader(buffer)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("compress gzip file.gz")
	fmt.Println("compress bzip file.bz2")
	fmt.Println("compress zlib file.zlib")
}
