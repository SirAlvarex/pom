package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"

	"golang.org/x/tools/imports"
)

var format = `// Code generated DO NOT EDIT
// This file was generated by robots at
// {{ .Timestamp }}
package pom

// XMLMap is a custom key used to let XML data parse maps
// Because it doesnt do that by default...for some reason.
type XMLMap map[string]string

type xmlMapEntry struct {
    XMLName xml.Name
    Value   string ` + "`xml:\",chardata\"`" + `
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m XMLMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if len(m) == 0 {
        return nil
    }

    err := e.EncodeToken(start)
    if err != nil {
        return err
    }

    for k, v := range m {
        e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
    }

    return e.EncodeToken(start.End())
}


// UnmarshalXML takes a key and turns it into a map
func (m *XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    *m = XMLMap{}
    for {
        var e xmlMapEntry

        err := d.Decode(&e)
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        (*m)[e.XMLName.Local] = e.Value
    }
    return nil
}

{{ range .Types }}
{{ . }}
{{ end }}
`

func main() {
	URL := "https://maven.apache.org/xsd/maven-4.0.0.xsd"
	out := new(bytes.Buffer)
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	data := out.Bytes()
	schema := Schema{}
	err = xml.Unmarshal(data, &schema)
	if err != nil {
		fmt.Println(err)
	}

	types := schema.GetTypes()
	f, _ := os.Create("gen_models.go")
	t, _ := template.New("tmp").Parse(format)
	t.Execute(f, struct {
		Types     []string
		Timestamp time.Time
	}{types, time.Now()})
	res, _ := imports.Process("gen_models.go", nil, nil)

	err = ioutil.WriteFile("gen_models.go", res, 0)

	return
}