package main

import (
	"encoding/xml"
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
