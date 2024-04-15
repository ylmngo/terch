package utils_test

import (
	"fmt"
	"os"
	"terch/utils"
	"testing"
)

func TestCreateDocument(t *testing.T) {
	var TEST_FILE string = "bgnet.pdf"
	file, err := os.Open(TEST_FILE)
	if err != nil {
		t.Fatalf("unable to open test file: %v\n", err)
		return
	}
	defer file.Close()

	wvf, err := os.Open("oup_prc.txt")
	if err != nil {
		t.Fatalf("Unable to open embeddings file: %v\n", err)
		return
	}
	defer wvf.Close()

	cli := utils.InitTikaClient("http://localhost:9998")
	utils.InitMap(wvf)

	doc, err := utils.CreateDocument(file, cli)
	if err != nil {
		t.Fatalf("Unable to create document from file: %v\n", err)
		return
	}

	fmt.Println(doc.Vec)
}
