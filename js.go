package main

import (
	"os"
	"strings"
	"io/ioutil"

	"github.com/gabriel-vasile/mimetype"
)

func getHTML() string {
	f, err := os.Open("decrypt.html")
	HandleErr(err)
	data, err = ioutil.ReadAll(f)
	HandleErr(err)

	return string(data)
}

func DecrypterFromFile(fname string, b64 string) []byte {
	html := getHTML()

	html = strings.Replace(html, "{{DATA}}", b64, 1)
	html = strings.Replace(html, "{{BINARY}}", "true", 1)
	
	mtype, err := mimetype.DetectFile(fname)
	HandleErr(err)

	html = strings.Replace(html, "{{MIME}}", mtype.String(), 1)
	html = strings.Replace(html, "{{EXT}}", mtype.Extension(), 1)

	return []byte(html)
}

func DecrypterFromPaste(b64 string) []byte {
	html := getHTML()

	html = strings.Replace(html, "{{DATA}}", b64, 1)
	html = strings.Replace(html, "{{BINARY}}", "false", 1)

	return []byte(html)
}