package filestorage

import (
	"fmt"
	"io"
	"os"

	"github.com/prondos/axdb/pkg/utils"
)

func (p *FileStorage[IT, DT]) indexToBytes(index IT, offset int64) []byte {
	var indexWrapper = IndexWrapper[IT]{I: index, O: offset}
	indexBytes, indexBytesSizeBytes, err := utils.ToBytesWithSize(indexWrapper)
	if err != nil {
		panic(err)
	}
	return append(indexBytesSizeBytes, indexBytes...)
}

func (p *FileStorage[IT, DT]) indexFromReader(reader *os.File) (IT, int64, error) {
	indexWrapper, err := utils.FromReader[IndexWrapper[IT]](reader)
	if err == io.EOF {
		return *new(IT), 0, err
	}
	if err != nil {
		panic(err)
	}
	return indexWrapper.I, indexWrapper.O, nil
}

func (p *FileStorage[IT, DT]) dataToBytes(data DT, dataBlockSize int) []byte {
	var dataWrapper = &DataWrapper[DT]{D: data}
	dataBytes, dataBytesSizeBytes, err := utils.ToBytesWithSize(*dataWrapper)
	if err != nil {
		panic(err)
	}
	dataBytesPadded, err := utils.PadBytes(dataBytes, dataBlockSize)
	if err != nil {
		panic(err)
	}
	return append(dataBytesSizeBytes, dataBytesPadded...)
}

func (p *FileStorage[IT, DT]) dataFromReader(reader *os.File, offset int64) (DT, error) {
	if _, err := reader.Seek(offset, io.SeekStart); err != nil {
		if err == io.EOF {
			return *new(DT), err
		}
		panic(fmt.Sprintf("error seeking to offset %v", offset))
	}

	dataWrapper, err := utils.FromReader[DataWrapper[DT]](reader)
	if err != nil {
		panic(err)
	}

	return dataWrapper.D, nil
}
