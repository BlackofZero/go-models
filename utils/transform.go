package utils

import "unsafe"

func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
