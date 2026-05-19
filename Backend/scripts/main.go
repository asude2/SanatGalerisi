package main

import (
	"fmt"
	"os"
)

var connString = "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;TrustServerCertificate=true;"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run . <clean|massive|test>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "clean":
		seedCleanData()
	case "massive":
		seedMassiveData()
	case "test":
		seedTestData()
	default:
		fmt.Printf("Bilinmeyen seçenek: %s\n", os.Args[1])
		fmt.Println("Kullanım: go run . <clean|massive|test>")
		os.Exit(1)
	}
}
