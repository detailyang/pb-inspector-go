package inspector

import (
	"io"
)

// Inspector represents the protobuf inspector
type Inspector struct {
	decoder    *Decoder
	definition *Definition
}

// NewInspector returns the inspector to inpsect protobuf
func NewInspector() *Inspector {
	return &Inspector{
		definition: NewDefinition(),
	}
}

// ReadSchemaFromReader reads schema from reader which named f.
func (p *Inspector) ReadSchemaFromReader(f string, r io.Reader) error {
	return p.definition.ReadFrom(f, r)
}

// ReadSchemaFromFile reads scheam from file.
func (p *Inspector) ReadSchemaFromFile(f string) error {
	return p.definition.ReadFile(f)
}

// ToMapWithSchema maps raw bytes to map[string]interface{} by self definition
func (p *Inspector) ToMapWithSchema(pkg, name string, raw []byte) (map[string]interface{}, error) {
	return NewDecoder(p.definition, NewBuffer(raw)).Decode(pkg, name)
}

// ToMapWithSchema maps raw bytes to map[string]interface{} by specified definition
func (p *Inspector) ToMapWithSchemaByDefinition(d *Definition, pkg, name string, raw []byte) (map[string]interface{}, error) {
    return NewDecoder(d, NewBuffer(raw)).Decode(pkg, name)
}

// InspectWithoutSchema inspects raw protobuf binary data and write to w
func (p *Inspector) InspectWithoutSchema(verbose bool, raw []byte, w io.Writer) error {
	return NewBuffer(raw).InspectWithoutSchema(verbose, raw, w)
}
