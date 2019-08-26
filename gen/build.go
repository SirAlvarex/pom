package main

import (
	"bytes"
	"fmt"
	"strings"
)

var existingTypes = make(map[string]bool, 0)

// GetTypes returns the formatted struct definitions of each type
func (s Schema) GetTypes() []string {
	result := make([]string, 0)
	for _, sType := range s.ComplexType {
		result = append(result, s.GetTypeAsString(sType))
	}
	return result
}

type pomType struct {
	Name   string
	Doc    string
	Fields []pomTypeField
}

type pomTypeField struct {
	Name         string
	Doc          string
	Tag          string
	Type         string
	DefaultValue string
	IsPointer    bool
	IsSlice      bool
}

// GetTypeAsString applies a type to a struct template
func (s Schema) GetTypeAsString(target ComplexType) string {
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

	types := make([]pomType, 0)
	myType := pomType{
		Name: typeName,
		Doc:  typeDoc,
	}
	// fields will be the fields in the struct
	fields := make([]string, 0)
	myType.Fields = make([]pomTypeField, 0)
	for _, elem := range target.All.Element {
		abc := pomTypeField{}
		// Time to clean up the field name
		field := strings.Title(elem.Name)
		// GoLint spec
		field = strings.Replace(field, "Url", "URL", -1)
		// GoLint spec
		field = strings.Replace(field, "Id", "ID", -1)
		abc.Name = field
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
			subTypeName := fmt.Sprintf("Sequence%s", strings.Title(seqName))
			if ok := existingTypes[subTypeName]; !ok {
				subTypeType := strings.Replace(strings.Replace(seqType, "xs:", "", -1), "boolean", "bool", -1)
				subTypeDefault := fmt.Sprintf("%s{}", subTypeType)
				subType := pomType{
					Name: subTypeName,
					Doc:  "// Sequence Type",
					Fields: []pomTypeField{
						pomTypeField{
							Name: "Comment",
							Type: "string",
							Tag:  "`xml:\",comment\"`",
						},
						pomTypeField{
							Name:         strings.Title(seqName),
							Type:         subTypeType,
							Tag:          fmt.Sprintf("`xml:\"%s,omitempty\"`", seqName),
							IsPointer:    true,
							IsSlice:      true,
							DefaultValue: subTypeDefault,
						},
					},
				}
				types = append(types, subType)
				existingTypes[subTypeName] = true
			}
			abc.Type = subTypeName
			abc.IsPointer = true
			abc.DefaultValue = fmt.Sprintf("%s{}", subTypeName)
		}

		// If MaxOccurs is set, then that means Any is set.
		// An "Any" element is XMLs type of Generic
		// XMLInner is the only way we can do generics -- except that means we cannot modify the subxml
		// XMLAnyElement, however, is like a map[string]string, but ordered
		// Properties has a consistent map-like format, so we have a special case there
		if len(elem.ComplexType.Sequence.Any.MaxOccurs) > 0 {
			if elem.Name == "properties" {
				field += " *XMLAnyElement"
				abc.Type = "XMLAnyElement"
				abc.DefaultValue = "XMLAnyElement{}"
				abc.IsPointer = true
			} else {
				field += " *XMLInner"
				abc.Type = "XMLInner"
				abc.DefaultValue = "XMLInner{}"
				abc.IsPointer = true
			}
		}

		// If the element itself has a type, set it here.
		// This value is unset if the type is a sequence, so no conflict with values above
		if len(elem.Type) > 0 {
			field += " *" + elem.Type
			abc.IsPointer = true
			abc.Type = strings.Replace(strings.Replace(elem.Type, "xs:", "", -1), "boolean", "bool", -1)
			abc.DefaultValue = fmt.Sprintf("%s{}", abc.Type)
			if abc.Type == "bool" {
				abc.DefaultValue = "false"
			} else if abc.Type == "string" {
				abc.DefaultValue = `""`
			}
		}

		abc.Tag = fmt.Sprintf(" `xml:\"%s,omitempty\"`", elem.Name)
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
			abc.Doc = documentation
			field = fmt.Sprintf("%s\n%s", documentation, field)
		}
		myType.Fields = append(myType.Fields, abc)
		fields = append(fields, field)
	}

	// Add a comment field to the bottom of each subtype.  This way we keep comments
	fields = append(fields, "Comment string `xml:\",comment\"`")
	myType.Fields = append(myType.Fields, pomTypeField{
		Name: "Comment",
		Type: "string",
		Tag:  "`xml:\",comment\"`",
	})
	types = append(types, myType)

	// Parsing template out to a buffer
	buff := &bytes.Buffer{}
	/*
		structFormat.Execute(buff, struct {
			TypeDoc string
			Name    string
			Fields  []string
		}{
			typeDoc,
			typeName,
			fields,
		})
	*/
	structFormatv2.Execute(buff, types)

	// Return the string representation of this struct definition
	return buff.String()
}
