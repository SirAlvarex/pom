package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"text/template"
)

// Schema is the root element of the POM XSD
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

// Element is a description of value inside of a type
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

// SeqElement is the definition of an element inside a sequence
type SeqElement struct {
	Text      string `xml:",chardata"`
	Name      string `xml:"name,attr"`
	MinOccurs string `xml:"minOccurs,attr"`
	MaxOccurs string `xml:"maxOccurs,attr"`
	Type      string `xml:"type,attr"`
}

// Annotation contains documentation information
type Annotation struct {
	Text          string          `xml:",chardata"`
	Documentation []Documentation `xml:"documentation"`
}

// Documentation is...documentation
type Documentation struct {
	Text   string `xml:",chardata"`
	Source string `xml:"source,attr"`
}

// ComplexType is the unique types declared in the XSD
type ComplexType struct {
	Text       string      `xml:",chardata"`
	Name       string      `xml:"name,attr"`
	Type       string      `xml:"type,attr"`
	Annotation Annotation  `xml:"annotation"`
	All        All         `xml:"all"`
	Attribute  []Attribute `xml:"attribute"`
}

// ElemComplexType contains type information if this element is part of a list
type ElemComplexType struct {
	Text     string   `xml:",chardata"`
	Sequence Sequence `xml:"sequence"`
}

// Sequence described how a sequence is configured
type Sequence struct {
	Text    string     `xml:",chardata"`
	Element SeqElement `xml:"element"`
	Any     Any        `xml:"any"`
}

// Attribute is XML attributes
type Attribute struct {
	Text       string     `xml:",chardata"`
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	Use        string     `xml:"use,attr"`
	Annotation Annotation `xml:"annotation"`
}

// Any is the generic untyped type of XML
type Any struct {
	Text            string `xml:",chardata"`
	MinOccurs       string `xml:"minOccurs,attr"`
	MaxOccurs       string `xml:"maxOccurs,attr"`
	ProcessContents string `xml:"processContents,attr"`
}

// All is the list of elements in a type
type All struct {
	Text    string    `xml:",chardata"`
	Element []Element `xml:"element"`
}

// GetTypes returns the formatted struct definitions of each type
func (s Schema) GetTypes() []string {
	result := make([]string, 0)
	for _, sType := range s.ComplexType {
		result = append(result, s.GetTypeAsString(sType))
	}
	return result
}

// GetTypeAsString applies a type to a struct template
func (s Schema) GetTypeAsString(target ComplexType) string {
	format := `
{{ .TypeDoc }}
type {{ .Name }} struct {
{{ range .Fields }}
    {{ . }}
{{ end }}
}
	`
	typeName := target.Name

	// We want the root object to be a `project`, not a `Model`
	if typeName == "Model" {
		typeName = "project"
	}

	// Format the Type documentation string.
	// Type declarations must use //, so we remove newlines and smush things together
	// For the Maven XSD, the first element in a doc is the version of the pom it was added.  So we take just the second element
	var typeDoc string
	if len(target.Annotation.Documentation) > 1 {
		doc := strings.Split(strings.Replace(strings.TrimSpace(target.Annotation.Documentation[1].Text), "\r\n", "\n", -1), "\n")
		typeDoc = fmt.Sprintf("\n// %s %s ", typeName, strings.Join(doc, "\n//"))
	}

	// fields will be the fields in the struct
	fields := make([]string, 0)
	for _, elem := range target.All.Element {
		// Time to clean up the field name
		field := strings.Title(elem.Name)
		// GoLint spec
		field = strings.Replace(field, "Url", "URL", -1)
		// GoLint spec
		field = strings.Replace(field, "Id", "ID", -1)

		// Sequence is set if the this type is a list of elements
		seqType := elem.ComplexType.Sequence.Element.Type
		seqName := elem.ComplexType.Sequence.Element.Name
		if len(seqType) > 0 {
			// Converting these types to work with XML
			// <models>
			//    <model>thing<model>
			// </models>
			// For the <model> tag to work, we need to create a subelement struct
			field += fmt.Sprintf(" *struct { Comment string `xml:\",comment\"`"+"\n%s []*%s `xml:\"%s,omitempty\"` }",
				strings.Title(seqName),
				seqType,
				seqName,
			)
		}

		// If MaxOccurs is set, then that means Any is set.
		// An "Any" element is XMLs type of Generic
		// XMLInner is the only way we can do generics -- except that means we cannot modify the subxml
		// XMLAnyElement, however, is like a map[string]string, but ordered
		// Properties has a consistent map-like format, so we have a special case there
		if len(elem.ComplexType.Sequence.Any.MaxOccurs) > 0 {
			if elem.Name == "properties" {
				field += " *XMLAnyElement"
			} else {
				field += " *XMLInner"
			}
		}

		// If the element itself has a type, set it here.
		// This value is unset if the type is a sequence, so no conflict with values above
		if len(elem.Type) > 0 {
			field += " *" + elem.Type
		}

		// Adding the XML tags to the end of the field
		field += fmt.Sprintf(" `xml:\"%s,omitempty\"`", elem.Name)
		// Removing xs:, which is an xml standard type (string, boolean)
		field = strings.Replace(field, "xs:", "", -1)
		// Rename boolean to golang type bool
		field = strings.Replace(field, "boolean", "bool", -1)

		// Format the documentation for the field
		// For the Maven XSD, the first element in a doc is the version of the pom it was added.  So we take just the second element
		var documentation string
		if len(elem.Annotation.Documentation) > 1 {
			documentation = fmt.Sprintf("\n/* %s %s*/ ", strings.Title(elem.Name), strings.TrimSpace(elem.Annotation.Documentation[1].Text))
		}
		// Only add a documentation line if we do in fact have docs
		if len(documentation) > 0 {
			field = fmt.Sprintf("%s\n%s", documentation, field)
		}
		fields = append(fields, field)
	}
	// Add a comment field to the bottom of each subtype.  This way we keep comments
	fields = append(fields, "Comment string `xml:\",comment\"`")

	// Parsing template out to a buffer
	buff := &bytes.Buffer{}
	tmpl, _ := template.New("tmp").Parse(format)
	tmpl.Execute(buff, struct {
		TypeDoc string
		Name    string
		Fields  []string
	}{
		typeDoc,
		typeName,
		fields,
	})

	// Return the string representation of this struct definition
	return buff.String()
}
