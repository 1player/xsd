package xsd

import (
	"encoding/xml"
	"io"
)

type Type interface{}

type Annotation struct {
	Documentation string `xml:"http://www.w3.org/2001/XMLSchema documentation"`
}

type Element struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Sequence struct {
	Elements []Element `xml:"http://www.w3.org/2001/XMLSchema element"`
}

type ComplexType struct {
	Name       string      `xml:"name,attr"`
	Annotation *Annotation `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	Sequence   *Sequence   `xml:"http://www.w3.org/2001/XMLSchema sequence"`
}

type Restriction struct {
	Base string `xml:"base,attr"`
}

type SimpleType struct {
	Name        string      `xml:"name,attr"`
	Annotation  *Annotation `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	Restriction Restriction `xml:"http://www.w3.org/2001/XMLSchema restriction"`
}

type Document struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2001/XMLSchema schema"`
	TargetNamespace string   `xml:"targetNamespace,attr"`
	Version         string   `xml:"version,attr"`

	Annotation   *Annotation   `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	SimpleTypes  []SimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	ComplexTypes []ComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Elements     []Element     `xml:"http://www.w3.org/2001/XMLSchema element"`

	lookup map[string]Type
}

func (d *Document) fillLookupTable() {
	d.lookup = make(map[string]Type)

	// Simple types
	for _, t := range d.SimpleTypes {
		d.lookup[t.Name] = t
	}

	// Complex types
	for _, t := range d.ComplexTypes {
		d.lookup[t.Name] = t
	}

	// All remaining elements
	for _, t := range d.Elements {
		d.lookup[t.Name] = t
	}
}

// Parse and return an XSD document
func NewFromReader(r io.Reader) (*Document, error) {
	var d Document

	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&d)
	if err != nil {
		return nil, err
	}

	d.fillLookupTable()
	return &d, nil
}

func (d *Document) Lookup(element string) Type {
	if t, ok := d.lookup[element]; ok {
		return t
	}
	return nil
}
