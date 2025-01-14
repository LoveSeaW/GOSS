package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetOffsetFromHeader(header http.Header) int64 {
	byteRange := header.Get("range")
	if len(byteRange) < 7 {
		return 0
	}

	// 是否以 bytes= 开头
	if !strings.HasPrefix(byteRange, "bytes=") {
		return 0
	}

	bytePos := strings.Split(byteRange[6:], "_")
	offset, err := strconv.ParseInt(bytePos[0], 10, 64)
	if err != nil {
		return 0
	}
	return offset
}

// 获取对象散列值的Base64编码
func GetHashFromHeader(header http.Header) string {
	digest := header.Get("digest")
	if !strings.HasPrefix(digest, "SHA-256=") {
		return ""
	}
	return digest[8:]
}

// 获取对象数据的长度
func GetSizeFromHeader(header http.Header) int64 {
	sizeStr := header.Get("content-length")
	if sizeStr == "" {
		return 0
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0
	}
	return size
}

func CalculateHash(reader io.Reader) string {
	hash := sha256.New()
	_, err := io.Copy(hash, reader)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
