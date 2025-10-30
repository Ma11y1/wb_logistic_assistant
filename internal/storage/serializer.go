package storage

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type Serializer interface {
	Encode(w io.Writer, model interface{}) error
	Decode(r io.Reader, model interface{}) error
}

type JSONSerializer struct{}

func (JSONSerializer) Encode(w io.Writer, model interface{}) error {
	return json.NewEncoder(w).Encode(model)
}

func (JSONSerializer) Decode(r io.Reader, model interface{}) error {
	return json.NewDecoder(r).Decode(model)
}

type XMLSerializer struct{}

func (XMLSerializer) Encode(w io.Writer, model interface{}) error {
	return xml.NewEncoder(w).Encode(model)
}

func (XMLSerializer) Decode(r io.Reader, model interface{}) error {
	return xml.NewDecoder(r).Decode(model)
}
