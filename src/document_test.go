package main

// func TestInsert(t *testing.T) {
// 	em, err := os.Open("tests/oup_prc.txt")
// 	if err != nil {
// 		t.Fatalf("Unable to open tests file: %v\n", err)
// 	}

// 	app := InitApp(em, DSN)
// 	defer app.DbPool.Close()

// 	testFiles := []string{
// 		"tests/testInsert.txt",
// 	}

// 	for _, tf := range testFiles {
// 		f, err := os.Open(tf)
// 		if err != nil {
// 			t.Fatalf("Unable to Open file: %v\n", err)
// 		}
// 		defer f.Close()
// 		if err := app.Insert(f); err != nil {
// 			t.Fatalf("Unable to Insert file to DB: %v\n", err)
// 			return
// 		}

// 	}

// }
