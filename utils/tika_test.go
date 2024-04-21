package utils_test

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"terch/utils"
	"testing"
)

func TestCreateDocument(t *testing.T) {
	var TEST_FILE string = "tests/bgnet.pdf"
	file, err := os.Open(TEST_FILE)
	if err != nil {
		t.Fatalf("unable to open test file: %v\n", err)
		return
	}
	defer file.Close()

	id := rand.Intn(100)

	fname := fmt.Sprintf("tests/%s%s", strconv.Itoa(id), filepath.Ext(file.Name()))

	ff, err := os.Create(fname)
	if err != nil {
		ff.Close()
		t.Fatalf("Unable to open new file: %v\n", err)
		return
	}

	io.Copy(ff, file)

	ff.Close()

	nff, err := os.Open(fname)
	if err != nil {
		t.Fatalf("Unable to open new file: %v\n", err)
		return
	}
	defer nff.Close()

	wvf, err := os.Open("tests/oup_prc.txt")
	if err != nil {
		t.Fatalf("Unable to open embeddings file: %v\n", err)
		return
	}
	defer wvf.Close()

	cli := utils.InitTikaClient("http://localhost:9998")
	utils.InitMap(wvf)

	doc, err := utils.CreateDocument(nff, cli, file.Name())
	if err != nil {
		t.Fatalf("Unable to create document from file: %v\n", err)
		return
	}

	fmt.Println(doc.Vec)
}
