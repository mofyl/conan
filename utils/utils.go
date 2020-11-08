package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
)

const (
	PRO_NAME       = "conan"
	MM_PRIVATE_KEY = "32ed87bdb5fdc5e9cba88547376818d4"
	MM_SERVER_URL  = "http://127.0.0.1:8080"
)

var (
	Json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func GetProDirAbs() string {

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic("Can not Read GOPATH")
	}
	sbuilder := strings.Builder{}
	sbuilder.WriteString(goPath)
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString("src")
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString(PRO_NAME)
	sbuilder.WriteString(string(filepath.Separator))

	path := sbuilder.String()
	return path
}

func GetUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func AESCBCEncode(key string, data []byte) ([]byte, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return nil, err
	}
	content := pKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyb[:block.BlockSize()])
	crypted := make([]byte, len(content))
	blockMode.CryptBlocks(crypted, content)
	return crypted, nil
}

func AESCBCDecode(key string, data []byte) ([]byte, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyb[:blockSize])
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)

	origData = pKCS5UnPadding(origData)

	return origData, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

func AESCBCBase64Encode(key string, data []byte) (string, error) {
	crypted, err := AESCBCEncode(key, data)

	if err != nil {
		return "", err
	}

	return Base64Encode(crypted), nil
}

func AESCBCBase64Decode(key string, data string) ([]byte, error) {
	crypted, err := Base64Decode(data)

	if err != nil {
		return nil, err
	}

	orig, err := AESCBCDecode(key, crypted)
	if err != nil {
		return nil, err
	}
	return orig, nil
}

func StrUUID() (string, error) {

	uuid, err := GetUUID()
	if err != nil {
		return "", nil
	}

	uuidB, _ := uuid.MarshalText()
	res := md5.Sum(uuidB)
	return hex.EncodeToString(res[:]), nil
}

func CheckUrl(url string) bool {
	//b, _ := regexp.MatchString("^http:////[a-zA-Z0-9](/.[a-z]*)|https:////[a-zA-Z0-9](/.[a-z]*)|[1-9]{1,3}(/.[0-9]{1,3}){3}$", url)
	b, _ := regexp.MatchString(`^http:\/\/www\.[a-zA-Z0-9]+(\.[a-z]+)+$|^https:\/\/www\.[a-zA-Z0-9]+(\.[a-z]+)+$|^[1-9]{1,3}(\.[0-9]{1,3}){3}$`, url)
	return b
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("str can not empty")
	}

	return str[len(str)-1]

}

func JoinPath(absPath string, relativePath string) string {
	if relativePath == "" {
		return absPath
	}

	finalPath := path.Join(absPath, relativePath)

	appendSlash := lastChar(finalPath) != '/' && lastChar(relativePath) == '/'

	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func RemoveSpace(str string) string {
	return strings.TrimSpace(str)
}

func StrToByte(str string) []byte {
	var res []byte

	*(*string)(unsafe.Pointer(&res)) = str

	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&res)) + 2*unsafe.Sizeof(&res))) = len(str)

	return res
}

func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BCryptEnCodePwd(data []byte) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return hash, nil
}

func BCryptCheckPwd(hash []byte, pwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, pwd)

	return err == nil
}
