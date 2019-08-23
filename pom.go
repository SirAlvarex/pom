package pom

//go:generate go run gen/gen.go gen/models.go

import (
	"encoding/xml"
)

func Unmarshal(rawPom []byte) (project, error) {
	pom := project{}
	err := xml.Unmarshal(rawPom, &pom)
	if err == nil && pom.Properties != nil {
		for index, prop := range pom.Properties.Elements {
			pom.Properties.Elements[index] = xmlMapEntry{
				xml.Name{Local: prop.XMLName.Local},
				prop.Value,
				prop.Comment,
			}
		}
	}
	return pom, err
}

func Marshal(pom project) ([]byte, error) {
	data, err := xml.MarshalIndent(pom, " ", "    ")
	data = append([]byte(xml.Header), data...)
	return data, err
}
