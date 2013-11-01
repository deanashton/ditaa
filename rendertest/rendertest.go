package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

const (
	sources = "../orig-java/tests/xmls"
	results = "imgs"
)

type Grid struct {
	XMLName xml.Name `xml:"grid"`
	W       int      `xml:"width,attr"`
	H       int      `xml:"height,attr"`
}

type Diagram struct {
	XMLName xml.Name `xml:"diagram"`
}

func LoadDiagram(path string) (*Diagram, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loading diagram from '%s': %s", path, err)
	}
	diagram := Diagram{}
	err = xml.NewDecoder(r).Decode(&diagram)
	if err != nil {
		return nil, fmt.Errorf("decoding diagram from '%s': %s", path, err)
	}
	return &diagram, nil
}

func runRender(src, dst string) error {
	diagram, err := LoadDiagram(src)
	_ = diagram
	return err
}

func run() error {
	fnames := []string{}

	os.Mkdir(results, 0666)

	//todo: load files from ../orig-java/tests/xmls/*.xml, then try to render them into some output dir, and link them all on one html page
	err := filepath.Walk(sources, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}
		fnames = append(fnames, info.Name())
		return runRender(path, filepath.Join(results, info.Name()))
	})

	if err != nil {
		return err
	}

	return err
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
