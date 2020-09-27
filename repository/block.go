package repository

import (
	"fmt"
	"log"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/go-redis/redis"
	"github.com/marpme/digibyte-rosetta-node/utils"
)

const IndexKey = "INDEX_LOCK"

type BlockRepository struct {
	rdb *redis.Client
}

func NewBlockRepository(rdb *redis.Client) *BlockRepository {
	repo := &BlockRepository{
		rdb: rdb,
	}
	count := rdb.DBSize().Val()
	exists := rdb.Exists(IndexKey).Val()
	if exists == 0 {
		err := repo.SetSyncedIndex(0)
		if err != nil {
			panic(err)
		}
		log.Println("initially set last synced block index to zero")
	}

	log.Printf("Initialized block repository. Loaded %v blocks", count)
	return repo
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

func (b *BlockRepository) GetLastSyncedIndex() (*int64, error) {
	val, err := b.rdb.Get(IndexKey).Int64()
	if err != nil {
		return nil, fmt.Errorf("couldn't load last block index from db")
	}

	return &val, nil
}

func (b *BlockRepository) SetSyncedIndex(index int64) error {

	err := b.rdb.Set(IndexKey, index, 0).Err()
	if err != nil {
		return fmt.Errorf("couldn't set last block index in db, cause %v", err)
	}

	return nil
}

func (b *BlockRepository) StoreSyncedBlock(keyHash string, block *types.BlockResponse, index int64) error {

	err := b.StoreBlock(keyHash, block)
	if err != nil {
		return err
	}

	err = b.SetSyncedIndex(index)
	if err != nil {
		return err
	}

	return nil
}
