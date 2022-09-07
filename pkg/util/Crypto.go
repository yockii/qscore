package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/forgoer/openssl"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Des3EcbPkcs5PaddingEncryptBase64String(src string, key []byte) (string, error) {
	rst, err := openssl.Des3ECBEncrypt([]byte(src), key, openssl.PKCS5_PADDING)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(rst), nil
}

func Des3EcbPkcs5PaddingDecryptBase64String(src string, key []byte) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	rst, err := openssl.Des3ECBDecrypt(bs, key, openssl.PKCS5_PADDING)
	if err != nil {
		return "", err
	}
	return string(rst), nil
}
