package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

const (
	BLOCK_SIZE_TAG = "maxBytes"
)

type ByteSize interface {
	int64 | uint64
}

func FromReader[DT any](reader io.Reader) (DT, error) {
	object := new(DT)

	var size int64
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return *object, err
	}

	var objectBytes []byte = make([]byte, int(size))
	if err := binary.Read(reader, binary.LittleEndian, &objectBytes); err != nil {
		return *object, err
	}

	objectBuffer := bytes.NewBuffer(objectBytes)
	decoder := gob.NewDecoder(objectBuffer)
	if err := decoder.Decode(object); err != nil {
		return *object, err
	}

	return *object, nil
}

func ToBytes(o any) ([]byte, error) {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Int, reflect.Int64, reflect.Uint64:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, o)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(o); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func ToBytesWithSize(o any) ([]byte, []byte, error) {
	oBytes, err := ToBytes(o)
	if err != nil {
		return nil, nil, err
	}
	oBytesSize := int64(len(oBytes))
	oBytesSizeBytes, err := ToBytes(oBytesSize)
	if err != nil {
		return nil, nil, err
	}
	return oBytes, oBytesSizeBytes, nil
}

func PadBytes(inputBytes []byte, blockSize int) ([]byte, error) {
	inputSize := len(inputBytes)

	if inputSize == blockSize {
		return inputBytes, nil
	}
	if inputSize > blockSize {
		return nil, fmt.Errorf("inputBytes %d is larger than blockSize %d", inputSize, blockSize)
	}

	padding := bytes.Repeat([]byte{byte('-')}, blockSize-inputSize)
	paddedBytes := append(inputBytes, padding...)

	return paddedBytes, nil
}

func GetBlockSize[T any]() (int, error) {
	var recDataSample T
	var blockSize int = 0
	recDataType := reflect.TypeOf(recDataSample)
	recDataTypeFields := reflect.VisibleFields(recDataType)
	for _, field := range recDataTypeFields {
		if field.Tag.Get(BLOCK_SIZE_TAG) != "" {
			fieldSize, err := strconv.Atoi(field.Tag.Get(BLOCK_SIZE_TAG))
			if err != nil {
				return 0, fmt.Errorf("can not read %v tag from %T: %v", BLOCK_SIZE_TAG, recDataSample, err)
			}
			blockSize += fieldSize
		} else {
			return 0, fmt.Errorf("all visible fields of %T must have the '%v' tag", recDataSample, BLOCK_SIZE_TAG)
		}

	}
	return blockSize, nil
}
