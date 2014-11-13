package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

type everything struct {
	is      string
	Awesome string
}

func main() {
	e := everything{is: "Everything is", Awesome: "Awesome!!!"}

	s := reflect.ValueOf(&e).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.CanSet() {
			fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		} else {
			fmt.Printf("%d: %s %s\n", i, typeOfT.Field(i).Name, f.Type())
		}
	}
	b, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to encode to JSON.")
	}
	fmt.Println(string(b))

	b, err = xml.Marshal(e)
	if err != nil {
		fmt.Println("failed to encode to XML.")
	}
	fmt.Println(string(b))

	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	if err := enc.Encode(e); err != nil {
		fmt.Println("failed to encode to Binary.", err)
	}
	fmt.Println(hex.EncodeToString(buffer.Bytes()))
}
