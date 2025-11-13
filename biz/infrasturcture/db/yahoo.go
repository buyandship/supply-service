package db

import (
	"context"

	"github.com/buyandship/bns-golib/trace"
	"github.com/buyandship/supply-svr/biz/model/yahoo"
	"gorm.io/gorm"
)

func (h *H) InsertBuyoutBidRequest(ctx context.Context, order *yahoo.BidRequest) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertYahooOrder")
	defer trace.EndSpan(span, err)

	db := h.cli.WithContext(ctx)

	var existing yahoo.BidRequest
	if err := db.Where("order_id = ?", order.OrderID).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		tx := db.Begin()
		// create bid request
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return err
		}
		// create transaction
		if err := tx.Create(&yahoo.YahooTransaction{
			BidRequestID: order.OrderID,
			Price:        order.MaxBid,
			Status:       "CREATED",
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
	*order = existing
	return nil
}

func (h *H) UpdateBuyoutBidRequest(ctx context.Context, order *yahoo.BidRequest) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpdateYahooOrder")
	defer trace.EndSpan(span, err)

	tx := h.cli.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&yahoo.BidRequest{}).Where("order_id = ?", order.OrderID).Updates(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&yahoo.YahooTransaction{
		BidRequestID: order.OrderID,
		Price:        order.MaxBid,
		Status:       order.Status, // TODO: transaction statu
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
