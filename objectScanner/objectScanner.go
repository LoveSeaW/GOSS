package main

import (
	"goss/apiServer/objects"
	"goss/pkg/es7"
	"goss/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
	}
}

func verify(hash string) {
	log.Println("verify:", hash)
	size, err := es7.SearchHashSize(hash)
	if err != nil {
		log.Println(err)
		return
	}
	stream, err := objects.GetStream(hash, size)
	if err != nil {
		log.Println(err)
		return
	}
	document := utils.CalculateHash(stream)
	if document != hash {
		log.Printf("object hash mismatch,calculated=%s,request=%s", document, hash)
	}
	stream.Close()
}
