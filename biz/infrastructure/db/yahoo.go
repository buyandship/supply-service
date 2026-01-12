package db

import (
	"context"
	"time"

	"github.com/buyandship/bns-golib/trace"
	"github.com/buyandship/supply-service/biz/model/yahoo"
	YahooModel "github.com/buyandship/supply-service/biz/model/yahoo"
	"gorm.io/gorm"
)

func (h *H) InsertBidRequest(ctx context.Context, order *yahoo.BidRequest) (err error) {
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

func (h *H) GetBidRequestByAuctionID(ctx context.Context, auctionID string) (order *yahoo.BidRequest, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetBidRequestByAuctionID")
	defer trace.EndSpan(span, err)

	db := h.cli.WithContext(ctx)

	order = &yahoo.BidRequest{}
	if err := db.Where("auction_id = ? AND status IN (?)", auctionID, []string{yahoo.StatusCreated, yahoo.StatusBiddingInProgress}).First(order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return order, nil
}

// For Buyout Request
func (h *H) UpdateBuyoutRequest(ctx context.Context, order *yahoo.BidRequest) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpdateBuyoutRequest")
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

// This function is used to add bid when the current bid is less than the current max bid.
func (h *H) AddBidRequest(ctx context.Context, bid *yahoo.YahooTransaction) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "AddBidRequest")
	defer trace.EndSpan(span, err)

	tx := h.cli.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// get transaction
	var transaction yahoo.YahooTransaction
	if err := tx.Model(&yahoo.YahooTransaction{}).Where("bid_request_id = ? and status = ?", bid.BidRequestID, yahoo.StatusBiddingInProgress).First(&transaction).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}
	} else {
		// out bid old transaction
		if err := tx.Model(&yahoo.YahooTransaction{}).Where("id = ?", transaction.ID).Updates(&yahoo.YahooTransaction{
			Status: yahoo.StatusOutBid,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Create(&yahoo.YahooTransaction{
		BidRequestID:  bid.BidRequestID,
		Price:         bid.Price,
		Status:        bid.Status,
		TransactionID: bid.TransactionID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// This function is used to out bid the bid request
func (h *H) OutBidRequest(ctx context.Context, orderID string) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "OutBidRequest")
	defer trace.EndSpan(span, err)

	tx := h.cli.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// get transaction
	var transaction yahoo.YahooTransaction
	if err := tx.Model(&yahoo.YahooTransaction{}).Where("bid_request_id = ? and status = ?", orderID, yahoo.StatusBiddingInProgress).First(&transaction).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}
	} else {
		// out bid old transaction
		if err := tx.Model(&yahoo.YahooTransaction{}).Where("id = ?", transaction.ID).Updates(&yahoo.YahooTransaction{
			Status: yahoo.StatusOutBid,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&yahoo.BidRequest{}).Where("order_id = ?", orderID).Updates(&yahoo.BidRequest{
		Status: yahoo.StatusLostBid,
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

	var existing yahoo.BidAuctionItem
	if err := db.Where("bid_request_id = ?", item.BidRequestID).First(&existing).Error; err == nil {
		return nil
	}

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

func (h *H) SwitchYahooAccount(ctx context.Context, accountId int32) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "SwitchYahooAccount")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	tx := sql.Begin()

	if err := tx.Model(&YahooModel.Account{}).Where("id = ?", accountId).First(&YahooModel.Account{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&YahooModel.Account{}).Where("active_at is not null").Update("active_at", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&YahooModel.Account{}).Where("id = ?", accountId).Update("active_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
