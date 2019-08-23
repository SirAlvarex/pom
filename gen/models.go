package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

type Schema struct {
	XMLName            xml.Name      `xml:"schema"`
	Text               string        `xml:",chardata"`
	Xs                 string        `xml:"xs,attr"`
	ElementFormDefault string        `xml:"elementFormDefault,attr"`
	Xmlns              string        `xml:"xmlns,attr"`
	TargetNamespace    string        `xml:"targetNamespace,attr"`
	Element            Element       `xml:"element"`
	ComplexType        []ComplexType `xml:"complexType"`
}

type Element struct {
	Text        string          `xml:",chardata"`
	Name        string          `xml:"name,attr"`
	Type        string          `xml:"type,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	Annotation  Annotation      `xml:"annotation"`
	Default     string          `xml:"default,attr"`
	ComplexType ElemComplexType `xml:"complexType"`
}

type SeqElement struct {
	Text      string `xml:",chardata"`
	Name      string `xml:"name,attr"`
	MinOccurs string `xml:"minOccurs,attr"`
	MaxOccurs string `xml:"maxOccurs,attr"`
	Type      string `xml:"type,attr"`
}

type Annotation struct {
	Text          string          `xml:",chardata"`
	Documentation []Documentation `xml:"documentation"`
}

type Documentation struct {
	Text   string `xml:",chardata"`
	Source string `xml:"source,attr"`
}

type ComplexType struct {
	Text       string      `xml:",chardata"`
	Name       string      `xml:"name,attr"`
	Type       string      `xml:"type,attr"`
	Annotation Annotation  `xml:"annotation"`
	All        All         `xml:"all"`
	Attribute  []Attribute `xml:"attribute"`
}

type ElemComplexType struct {
	Text     string   `xml:",chardata"`
	Sequence Sequence `xml:"sequence"`
}

type Sequence struct {
	Text    string     `xml:",chardata"`
	Element SeqElement `xml:"element"`
	Any     Any        `xml:"any"`
}
type Attribute struct {
	Text       string     `xml:",chardata"`
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	Use        string     `xml:"use,attr"`
	Annotation Annotation `xml:"annotation"`
}

type Any struct {
	Text            string `xml:",chardata"`
	MinOccurs       string `xml:"minOccurs,attr"`
	MaxOccurs       string `xml:"maxOccurs,attr"`
	ProcessContents string `xml:"processContents,attr"`
}

type All struct {
	Text    string    `xml:",chardata"`
	Element []Element `xml:"element"`
}

func (s Schema) FindType(target string) ComplexType {
	for _, complexType := range s.ComplexType {
		if complexType.Name == target {
			return complexType
		}
	}
	return ComplexType{}
}

func (s Schema) GetTypes() []string {
	result := make([]string, 0)
	for _, sType := range s.ComplexType {
		result = append(result, s.GetTypeAsString(sType))
	}
	return result
}

func (s Schema) GetTypeAsString(target ComplexType) string {
	format := `
{{ .TypeDoc }}
type {{ .Name }} struct {
	{{ if eq .Name "project" }}
	XMLName        xml.Name
	Xmlns          string   ` + "`xml:\"xmlns,attr\"`" + `
	Xsi            string   ` + "`xml:\"xsi,attr\"`" + `
	SchemaLocation string   ` + "`xml:\"schemaLocation,attr\"`" + `
	{{end}}

{{ range .Elem }}
    {{ . }}
{{ end }}
}
	`
	typeName := target.Name
	if typeName == "Model" {
		typeName = "project"
	}
	var typeDoc string
	if len(target.Annotation.Documentation) > 1 {
		doc := strings.Split(strings.Replace(strings.TrimSpace(target.Annotation.Documentation[1].Text), "\r\n", "\n", -1), "\n")
		typeDoc = fmt.Sprintf("\n// %s %s ", typeName, strings.Join(doc, "\n//"))
	}
	elements := make([]string, 0)
	for _, elem := range target.All.Element {
		valueToPrint := strings.Title(elem.Name)
		valueToPrint = strings.Replace(valueToPrint, "Url", "URL", -1)
		valueToPrint = strings.Replace(valueToPrint, "Id", "ID", -1)
		seqType := elem.ComplexType.Sequence.Element.Type
		seqName := elem.ComplexType.Sequence.Element.Name
		if len(seqType) > 0 {
			if strings.HasPrefix(seqType, "xs:") {
				if len(seqName) > 0 {
					valueToPrint += fmt.Sprintf(" *struct { Comment string `xml:\",comment\"`"+"\n%s *[]%s `xml:\"%s,omitempty\"` }",
						strings.Title(seqName),
						seqType,
						seqName,
					)

				} else {
					valueToPrint += " *[]" + seqType
				}
			} else {
				seqRune := []rune(seqType)
				seqRune[0] = unicode.ToLower(seqRune[0])
				seqLower := string(seqRune)
				valueToPrint += fmt.Sprintf(" *struct { Comment string `xml:\",comment\"`"+" \n%s *[]%s `xml:\"%s,omitempty\"`}",
					seqType,
					seqType,
					seqLower,
				)
			}
		}
		if len(elem.ComplexType.Sequence.Any.MaxOccurs) > 0 {
			valueToPrint += " *XMLAnyElement"
		}
		if len(elem.Type) > 0 {
			valueToPrint += " *" + elem.Type
		}

		valueToPrint += fmt.Sprintf(" `xml:\"%s,omitempty\"`", elem.Name)
		valueToPrint = strings.Replace(valueToPrint, "xs:", "", -1)
		valueToPrint = strings.Replace(valueToPrint, "boolean", "bool", -1)
		var documentation string
		if len(elem.Annotation.Documentation) > 1 {
			documentation = fmt.Sprintf("\n/* %s %s*/ ", strings.Title(elem.Name), strings.TrimSpace(elem.Annotation.Documentation[1].Text))
		}
		if len(documentation) > 0 {
			valueToPrint = fmt.Sprintf("%s\n%s", documentation, valueToPrint)
		}
		elements = append(elements, valueToPrint)
		//fmt.Println(strings.Join(level, ""), valueToPrint)
		//elemType := s.FindType(elem.ComplexType.Sequence.Element.Type)
		//s.PrintType(elemType, level...)
		//elemType = s.FindType(elem.Type)
		//s.PrintType(elemType, level...)
	}
	elements = append(elements, "Comment string `xml:\",comment\"`")
	buff := &bytes.Buffer{}
	tmpl, _ := template.New("tmp").Parse(format)
	tmpl.Execute(buff, struct {
		TypeDoc string
		Name    string
		Elem    []string
	}{
		typeDoc,
		typeName,
		elements,
	})
	return buff.String()
}
