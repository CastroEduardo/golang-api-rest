package util

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
func DecodeMD5(s string) string {
	//src := []byte(s)
	// We can use the source slice itself as the destination
	// because the decode loop increments by one and then the 'seen' byte is not used anymore.
	str, _ := hex.DecodeString(s)
	return string([]byte(str)) //string([]byte(str))
}

// // EncodeMD5 md5 encryption
// func DecodeMD5(value string) string {
// 	m := md5.New()
// 	m.Write([]byte(value))

// 	return hex.DecodeString(value)
// }
