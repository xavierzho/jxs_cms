package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func EncodeMD5Hex(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func EncodeMD5(text string) string {
	w := md5.New()
	io.WriteString(w, text)                  //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}
