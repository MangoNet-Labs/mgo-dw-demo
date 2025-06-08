package third

import (
	"context"
	"errors"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"
	"gorm.io/gorm"
	"log"
	"user/internal/svc"
	"user/model"
)

func GetLatestEpoch(ctx *svc.ServiceContext) (string, error) {
	var checkpoint model.MgoCheckpoint
	err := ctx.DB.Select("sequence_number").Order("sequence_number DESC").First(&checkpoint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "0", nil // No record, return default value
		}
		return "0", err
	}
	return checkpoint.SequenceNumber, nil
}

func GetCheckpoints(SequenceNumber string, svcCtx *svc.ServiceContext) (*response.PaginatedCheckpointsResponse, error) {

	var ctx = context.Background()
	checkpointsResponse, err := svcCtx.MgoCli.MgoGetCheckpoints(ctx, request.MgoGetCheckpointsRequest{
		Cursor:          SequenceNumber,
		Limit:           50,
		DescendingOrder: false,
	})
	if err != nil {
		return nil, err
	}
	return &checkpointsResponse, nil

}

func GetCheckpointList(ctx *svc.ServiceContext) ([]model.MgoCheckpoint, error) {
	var checkpoint []model.MgoCheckpoint
	err := ctx.DB.Select("sequence_number, status").
		Where("status = 0 ").
		Limit(100).Find(&checkpoint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record, return default value
		}
		return nil, err
	}
	return checkpoint, nil
}

func QueryTransactionBlocks(SequenceNumber string, svcCtx *svc.ServiceContext) ([]response.MgoTransactionBlockResponse, error) {

	var ctx = context.Background()
	allTxs, err := AllTransactions(ctx, svcCtx, SequenceNumber, "", nil)
	if err != nil {
		log.Println("fetch error:", err)
	}
	return allTxs, nil

}

func AllTransactions(ctx context.Context, svcCtx *svc.ServiceContext, sequenceNumber string, cursor string, allData []response.MgoTransactionBlockResponse) ([]response.MgoTransactionBlockResponse, error) {

	req := request.MgoXQueryTransactionBlocksRequest{
		MgoTransactionBlockResponseQuery: request.MgoTransactionBlockResponseQuery{
			TransactionFilter: map[string]interface{}{
				"Checkpoint": sequenceNumber,
			},
			Options: request.MgoTransactionBlockOptions{
				ShowInput:          true,
				ShowBalanceChanges: true,
			},
		},
		Limit:           50,
		DescendingOrder: false,
	}
	if cursor != "" {
		req.Cursor = cursor
	}
	resp, err := svcCtx.MgoCli.MgoXQueryTransactionBlocks(ctx, req)
	if err != nil {
		return nil, err
	}

	allData = append(allData, resp.Data...)

	// Recursively get the next page
	if resp.HasNextPage {
		return AllTransactions(ctx, svcCtx, sequenceNumber, resp.NextCursor, allData)
	}

	return allData, nil
}
