package purego

import (
	"fmt"
	"reflect"
)

type primitiveType int

const (
	ptFloat  primitiveType = FLOAT
	ptInt    primitiveType = INT
	ptUint   primitiveType = UNSIGNED_INT
	ptUshort primitiveType = UNSIGNED_SHORT
	ptUbyte  primitiveType = UNSIGNED_BYTE
)

func (p primitiveType) sizeInBytes() int {
	switch p {
	case ptFloat:
		return 4
	case ptInt:
		return 4
	case ptUint:
		return 4
	case ptUshort:
		return 2
	case ptUbyte:
		return 1
	default:
		panic(fmt.Errorf("unknown primitiveType 0x%.4x", p))
	}
}

func (p primitiveType) isArrayOfType(array interface{}) bool {
	ty := reflect.TypeOf(array).Elem()
	switch p {
	case ptFloat:
		return ty.Name() == "float32"
	case ptInt:
		return ty.Name() == "int32"
	case ptUint:
		return ty.Name() == "uint32"
	case ptUshort:
		return ty.Name() == "uint16"
	case ptUbyte:
		return ty.Name() == "uint8"
	default:
		panic(fmt.Errorf("unknown primitiveType 0x%.4x", p))
	}
}
