package routines

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marpme/digibyte-rosetta-node/client"
	"github.com/marpme/digibyte-rosetta-node/repository"
	"github.com/marpme/digibyte-rosetta-node/utils"
)

var ctx = context.Background()

func StartSyncingCapabilities(wg *sync.WaitGroup, blockRepo *repository.BlockRepository, client client.DigibyteClient) {
	index, err := blockRepo.GetLastSyncedIndex()
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("tip is becoming stale, scanning for new blocks (last known index: %v)", *index)
		lastBlock, err := client.GetLatestBlock(ctx)
		if err != nil {
			log.Panicf("The sync process panic'd unexpectly %s", err)
		}
		for *index < lastBlock.Height {
			block, err := client.GetBlock(ctx, *index)
			err = blockRepo.StoreSyncedBlock(block.Hash, utils.MapBlockWithTransaction(block), block.Height)
			if err != nil {
				log.Panicf("The sync process panic'd unexpectly %s", err)
			}

			if *index%100 == 0 {
				log.Printf("synced 100 blocks successfully - current block (%v, %s)\n", block.Height, block.Hash)
			}
			*index++
		}
		time.Sleep(15 * time.Second) // wait at least a block time.
	}
}
