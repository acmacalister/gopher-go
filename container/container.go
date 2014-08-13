package main

import (
	"container/list"
	"fmt"
	"time"
)

const (
	vulcanQuote = "Live long and prosper."
	size        = 1000
)

func main() {
	l := list.List{}
	s := []string{}
	m := make(map[int]string)
	insertLinkList(&l)
	iterateLinkList(&l)
	insertSlice(s)
	iterateSlice(s)
	insertMap(m)
	iterateMap(m)
}

func insertLinkList(l *list.List) {
	defer timeTrack(time.Now(), "insertLinkList")
	for i := 0; i < size; i++ {
		l.PushBack(vulcanQuote)
	}
}

func iterateLinkList(l *list.List) {
	defer timeTrack(time.Now(), "iterateLinkList")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Sprintln(e.Value)
	}
}

func insertSlice(slice []string) {
	defer timeTrack(time.Now(), "insertSlice")
	for i := 0; i < size; i++ {
		slice = append(slice, vulcanQuote)
	}
}

func iterateSlice(slice []string) {
	defer timeTrack(time.Now(), "iterateSlice")
	for _, value := range slice {
		fmt.Sprintln(value)
	}
}

func insertMap(m map[int]string) {
	defer timeTrack(time.Now(), "insertMap")
	for i := 0; i < size; i++ {
		m[i] = vulcanQuote
	}
}

func iterateMap(m map[int]string) {
	defer timeTrack(time.Now(), "iterateMap")
	for _, v := range m {
		fmt.Sprintln(v)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
