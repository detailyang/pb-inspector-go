package inspector

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	pp "github.com/emicklei/proto"
)

// Definition holds the protobuf definition
type definition struct {
	messages          map[string]*pp.Message
	enums             map[string]*pp.Enum
	filenamesRead     []string
	filenameToPackage map[string]string
}

func newDefinition() *definition {
	return &definition{
		messages:          map[string]*pp.Message{},
		enums:             map[string]*pp.Enum{},
		filenamesRead:     []string{},
		filenameToPackage: map[string]string{},
	}
}

// ReadFile reads the proto definition from a filename.
func (d *definition) ReadFile(filename string) error {
	fileReader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	return d.ReadFrom(filename, fileReader)
}

// ReadFrom reads from reader which named filename
func (d *definition) ReadFrom(filename string, reader io.Reader) error {
	for _, each := range d.filenamesRead {
		if each == filename {
			return nil
		}
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	d.filenamesRead = append(d.filenamesRead, filename)
	parser := pp.NewParser(bytes.NewReader(data))
	def, err := parser.Parse()
	if err != nil {
		return err
	}

	pkg := packageOf(def)
	d.filenameToPackage[filename] = pkg
	pp.Walk(def, pp.WithMessage(func(each *pp.Message) {
		d.AddMessage(pkg, each.Name, each)
	}))
	pp.Walk(def, pp.WithEnum(func(each *pp.Enum) {
		d.AddEnum(pkg, each.Name, each)
	}))
	return nil
}

// Package returns the proto package name as declared in the proto filename.
func (d *definition) Package(filename string) (pkg string, ok bool) {
	pkg, ok = d.filenameToPackage[filename]
	return
}

// MessagesInPackage returns the messages
func (d *definition) MessagesInPackage(pkg string) (list []*pp.Message) {
	for k, v := range d.messages {
		if strings.HasPrefix(k, pkg+".") {
			list = append(list, v)
		}
	}
	return
}

// Message returns the message
func (d *definition) Message(pkg string, name string) (m *pp.Message, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	m, ok = d.messages[key]
	return
}

// Enum returns the enum
func (d *definition) Enum(pkg string, name string) (e *pp.Enum, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	e, ok = d.enums[key]
	return
}

// AddEnum adds the Enum
func (d *definition) AddEnum(pkg string, name string, enu *pp.Enum) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.enums[key] = enu
}

// AddMessage adds the message
func (d *definition) AddMessage(pkg string, name string, message *pp.Message) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.messages[key] = message
}

func packageOf(def *pp.Proto) string {
	for _, each := range def.Elements {
		if p, ok := each.(*pp.Package); ok {
			return p.Name
		}
	}
	return ""
}
