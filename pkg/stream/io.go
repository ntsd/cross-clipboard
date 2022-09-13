package stream

import (
	"bufio"
	"unsafe"
)

const dataSizeLength = 8 // int64 byte length

func readDataSize(r *bufio.Reader) (int, error) {
	bytes := make([]byte, dataSizeLength)
	for i := 0; i < dataSizeLength; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		bytes[i] = b
	}
	return bytesToInt(bytes), nil
}

func intToBytes(num int) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func bytesToInt(arr []byte) int {
	val := int(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}
