package services

import (
	"context"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/marpme/digibyte-rosetta-node/client"
	"github.com/marpme/digibyte-rosetta-node/repository"
)

// BlockAPIService client based implementation of the block servicer
type blockAPIService struct {
	server.BlockAPIServicer
	client          client.DigibyteClient
	blockRepository *repository.BlockRepository
}

// NewBlockAPIService creates a new block API service
func NewBlockAPIService(client client.DigibyteClient, blockRepository *repository.BlockRepository) server.BlockAPIServicer {
	return &blockAPIService{
		client:          client,
		blockRepository: blockRepository,
	}
}

func mapTransactions(txs []string) []*types.TransactionIdentifier {
	var transactionsIdentifiers []*types.TransactionIdentifier
	for i := 0; i < len(txs); i++ {
		transactionsIdentifiers = append(transactionsIdentifiers, &types.TransactionIdentifier{
			Hash: txs[i],
		})
	}

	return transactionsIdentifiers
}

func (blockService *blockAPIService) retriveBlock(ctx context.Context, blockRequest *types.BlockRequest) (*btcjson.GetBlockVerboseResult, *btcjson.GetBlockVerboseResult, *types.Error) {

	var block, prevBlock *btcjson.GetBlockVerboseResult
	var err error

	if blockRequest.BlockIdentifier.Index != nil {
		block, err = blockService.client.GetBlock(ctx, *blockRequest.BlockIdentifier.Index)

		if *blockRequest.BlockIdentifier.Index == 0 {
			prevBlock = &btcjson.GetBlockVerboseResult{
				Hash: "0x0",
			}
		} else {
			prevBlock, err = blockService.client.GetBlock(ctx, *blockRequest.BlockIdentifier.Index-1)
		}
	} else if blockRequest.BlockIdentifier.Hash != nil {
		block, err = blockService.client.GetBlockByHash(ctx, *blockRequest.BlockIdentifier.Hash)

		if err != nil {
			return nil, nil, ErrUnableToGetBlk
		}
		prevBlock, err = blockService.client.GetBlockByHash(ctx, block.PreviousHash)
	} else {
		block, err := blockService.client.GetLatestBlock(ctx)
		if err != nil {
			return nil, nil, ErrUnableToGetBlk
		}
		prevBlock, err = blockService.client.GetBlockByHash(ctx, block.Header.PrevBlock.String())
	}

	if err != nil {
		return nil, nil, ErrUnableToGetBlk
	}

	return block, prevBlock, nil
}

// Block retrieves the block for a given candidate
func (blockService *blockAPIService) Block(ctx context.Context, blockRequest *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	blockResponse, dbErr := blockService.blockRepository.GetBlock(*blockRequest.BlockIdentifier.Hash)

	if blockResponse == nil || dbErr != nil {
		block, prevBlock, err := blockService.retriveBlock(ctx, blockRequest)

		if err != nil {
			return nil, err
		}

		blockResponse = &types.BlockResponse{
			Block: &types.Block{
				BlockIdentifier: &types.BlockIdentifier{
					Hash: block.Hash,
				},
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash: prevBlock.Hash,
				},
				Timestamp: block.Time,
			},
			OtherTransactions: mapTransactions(block.Tx),
		}

		blockService.blockRepository.StoreBlock(block.Hash, blockResponse)

		return blockResponse, nil
	}

	log.Printf("Queried block from DB (hash: %s)\n", *blockRequest.BlockIdentifier.Hash)
	return blockResponse, nil
}

// BlockTransaction retrieves the block with the given transactions included
func (blockService *blockAPIService) BlockTransaction(ctx context.Context, blockTransaction *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {
	block, err := blockService.client.GetBlockByHashWithTransaction(ctx, blockTransaction.BlockIdentifier.Hash)

	if err != nil {
		return nil, ErrUnableToGetBlk
	}

	var networkIndex *int64 = new(int64)
	*networkIndex = 0

	for index, tx := range block.Tx {
		if tx.Hash == blockTransaction.TransactionIdentifier.Hash {
			txOperations := make([]*types.Operation, 0)

			for _, vOut := range block.Tx[index].Vout {

				if !client.IsValidPaymentType(vOut.ScriptPubKey.Type) {
					continue
				}

				for _, address := range vOut.ScriptPubKey.Addresses {
					txOperations = append(txOperations, &types.Operation{
						OperationIdentifier: &types.OperationIdentifier{
							Index:        int64(vOut.N),
							NetworkIndex: networkIndex,
						},
						Type:   client.Transfer,
						Status: client.StatusSuccess,
						Account: &types.AccountIdentifier{
							Address: address,
						},
						Amount: &types.Amount{
							Value: fmt.Sprintf("%d", int64(vOut.Value*client.BASE_CURRENCY_DECIMAL_DIVIDER)),
							Currency: &types.Currency{
								Decimals: client.BASE_CURRENCY_DECIMAL_COUNT,
								Symbol:   client.CURRENCY_SYMBOL,
							},
						},
					})
				}

			}

			return &types.BlockTransactionResponse{
				Transaction: &types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: tx.Hash,
					},
					Metadata: map[string]interface{}{
						"size":     tx.Size,
						"lockTime": tx.LockTime,
					},
					Operations: txOperations,
				},
			}, nil
		}
	}

	return nil, ErrUnableToGetTxns
}
