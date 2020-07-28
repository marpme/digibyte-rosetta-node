package repository

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/go-redis/redis"
	"github.com/marpme/digibyte-rosetta-node/utils"
)

type BlockRepository struct {
	rdb *redis.Client
}

func NewBlockRepository(rdb *redis.Client) *BlockRepository {
	return &BlockRepository{
		rdb: rdb,
	}
}

func (b *BlockRepository) StoreBlock(keyHash string, block *types.BlockResponse) error {

	encodedBlock, err := utils.EncodeBlockLikeStruct(block)
	if err != nil {
		panic(err)
	}

	err = b.rdb.Set(keyHash, encodedBlock, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *BlockRepository) GetBlock(keyHash string) (*types.BlockResponse, error) {
	val, err := b.rdb.Get(keyHash).Bytes()
	if err != nil {
		return nil, fmt.Errorf("couldn't load block from db with hash %s", keyHash)
	}

	decodedBlock, err := utils.DecodeBlockLikeStruct(val)
	if err != nil {
		return nil, fmt.Errorf("wasn't able to decode loaded block from db with hash %s", keyHash)
	}

	return decodedBlock, nil
}
