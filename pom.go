package pom

//go:generate go run gen/gen.go gen/models.go > gen_models.go

import "encoding/xml"

func Unmarshal(rawPom []byte) (project, error) {
	pom := project{}
	err := xml.Unmarshal(rawPom, &pom)
	return pom, err
}

func Marshal(pom project) ([]byte, error) {
	return xml.MarshalIndent(pom, " ", "    ")
}
