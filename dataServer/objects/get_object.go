package objects

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"goss/dataServer/locate"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getFile(name string) string {
	fmt.Println("objects try getFile:", name)
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}

	file := files[0]
	hash := sha256.New()
	sendFile(hash, file)

	calculated := url.PathEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	hash2 := strings.Split(file, ".")[2]
	if calculated != hash2 {
		log.Println("object hash mismatch, removed,", file)
		locate.Del(file)
		os.Remove(file)
		return ""
	}
	return file
}
