package utils

import (
	"bytes"
	"encoding/gob"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// DecodeBlockLikeStruct decodes a byte array into a valid block response using gob
func DecodeBlockLikeStruct(anyVal []byte) (*types.BlockResponse, error) {
	var blockResponse *types.BlockResponse
	d := gob.NewDecoder(bytes.NewReader(anyVal))
	if err := d.Decode(&blockResponse); err != nil {
		return nil, err
	}

	return blockResponse, nil
}

// EncodeBlockLikeStruct encodes valid block response into a byte array using gob
func EncodeBlockLikeStruct(anyBlock *types.BlockResponse) ([]byte, error) {
	var buffer bytes.Buffer
	e := gob.NewEncoder(&buffer)
	if err := e.Encode(anyBlock); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
