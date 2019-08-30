package pom

//go:generate go run gen/main.go gen/models.go gen/templates.go gen/build.go

import (
	"encoding/xml"
	"strings"
)

func Unmarshal(rawPom []byte) (project, error) {
	pom := project{}
	err := xml.Unmarshal(rawPom, &pom)
	return pom, err
}

var pomProjectHeader = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">`

func Marshal(pom project) ([]byte, error) {
	data, err := xml.MarshalIndent(pom, "", "    ")
	if err != nil {
		return data, err
	}
	data = append([]byte(xml.Header), data...)
	data = []byte(strings.Replace(string(data), "<project>", pomProjectHeader, 1))
	return data, err
}
