package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
)

func StringifyToFile[A any](m A, file string) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	asStr := base64.StdEncoding.EncodeToString(b.Bytes())
	os.WriteFile(file, []byte(asStr), 0644)
}

func ParseFromFile[A any](file string) A {
	str, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var m A
	by, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}
