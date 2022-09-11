package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"

	"github.com/vmihailenco/msgpack"
)

func encodeTicks(object interface{}, bCompress bool) ([]byte, error) {
	b, err := msgpack.Marshal(object)
	if err != nil {
		return nil, err
	}
	if bCompress {
		var zb bytes.Buffer
		zw := zlib.NewWriter(&zb)
		zw.Write(b)
		zw.Close()
		return zb.Bytes(), nil
	}
	return b, nil
}

func decodeTicks(b io.Reader, object interface{}, bCompress bool) error {
	var out bytes.Buffer
	if bCompress {
		reader, err := zlib.NewReader(b)
		if err != nil {
			return err
		}
		io.Copy(&out, reader)
		reader.Close()
	} else {
		io.Copy(&out, b)
	}
	return msgpack.Unmarshal(out.Bytes(), object)
}

func writeMsgpackFile(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		b, err := encodeTicks(object, true)
		if err == nil {
			file.Write(b)
		}
		file.Close()
	}
	return err
}

func readMsgpackFile(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	err = decodeTicks(file, object, true)
	file.Close()
	return err
}
