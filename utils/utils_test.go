package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"unicode"
)

func TestSanitizeWords(t *testing.T) {
	words := []string{
		"Can't",
		"Don't",
		"Isn't",
		"Won't",
		"I'm",
		"You'r",
		"They'r",
		"It's",
		"He's",
		"She's",
		"Cafe360",
		"Web2.0",
		"Room101",
		"3DPrinter",
		"IPv6",
	}

	for _, wd := range words {
		for _, r := range sanitizeWord(wd) {
			if !unicode.IsLetter(r) {
				t.Fatalf("Invalid Rune: %v\n", r)
			}
		}
		fmt.Println(sanitizeWord(wd))
	}
}

func TestCalcDocVec(t *testing.T) {
	f, err := os.Open("tests/oup_prc.txt")
	if err != nil {
		t.Fatal("Unable to open file")
		return
	}
	defer f.Close()

	InitMap(f)
	testFiles := []string{
		"tests/1.txt",
		"tests/2.txt",
		"tests/3.txt",
		"tests/4.txt",
	}
	for _, tf := range testFiles {
		tff, err := os.Open(tf)
		if err != nil {
			tff.Close()
			t.Fatalf("Unable to open file: %v\n", err)
		}
		vec := CalcDocVec(tff)
		log.Printf("%s ---- %v\n", tf, vec)
		tff.Close()
	}
}

func BenchmarkCalcDocVec(b *testing.B) {
	f, _ := os.Open("tests/oup_prc.txt")
	defer f.Close()
	InitMap(f)
	sample, _ := os.Open("tests/1.txt")
	defer sample.Close()

	for i := 0; i < b.N; i++ {
		data, _ := io.ReadAll(sample)
		CalcDocVec(bytes.NewBuffer(data))
		sample.Seek(0, io.SeekStart)
	}
}

func BenchmarkNCalcDocVec(b *testing.B) {
	f, _ := os.Open("tests/oup_prc.txt")
	defer f.Close()
	InitMap(f)
	sample, _ := os.Open("tests/1.txt")
	defer sample.Close()

	for i := 0; i < b.N; i++ {
		data, _ := io.ReadAll(sample)
		NCalcDocVec(bytes.NewBuffer(data))
		sample.Seek(0, io.SeekStart)
	}

}
