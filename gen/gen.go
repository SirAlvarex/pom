package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/tools/imports"
)

var fileName = "gen_models.go"

func main() {
	// This is the URL to the POM Schema Definition
	URL := "https://maven.apache.org/xsd/maven-4.0.0.xsd"
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Write the body to a buffer
	out := new(bytes.Buffer)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Unmarshal the schema definition into a schema object
	data := out.Bytes()
	schema := Schema{}
	err = xml.Unmarshal(data, &schema)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Create and open the file we are writing
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Create the unformatted version of the models
	types := schema.GetTypes()
	err = modelFormat.Execute(f, struct {
		Types     []string
		Timestamp time.Time
	}{types, time.Now()})

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Pass the models through go-imports, which will handle any import paths
	// and run a go fmt on the file
	res, err := imports.Process(fileName, nil, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	err = ioutil.WriteFile(fileName, res, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return
}
