package main

// echo 08ffff01100840f7d438 | xxd -r -p | bin/pb-inspector -

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	inspector "github.com/detailyang/pb-inspector-go/protobuf-inspector"
	"github.com/k0kubun/pp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pb-inspector"
	app.Description = "pb-protobuf inspects protobuf binary file"
	app.Usage = "a protobuf inspector"
	app.UsageText = "pb-inspector [file] <package> <name>"
	app.Author = "detailyang"
	app.Email = "detailyang@gmail.com"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file-type",
			Value: "binary",
			Usage: "Specify the file type (binary or hex)",
		},
		cli.StringSliceFlag{
			Name:  "pb-file",
			Value: nil,
			Usage: "Load pb file to decode",
		},
		cli.StringFlag{
			Name:  "pb-dir",
			Value: "",
			Usage: "Load *.pb *.proto from directory",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}
		return run(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	file := c.Args().First()
	pkg := c.Args().Get(1)
	name := c.Args().Get(2)

	var r *bufio.Reader

	if file == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		r = bufio.NewReader(f)
	}

	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	filetype := c.String("file-type")
	switch filetype {
	case "binary":
		break
	case "hex":
		var err error
		raw, err = hex.DecodeString(strings.TrimSpace(string(raw)))
		if err != nil {
			return err
		}
	default:
		return errors.New("unknow file type")
	}

	pbfiles := c.StringSlice("pb-file")
	pbdir := c.String("pb-dir")
	pbfiles, err = walkpb(pbdir, pbfiles)
	if err != nil {
		return err
	}

	w := bytes.NewBuffer(nil)
	in := inspector.NewInspector()

	if len(pbfiles) == 0 {
		if err := in.InspectWithoutSchema(false, raw, w); err != nil {
			return err
		}

		fmt.Println(hex.EncodeToString(raw))
		fmt.Println("=>")
		fmt.Println(w.String())

	} else {
		for _, file := range pbfiles {
			if err = in.ReadSchemaFromFile(file); err != nil {
				return err
			}
		}

		m, err := in.ToMapWithSchema(pkg, name, raw)
		if err != nil {
			return err
		}

		fmt.Println(hex.EncodeToString(raw))
		fmt.Printf("=> (pkg=%s name=%s)\n", pkg, name)
		pp.Println(m)
	}

	return nil
}

func walkpb(dir string, files []string) ([]string, error) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pb") || strings.HasSuffix(path, ".proto") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
