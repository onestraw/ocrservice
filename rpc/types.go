package rpc

import (
	"bytes"
	"encoding/gob"
)

type OCRImageRequest struct {
	Lang      string
	Whitelist string
	Image     []byte
}

type OCRImageResponse struct {
	Version string
	Text    string
}

func Encode(object interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(object)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *OCRImageRequest) Encode() ([]byte, error) {
	return Encode(*r)
}

func (r *OCRImageRequest) Decode(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(r)
}

func (r *OCRImageResponse) Encode() ([]byte, error) {
	return Encode(*r)
}

func (r *OCRImageResponse) Decode(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(r)
}
