package pom

//go:generate go run gen/main.go gen/models.go gen/templates.go gen/build.go

import (
	"encoding/xml"
	"strings"
)

var pomProjectHeader = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">`

// Unmarshal takes in the raw data of a POM, and returns a project in the form of a Model
func Unmarshal(rawPom []byte) (Model, error) {
	pom := project{}
	err := xml.Unmarshal(rawPom, &pom)
	return pom.Model, err
}

// Marshal turns a POM project into the raw bytes of a pom, ready for export
func Marshal(pom Model) ([]byte, error) {
	p := project{pom}
	data, err := xml.MarshalIndent(p, "", "    ")
	if err != nil {
		return data, err
	}
	data = append([]byte(xml.Header), data...)
	data = []byte(strings.Replace(string(data), "<project>", pomProjectHeader, 1))
	return data, err
}
