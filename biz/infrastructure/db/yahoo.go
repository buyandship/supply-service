package db

import (
	"context"

	"github.com/buyandship/bns-golib/trace"
	"github.com/buyandship/supply-service/biz/model/yahoo"
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

		return tx.Commit().Error
	}
	*order = existing
	return nil
}

func (h *H) GetBidRequest(ctx context.Context, orderID string) (order *yahoo.BidRequest, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetBidRequest")
	defer trace.EndSpan(span, err)

	db := h.cli.WithContext(ctx)

	order = &yahoo.BidRequest{}
	if err := db.Where("order_id = ?", orderID).First(order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, nil
	}
	return order, nil
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
		BidRequestID:  order.OrderID,
		Price:         order.MaxBid,
		Status:        order.Status, // TODO: transaction statu
		TransactionID: order.TransactionID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (h *H) InsertBidAuctionItem(ctx context.Context, item *yahoo.BidAuctionItem) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertBidAuctionItem")
	defer trace.EndSpan(span, err)

	db := h.cli.WithContext(ctx)

	if err := db.Create(&item).Error; err != nil {
		return err
	}

	return nil
}

func (h *H) GetShippingFee(ctx context.Context) (fees []yahoo.ShippingFee, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetShippingFee")
	defer trace.EndSpan(span, err)

	db := h.cli.WithContext(ctx)

	shippingFees := []yahoo.ShippingFee{}
	if err := db.Find(&shippingFees).Error; err != nil {
		return nil, err
	}
	return shippingFees, nil
}
