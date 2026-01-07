package yahoo

import (
	"context"
	"net/url"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// PlaceBidRequest represents a bid request
type PlaceBidRequest struct {
	YahooAccountID  string `json:"yahoo_account_id"`
	YsRefID         string `json:"ys_ref_id"`
	TransactionType string `json:"transaction_type"` // BID or BUYOUT
	AuctionID       string `json:"auction_id"`
	Price           int    `json:"price"`
	Signature       string `json:"signature"`
	Quantity        int    `json:"quantity,omitempty"`
	Partial         bool   `json:"partial,omitempty"`
	IsShoppingItem  bool   `json:"is_shopping_item,omitempty"`
}

type PlaceBidResult struct {
	AuctionID       string  `json:"AuctionID" example:"x12345"`
	Title           string  `json:"Title" example:"商品名１"`
	CurrentPrice    float64 `json:"CurrentPrice" example:"1300"`
	UnitOfBidPrice  float64 `json:"UnitOfBidPrice" example:"10"`
	IsCurrentWinner bool    `json:"IsCurrentWinner" example:"false"`
	IsBuyBid        bool    `json:"IsBuyBid" example:"false"`
	IsNewBid        bool    `json:"IsNewBid" example:"true"`
	UnderReserved   bool    `json:"UnderReserved" example:"false"`
	NextBidPrice    float64 `json:"NextBidPrice" example:"1400"`
	AuctionUrl      string  `json:"AuctionUrl" example:"https://auctions.yahooapis.jp/AuctionWebService/V2/auctionItem?auctionID=x12345678"`
	AuctionItemUrl  string  `json:"AuctionItemUrl" example:"https://page.auctions.yahoo.co.jp/jp/auction/x12345678"`

	Signature     string `json:"Signature,omitempty" example:"4mYveHoMr0fad9AS.Seqc6ys2BdqMyWTxA2VG_RJDbDyZjtIU5MX8k_xqg--"`
	TransactionId string `json:"TransactionId,omitempty" example:"1234567890"`
}

type PlaceBidResponse struct {
	ResultSet struct {
		Result                PlaceBidResult `json:"Result"`
		TotalResultsAvailable int            `json:"@totalResultsAvailable,omitempty"`
		TotalResultsReturned  int            `json:"@totalResultsReturned,omitempty"`
		FirstResultPosition   int            `json:"@firstResultPosition,omitempty"`
	} `json:"ResultSet"`
}

// PlaceBid executes a bid on Yahoo Auction
func (c *Client) PlaceBid(ctx context.Context, req *PlaceBidRequest) (*PlaceBidResponse, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)
	params.Set("ys_ref_id", req.YsRefID)
	params.Set("transaction_type", req.TransactionType)
	params.Set("auction_id", req.AuctionID)
	params.Set("price", strconv.Itoa(req.Price))
	params.Set("signature", req.Signature)

	if req.Quantity > 0 {
		params.Set("quantity", strconv.Itoa(req.Quantity))
	}
	if req.Partial {
		params.Set("partial", "true")
	}

	var path string
	if req.IsShoppingItem {
		path = "/api/v1/placeBid/shpAucItem"
	} else {
		path = "/api/v1/placeBid"
	}

	resp, err := c.makeRequest(ctx, "POST", path, params, nil, AuthTypeHMAC)
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid failed: %+v", err)
		return nil, err
	}

	placeBidResponse := PlaceBidResponse{}
	if err := c.parseResponse(resp, &placeBidResponse); err != nil {
		return nil, err
	}

	transactionId := resp.Header.Get("X-Transaction-ID")
	placeBidResponse.ResultSet.Result.TransactionId = transactionId

	return &placeBidResponse, nil
}

func (c *Client) MockPlaceBid(ctx context.Context, req *PlaceBidRequest) (*PlaceBidResponse, error) {
	placeBidResponse := PlaceBidResponse{
		ResultSet: struct {
			Result                PlaceBidResult `json:"Result"`
			TotalResultsAvailable int            `json:"@totalResultsAvailable,omitempty"`
			TotalResultsReturned  int            `json:"@totalResultsReturned,omitempty"`
			FirstResultPosition   int            `json:"@firstResultPosition,omitempty"`
		}{
			Result: PlaceBidResult{
				AuctionID:       req.AuctionID,
				Title:           "Mock Title",
				CurrentPrice:    float64(req.Price),
				UnitOfBidPrice:  float64(10),
				IsCurrentWinner: false,
				IsBuyBid:        false,
				IsNewBid:        true,
				UnderReserved:   false,
				// NextBidPrice:    req.Price + 100,
				AuctionUrl:     "https://auctions.yahooapis.jp/AuctionWebService/V2/auctionItem?auctionID=x12345678",
				AuctionItemUrl: "https://page.auctions.yahoo.co.jp/jp/auction/x12345678",
			},
			TotalResultsAvailable: 1,
			TotalResultsReturned:  1,
			FirstResultPosition:   1,
		},
	}
	return &placeBidResponse, nil
}
