package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"terch/utils/xml"

	"github.com/google/go-tika/tika"
)

type Document struct {
	Name string
	Vec  []float64
	buf  *bytes.Buffer
}

type Pdata struct {
	Data string `xml:",chardata"`
}

func InitTikaClient(url string) *tika.Client {
	return tika.NewClient(nil, url)
}

// TODO: Remove 'uploads/' from File names
func CreateDocument(file *os.File, cli *tika.Client) (*Document, error) {
	doc := &Document{}
	doc.buf = bytes.NewBuffer(nil)
	doc.Vec = make([]float64, 10)
	doc.Name = file.Name()

	res, err := cli.Parse(context.Background(), file)
	if err != nil {
		fmt.Printf("Unable to parse file: %v\n", err)
		return nil, err
	}
	d := xml.NewDecoder(strings.NewReader(res))
	d.Strict = false
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "p" {
				var pdata Pdata
				if err := d.DecodeElement(&pdata, &t); err != nil {
					fmt.Printf("unable to decode element: %v\n", err)
					return nil, err
				}
				if _, err := doc.buf.WriteString(pdata.Data); err != nil {
					fmt.Printf("Unable to write data to buffer: %v\n", err)
					return nil, err
				}
			}
		default:
		}
	}

	vec := NCalcDocVec(doc.buf)
	copy(doc.Vec, vec)
	return doc, nil
}
