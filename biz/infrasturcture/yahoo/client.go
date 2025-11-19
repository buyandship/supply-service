package yahoo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"

	"net/http"

	"github.com/buyandship/bns-golib/config"
	bnsHttp "github.com/buyandship/bns-golib/http"
	"github.com/buyandship/bns-golib/retry"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	once    sync.Once
	Handler *Client
)

// Client represents the Yahoo Auction Bridge API client
type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *bnsHttp.Client
}

// NewClient creates a new Yahoo Auction Bridge client
func GetClient() *Client {
	once.Do(func() {
		client := bnsHttp.NewClient(
			bnsHttp.WithTimeout(10 * time.Second), // TODO: change to actual timeout
		)
		var baseURL string
		switch config.GlobalAppConfig.Env {
		case "dev":
			baseURL = "http://staging.yahoo-bridge.internal" // TODO: change to actual url
			// baseURL = "https://internal-stagin20251027043053843000000001-645109195.ap-northeast-1.elb.amazonaws.com"
		case "prod":
			baseURL = "https://mock-api.yahoo-auction.jp" // TODO: change to actual url
		}
		apiKey := config.GlobalAppConfig.GetString("yahoo.api_key")
		secretKey := config.GlobalAppConfig.GetString("yahoo.secret_key")
		Handler = &Client{
			baseURL:    baseURL,
			apiKey:     apiKey,
			secretKey:  secretKey,
			httpClient: client,
		}
	})
	return Handler
}

// Authentication types
type AuthType string

const (
	AuthTypeHMAC  AuthType = "hmac"
	AuthTypeOAuth AuthType = "oauth"
	AuthTypeNone  AuthType = "none"
)

// Request/Response Models

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
		TotalResultsAvailable int            `json:"totalResultsAvailable"`
		TotalResultsReturned  int            `json:"totalResultsReturned"`
		FirstResultPosition   int            `json:"firstResultPosition"`
	} `json:"ResultSet"`
}

/*
type PlaceBidPreviewResult struct {
	Signature       string `json:"Signature"`
	BidPrice        int    `json:"BidPrice"`
	Tax             int    `json:"Tax"`
	Fee             int    `json:"Fee"`
	TotalPrice      int    `json:"TotalPrice"`
	CurrentPrice    int    `json:"CurrentPrice"`
	IsRestricted    bool   `json:"IsRestricted"`
	SignatureExpiry string `json:"SignatureExpiry"`
}
*/

type Error struct {
	Message string `json:"Message"`
	Code    int    `json:"Code"`
}

type PlaceBidPreviewResponse struct {
	ResultSet struct {
		Result                PlaceBidResult `json:"Result"`
		TotalResultsAvailable int            `json:"totalResultsAvailable"`
		TotalResultsReturned  int            `json:"totalResultsReturned"`
		FirstResultPosition   int            `json:"firstResultPosition"`
	} `json:"ResultSet,omitempty"`
	Error Error `json:"Error,omitempty"`
}

// PlaceBidPreviewRequest represents a bid preview request
type PlaceBidPreviewRequest struct {
	YahooAccountID  string `json:"yahoo_account_id"`
	YsRefID         string `json:"ys_ref_id"`
	TransactionType string `json:"transaction_type"`
	AuctionID       string `json:"auction_id"`
	Price           int    `json:"price"`
	Quantity        int    `json:"quantity,omitempty"`
	Partial         bool   `json:"partial,omitempty"`
}

// AuctionItemRequest represents a request for auction item information
type AuctionItemRequest struct {
	AuctionID string `json:"auctionID"`
	AppID     string `json:"appid,omitempty"`
}

// TransactionSearchRequest represents a transaction search request
type TransactionSearchRequest struct {
	YahooAccountID string `json:"yahoo_account_id"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	Status         string `json:"status,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
}

type Transaction struct {
	TransactionID    string  `json:"transaction_id,omitempty"`
	RequestGroupID   string  `json:"request_group_id,omitempty"`
	RetryCount       int     `json:"retry_count,omitempty"`
	YsRefID          string  `json:"ys_ref_id"`
	YahooAccountID   string  `json:"yahoo_account_id,omitempty"`
	AuctionID        string  `json:"auction_id,omitempty"`
	CurrentPrice     float64 `json:"current_price,omitempty"`
	TransactionType  string  `json:"transaction_type,omitempty"`
	Status           string  `json:"status,omitempty"`
	APIEndpoint      string  `json:"api_endpoint,omitempty"`
	HTTPStatus       int     `json:"http_status,omitempty"`
	ProcessingTimeMS int     `json:"processing_time_ms,omitempty"`
	ReqPrice         float64 `json:"req_price,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
	// RequestData      supply.YahooPlaceBidReq `json:"request_data"`
	// ResponseData     PlaceBidResponse        `json:"response_data,omitempty"`
}

type GetTransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type BidStatus struct {
	HasBid       bool `json:"HasBid,omitempty"`
	MyHighestBid int  `json:"MyHighestBid,omitempty"`
	IsWinning    bool `json:"IsWinning,omitempty"`
}

// Rating represents user rating information
type Rating struct {
	Point                   int  `json:"Point" example:"150"`
	TotalGoodRating         int  `json:"TotalGoodRating,omitempty" example:"145"`
	TotalNormalRating       int  `json:"TotalNormalRating,omitempty" example:"3"`
	TotalBadRating          int  `json:"TotalBadRating,omitempty" example:"2"`
	SellerTotalGoodRating   int  `json:"SellerTotalGoodRating,omitempty" example:"120"`
	SellerTotalNormalRating int  `json:"SellerTotalNormalRating,omitempty" example:"2"`
	SellerTotalBadRating    int  `json:"SellerTotalBadRating,omitempty" example:"1"`
	IsSuspended             bool `json:"IsSuspended" example:"false"`
	IsDeleted               bool `json:"IsDeleted" example:"false"`
}

// SellerInfo represents seller information
type SellerInfo struct {
	AucUserId            string                 `json:"AucUserId,omitempty" example:"sample_seller_123"`
	Rating               Rating                 `json:"Rating,omitempty"`
	AucUserIdItemListURL string                 `json:"AucUserIdItemListURL,omitempty" example:"https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=sample_seller_123"`
	AucUserIdRatingURL   string                 `json:"AucUserIdRatingURL,omitempty" example:"https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=sample_seller_123"`
	DisplayName          string                 `json:"DisplayName,omitempty" example:"サンプルセラー"`
	StoreName            string                 `json:"StoreName,omitempty" example:"サンプルストア"`
	IconUrl128           string                 `json:"IconUrl128,omitempty" example:"https://s.yimg.jp/images/auct/profile/icon/128/sample_seller_123.jpg"`
	IconUrl256           string                 `json:"IconUrl256,omitempty" example:"https://s.yimg.jp/images/auct/profile/icon/256/sample_seller_123.jpg"`
	IconUrl512           string                 `json:"IconUrl512,omitempty" example:"https://s.yimg.jp/images/auct/profile/icon/512/sample_seller_123.jpg"`
	IsStore              bool                   `json:"IsStore,omitempty" example:"true"`
	ShoppingSellerId     string                 `json:"ShoppingSellerId,omitempty" example:"store_12345"`
	Performance          map[string]interface{} `json:"Performance,omitempty"`
}

// ImageInfo represents image information
type ImageInfo struct {
	URL    string `json:"url,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890abc.jpg"`
	Width  int    `json:"width,omitempty" example:"600"`
	Height int    `json:"height,omitempty" example:"450"`
	Alt    string `json:"alt,omitempty" example:"Sample User Image"`
}

// Images represents collection of images
type Images struct {
	Image1 *ImageInfo `json:"Image1,omitempty"`
	Image2 *ImageInfo `json:"Image2,omitempty"`
	Image3 *ImageInfo `json:"Image3,omitempty"`
}

// Thumbnails represents thumbnail URLs
type Thumbnails struct {
	Thumbnail1 string `json:"Thumbnail1,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890abc.jpg"`
	Thumbnail2 string `json:"Thumbnail2,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890def.jpg"`
	Thumbnail3 string `json:"Thumbnail3,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890ghi.jpg"`
}

// Bidder represents bidder information
type Bidder struct {
	AucUserId            string  `json:"AucUserId,omitempty"`
	Rating               Rating  `json:"Rating,omitempty"`
	AucUserIdItemListURL string  `json:"AucUserIdItemListURL,omitempty"`
	AucUserIdRatingURL   string  `json:"AucUserIdRatingURL,omitempty"`
	DisplayName          string  `json:"DisplayName,omitempty"`
	IconUrl128           string  `json:"IconUrl128,omitempty"`
	IconUrl256           string  `json:"IconUrl256,omitempty"`
	IconUrl512           string  `json:"IconUrl512,omitempty"`
	IsStore              bool    `json:"IsStore,omitempty"`
	StoreName            *string `json:"StoreName,omitempty"`
}

// HighestBidders represents highest bidders information
type HighestBidders struct {
	TotalHighestBidders int      `json:"totalHighestBidders,omitempty"`
	Bidder              []Bidder `json:"Bidder,omitempty"`
	IsMore              bool     `json:"IsMore,omitempty"`
}

// ItemStatus represents item status
type ItemStatus struct {
	Condition string `json:"Condition,omitempty"`
	Comment   string `json:"Comment,omitempty"`
}

// ItemReturnable represents return policy
type ItemReturnable struct {
	Allowed bool   `json:"Allowed,omitempty"`
	Comment string `json:"Comment,omitempty"`
}

// Option represents auction options
type Option struct {
	StoreIcon            string `json:"StoreIcon,omitempty"`
	FeaturedIcon         string `json:"FeaturedIcon,omitempty"`
	FreeshippingIcon     string `json:"FreeshippingIcon,omitempty"`
	NewItemIcon          string `json:"NewItemIcon,omitempty"`
	EasyPaymentIcon      string `json:"EasyPaymentIcon,omitempty"`
	IsTradingNaviAuction bool   `json:"IsTradingNaviAuction,omitempty"`
}

// BankMethod represents bank payment method
type BankMethod struct {
	Name   string `json:"name,omitempty"`
	BankID string `json:"bank_id,omitempty"`
}

// BankPayment represents bank payment information
type BankPayment struct {
	TotalBankMethodAvailable int          `json:"totalBankMethodAvailable,omitempty"`
	Method                   []BankMethod `json:"Method,omitempty"`
}

// EasyPayment represents easy payment information
type EasyPayment struct {
	SafeKeepingPayment string `json:"SafeKeepingPayment,omitempty"`
	IsCreditCard       bool   `json:"IsCreditCard,omitempty"`
	AllowInstallment   bool   `json:"AllowInstallment,omitempty"`
	IsPayPay           bool   `json:"IsPayPay,omitempty"`
}

// OtherPayment represents other payment methods
type OtherPayment struct {
	TotalOtherMethodAvailable int      `json:"totalOtherMethodAvailable"`
	Method                    []string `json:"Method"`
}

// Payment represents payment information
type Payment struct {
	YBank            map[string]interface{} `json:"YBank,omitempty"`
	EasyPayment      *EasyPayment           `json:"EasyPayment,omitempty"`
	Bank             *BankPayment           `json:"Bank,omitempty"`
	CashRegistration string                 `json:"CashRegistration,omitempty"`
	PostalTransfer   string                 `json:"PostalTransfer,omitempty"`
	PostalOrder      string                 `json:"PostalOrder,omitempty"`
	CashOnDelivery   string                 `json:"CashOnDelivery,omitempty"`
	Other            *OtherPayment          `json:"Other,omitempty"`
}

// ShippingMethod represents shipping method
type ShippingMethod struct {
	Type                       string  `json:"type,omitempty"`
	Index                      int     `json:"Index,omitempty"`
	Name                       string  `json:"Name,omitempty"`
	ServiceCode                int     `json:"ServiceCode,omitempty"`
	IsOfficialDelivery         bool    `json:"IsOfficialDelivery,omitempty"`
	IsPrivacyDeliveryAvailable bool    `json:"IsPrivacyDeliveryAvailable"`
	SinglePrice                float64 `json:"SinglePrice"`
	PriceURL                   string  `json:"PriceURL,omitempty"`
	DeliveryFeeSize            string  `json:"DeliveryFeeSize,omitempty"`
}

// Shipping represents shipping information
type Shipping struct {
	TotalShippingMethodAvailable int              `json:"totalShippingMethodAvailable,omitempty"`
	LowestIndex                  int              `json:"LowestIndex,omitempty"`
	Method                       []ShippingMethod `json:"Method,omitempty"`
}

// BaggageInfo represents baggage information
type BaggageInfo struct {
	Size        string `json:"Size,omitempty"`
	SizeIndex   int    `json:"SizeIndex,omitempty"`
	Weight      string `json:"Weight,omitempty"`
	WeightIndex int    `json:"WeightIndex,omitempty"`
}

// CharityOption represents charity option
type CharityOption struct {
	Proportion int `json:"Proportion,omitempty"`
}

// ItemSpec represents item specifications
type ItemSpec struct {
	Size    string `json:"Size,omitempty"`
	Segment string `json:"Segment,omitempty"`
}

// CarRegist represents car registration information
type CarRegist struct {
	Model string `json:"Model,omitempty"`
}

// CarOptions represents car options
type CarOptions struct {
	Item []string `json:"Item,omitempty"`
}

// Car represents car auction information
type Car struct {
	TotalCosts              int        `json:"TotalCosts,omitempty"`
	TaxinTotalCosts         int        `json:"TaxinTotalCosts,omitempty"`
	TotalPrice              int        `json:"TotalPrice,omitempty"`
	TaxinTotalPrice         int        `json:"TaxinTotalPrice,omitempty"`
	TotalBidorbuyPrice      int        `json:"TotalBidorbuyPrice,omitempty"`
	TaxinTotalBidorbuyPrice int        `json:"TaxinTotalBidorbuyPrice,omitempty"`
	OverheadCosts           int        `json:"OverheadCosts,omitempty"`
	TaxinOverheadCosts      int        `json:"TaxinOverheadCosts,omitempty"`
	LegalCosts              int        `json:"LegalCosts,omitempty"`
	ContactTelNumber        string     `json:"ContactTelNumber,omitempty"`
	ContactReceptionTime    string     `json:"ContactReceptionTime,omitempty"`
	ContactUrl              string     `json:"ContactUrl,omitempty"`
	Regist                  CarRegist  `json:"Regist,omitempty"`
	Options                 CarOptions `json:"Options,omitempty"`
	TotalAmountComment      string     `json:"TotalAmountComment,omitempty"`
}

// ExternalFleaMarketInfo represents external flea market information
type ExternalFleaMarketInfo struct {
	IsWinner bool `json:"IsWinner,omitempty"`
}

// ShoppingSpec represents shopping specification
type ShoppingSpec struct {
	ID      int `json:"ID,omitempty"`
	ValueID int `json:"ValueID,omitempty"`
}

// ShoppingSpecs represents shopping specifications
type ShoppingSpecs struct {
	TotalShoppingSpecs int            `json:"totalShoppingSpecs,omitempty"`
	Spec               []ShoppingSpec `json:"Spec,omitempty"`
}

// ItemTagList represents item tag list
type ItemTagList struct {
	TotalItemTagList int      `json:"totalItemTagList,omitempty"`
	Tag              []string `json:"Tag,omitempty"`
}

// ShoppingItem represents shopping item information
type ShoppingItem struct {
	PostageSetId    int  `json:"PostageSetId,omitempty"`
	PostageId       int  `json:"PostageId,omitempty"`
	LeadTimeId      int  `json:"LeadTimeId,omitempty"`
	ItemWeight      int  `json:"ItemWeight,omitempty"`
	IsOptionEnabled bool `json:"IsOptionEnabled"`
}

// SellingInfo represents selling information
type SellingInfo struct {
	PageView                        int    `json:"PageView,omitempty"`
	WatchListNum                    int    `json:"WatchListNum,omitempty"`
	ReportedViolationNum            int    `json:"ReportedViolationNum,omitempty"`
	AnsweredQAndANum                int    `json:"AnsweredQAndANum,omitempty"`
	UnansweredQAndANum              int    `json:"UnansweredQAndANum,omitempty"`
	OfferNum                        int    `json:"OfferNum,omitempty"`
	UnansweredOfferNum              int    `json:"UnansweredOfferNum,omitempty"`
	AffiliateRatio                  int    `json:"AffiliateRatio,omitempty"`
	PageViewFromAff                 int    `json:"PageViewFromAff,omitempty"`
	WatchListNumFromAff             int    `json:"WatchListNumFromAff,omitempty"`
	BidsFromAff                     int    `json:"BidsFromAff,omitempty"`
	IsWon                           bool   `json:"IsWon,omitempty"`
	IsFirstSubmit                   bool   `json:"IsFirstSubmit,omitempty"`
	Duration                        int    `json:"Duration,omitempty"`
	FirstAutoResubmitAvailableCount int    `json:"FirstAutoResubmitAvailableCount,omitempty"`
	AutoResubmitAvailableCount      int    `json:"AutoResubmitAvailableCount"`
	FeaturedDpd                     string `json:"FeaturedDpd"`
	ResubmitPriceDownRatio          int    `json:"ResubmitPriceDownRatio"`
	IsNoResubmit                    bool   `json:"IsNoResubmit"`
	BidQuantityLimit                int    `json:"BidQuantityLimit"`
}

// WinnerRating represents winner rating information
type WinnerRating struct {
	Point       int  `json:"Point,omitempty"`
	IsSuspended bool `json:"IsSuspended,omitempty"`
	IsDeleted   bool `json:"IsDeleted,omitempty"`
	IsNotRated  bool `json:"IsNotRated,omitempty"`
}

// WinnerShoppingInfo represents winner shopping information
type WinnerShoppingInfo struct {
	OrderId string `json:"OrderId,omitempty"`
}

// Winner represents winner information
type Winner struct {
	AucUserId          string              `json:"AucUserId,omitempty"`
	Rating             WinnerRating        `json:"Rating,omitempty"`
	IsRemovable        bool                `json:"IsRemovable,omitempty"`
	RemovableLimitTime int64               `json:"RemovableLimitTime,omitempty"`
	WonQuantity        int                 `json:"WonQuantity,omitempty"`
	LastBidQuantity    int                 `json:"LastBidQuantity,omitempty"`
	WonPrice           int                 `json:"WonPrice,omitempty"`
	TaxinWonPrice      float64             `json:"TaxinWonPrice,omitempty"`
	LastBidTime        int64               `json:"LastBidTime,omitempty"`
	BuyTime            string              `json:"BuyTime,omitempty"`
	IsFnaviBundledDeal bool                `json:"IsFnaviBundledDeal,omitempty"`
	ShoppingInfo       *WinnerShoppingInfo `json:"ShoppingInfo,omitempty"`
	DisplayName        string              `json:"DisplayName,omitempty"`
	IconUrl128         string              `json:"IconUrl128,omitempty"`
	IconUrl256         string              `json:"IconUrl256,omitempty"`
	IconUrl512         string              `json:"IconUrl512,omitempty"`
	IsStore            bool                `json:"IsStore,omitempty"`
}

// WinnersInfo represents winners information
type WinnersInfo struct {
	WinnersNum int      `json:"WinnersNum,omitempty"`
	Winner     []Winner `json:"Winner,omitempty"`
}

// ReserveRating represents reserve rating information
type ReserveRating struct {
	Point       int  `json:"Point,omitempty"`
	IsSuspended bool `json:"IsSuspended,omitempty"`
	IsDeleted   bool `json:"IsDeleted,omitempty"`
}

// Reserve represents reserve information
type Reserve struct {
	AucUserId         string        `json:"AucUserId,omitempty"`
	Rating            ReserveRating `json:"Rating,omitempty"`
	LastBidQuantity   int           `json:"LastBidQuantity,omitempty"`
	LastBidPrice      int           `json:"LastBidPrice,omitempty"`
	TaxinLastBidPrice float64       `json:"TaxinLastBidPrice,omitempty"`
	LastBidTime       int64         `json:"LastBidTime,omitempty"`
	DisplayName       string        `json:"DisplayName,omitempty"`
	IconUrl128        string        `json:"IconUrl128,omitempty"`
	IconUrl256        string        `json:"IconUrl256,omitempty"`
	IconUrl512        string        `json:"IconUrl512,omitempty"`
	IsStore           bool          `json:"IsStore,omitempty"`
}

// ReservesInfo represents reserves information
type ReservesInfo struct {
	ReservesNum int       `json:"ReservesNum,omitempty"`
	Reserve     []Reserve `json:"Reserve,omitempty"`
}

// Cancel represents cancel information
type Cancel struct {
	// Add fields if needed when sample data is available
}

// CancelsInfo represents cancels information
type CancelsInfo struct {
	CancelsNum int      `json:"CancelsNum,omitempty"`
	Cancel     []Cancel `json:"Cancel,omitempty"`
}

// LastBid represents last bid information
type LastBid struct {
	Price              int     `json:"Price,omitempty"`
	TaxinPrice         float64 `json:"TaxinPrice,omitempty"`
	Quantity           int     `json:"Quantity,omitempty"`
	Partial            bool    `json:"Partial,omitempty"`
	IsFnaviBundledDeal bool    `json:"IsFnaviBundledDeal,omitempty"`
}

// NextBid represents next bid information
type NextBid struct {
	Price         int `json:"Price,omitempty"`
	LimitQuantity int `json:"LimitQuantity,omitempty"`
	UnitPrice     int `json:"UnitPrice,omitempty"`
}

// BidInfo represents bid information
type BidInfo struct {
	IsHighestBidder bool    `json:"IsHighestBidder,omitempty"`
	IsWinner        bool    `json:"IsWinner,omitempty"`
	IsDeletedWinner bool    `json:"IsDeletedWinner,omitempty"`
	IsNextWinner    bool    `json:"IsNextWinner,omitempty"`
	LastBid         LastBid `json:"LastBid,omitempty"`
	NextBid         NextBid `json:"NextBid,omitempty"`
}

// OfferInfo represents offer information
type OfferInfo struct {
	OfferCondition      int `json:"OfferCondition,omitempty"`
	SellerOfferredPrice int `json:"SellerOfferredPrice,omitempty"`
	BidderOfferredPrice int `json:"BidderOfferredPrice,omitempty"`
	RemainingOfferNum   int `json:"RemainingOfferNum,omitempty"`
}

// EasyPaymentDetail represents easy payment detail
type EasyPaymentDetail struct {
	AucUserId  string `json:"AucUserId,omitempty"`
	Status     string `json:"Status,omitempty"`
	LimitTime  int64  `json:"LimitTime,omitempty"`
	UpdateTime int64  `json:"UpdateTime,omitempty"`
}

// EasyPaymentInfo represents easy payment information
type EasyPaymentInfo struct {
	EasyPayment EasyPaymentDetail `json:"EasyPayment,omitempty"`
}

// StorePayment represents store payment information
type StorePayment struct {
	TotalStorePaymentMethodAvailable int      `json:"totalStorePaymentMethodAvailable,omitempty"`
	Method                           []string `json:"Method,omitempty"`
	UpdateTime                       int64    `json:"UpdateTime,omitempty"`
}

type Seller struct {
	ID          string  `json:"Id,omitempty"`
	Rating      float64 `json:"Rating,omitempty"`
	IsSuspended bool    `json:"IsSuspended,omitempty"`
	IsDeleted   bool    `json:"IsDeleted,omitempty"`
}

type ErrorResponse struct {
	Detail []ErrorDetail `json:"detail"`
}

type ErrorDetail struct {
	Type  string      `json:"type"`
	Loc   []string    `json:"loc"`
	Msg   string      `json:"msg"`
	Input interface{} `json:"input,omitempty"`
}

// Helper method to generate HMAC signature
func (c *Client) generateHMACSignature(timestamp, method, path, body string) string {
	message := timestamp + method + path + body
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Helper method to make authenticated requests
func (c *Client) makeRequest(ctx context.Context, method, path string, params url.Values, body interface{}, authType AuthType) (*http.Response, error) {
	hlog.CtxInfof(ctx, "makeRequest: %s %s %s %v", method, path, params.Encode(), body)
	// Build URL
	fullURL := c.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	// Prepare request body
	var bodyReader io.Reader
	var bodyStr string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyStr = string(bodyBytes)
		bodyReader = bytes.NewBufferString(bodyStr)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	if authType == AuthTypeHMAC {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		signature := c.generateHMACSignature(timestamp, method, path, bodyStr)

		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("X-Timestamp", timestamp)
		req.Header.Set("X-Signature", signature)
	}

	// Set content type for POST requests
	req.Header.Set("Content-Type", "application/json")

	// Make request with retry mechanism
	var resp *http.Response
	operation := func() (*http.Response, error) {
		var err error
		resp, err = c.httpClient.Do(ctx, req)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
			// TODO: handle retryable error
		default:
			respBody, _ := io.ReadAll(resp.Body)
			hlog.CtxInfof(ctx, "status code: [%d], response body: [%s]",
				resp.StatusCode, string(respBody))
			return resp, backoff.Permanent(fmt.Errorf("%s", string(respBody)))
		}
	}

	resp, err = backoff.Retry(ctx, operation, retry.GetDefaultRetryOpts()...)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to send request after retries: %v", err)
		return resp, fmt.Errorf("failed to send request after retries: %w", err)
	}

	return resp, nil
}

// API Methods
func (c *Client) parseResponse(resp *http.Response, v interface{}) error {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			hlog.Errorf("http close error: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	hlog.Debugf("parseResponse body: %s", string(body))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// PlaceBidPreview gets bid preview with signature
func (c *Client) PlaceBidPreview(ctx context.Context, req *PlaceBidPreviewRequest) (*PlaceBidPreviewResponse, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)
	params.Set("ys_ref_id", req.YsRefID)
	params.Set("transaction_type", req.TransactionType)
	params.Set("auction_id", req.AuctionID)
	params.Set("price", strconv.Itoa(req.Price))

	if req.Quantity > 0 {
		params.Set("quantity", strconv.Itoa(req.Quantity))
	}
	if req.Partial {
		params.Set("partial", "true")
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/placeBidPreview", params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	var placeBidPreviewResponse PlaceBidPreviewResponse
	if err := c.parseResponse(resp, &placeBidPreviewResponse); err != nil {
		return nil, err
	}

	if placeBidPreviewResponse.Error.Code != 0 {
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "You were not able to place your bid. This auction has already ended.",
		}
	}

	return &placeBidPreviewResponse, nil
}

func (c *Client) MockPlaceBidPreview(ctx context.Context, req *PlaceBidPreviewRequest) (*PlaceBidPreviewResponse, error) {
	placeBidPreviewResponse := PlaceBidPreviewResponse{
		ResultSet: struct {
			Result                PlaceBidResult `json:"Result"`
			TotalResultsAvailable int            `json:"totalResultsAvailable"`
			TotalResultsReturned  int            `json:"totalResultsReturned"`
			FirstResultPosition   int            `json:"firstResultPosition"`
		}{
			Result: PlaceBidResult{
				Signature: "abc123def456...",
			},
		},
	}
	return &placeBidPreviewResponse, nil
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
		return nil, err
	}

	placeBidResponse := PlaceBidResponse{}
	if err := c.parseResponse(resp, &placeBidResponse); err != nil {
		return nil, err
	}

	transactionId := resp.Header.Get("X-Transaction-ID")
	hlog.CtxDebugf(ctx, "transactionId: %s", transactionId)
	placeBidResponse.ResultSet.Result.TransactionId = transactionId

	return &placeBidResponse, nil
}

func (c *Client) MockPlaceBid(ctx context.Context, req *PlaceBidRequest) (*PlaceBidResponse, error) {
	placeBidResponse := PlaceBidResponse{
		ResultSet: struct {
			Result                PlaceBidResult `json:"Result"`
			TotalResultsAvailable int            `json:"totalResultsAvailable"`
			TotalResultsReturned  int            `json:"totalResultsReturned"`
			FirstResultPosition   int            `json:"firstResultPosition"`
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

// SearchTransactions searches for transactions
func (c *Client) SearchTransactions(ctx context.Context, req TransactionSearchRequest) (*http.Response, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)

	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}
	if req.Status != "" {
		params.Set("status", req.Status)
	}
	if req.Limit > 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.Offset > 0 {
		params.Set("offset", strconv.Itoa(req.Offset))
	}

	return c.makeRequest(ctx, "GET", "/api/v1/transactions", params, nil, AuthTypeHMAC)
}

func (c *Client) MockGetTransaction(ctx context.Context, req *supply.YahooGetTransactionReq, yahooAccountID string) (*Transaction, error) {
	resp := Transaction{
		TransactionID:   req.TransactionID,
		YsRefID:         "YS-REF-001",
		AuctionID:       "x12345",
		CurrentPrice:    1000,
		TransactionType: "BID",
		Status:          "completed",
		ReqPrice:        1000,
		CreatedAt:       "2025-10-22T12:00:00Z",
		UpdatedAt:       "2025-10-22T12:00:01Z",
	}
	return &resp, nil
}

// ExportTransactionsCSV exports transactions as CSV
func (c *Client) ExportTransactionsCSV(ctx context.Context, req TransactionSearchRequest) (*http.Response, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)

	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}
	if req.Status != "" {
		params.Set("status", req.Status)
	}

	return c.makeRequest(ctx, "GET", "/api/v1/transactions/export/csv", params, nil, AuthTypeHMAC)
}

// Health check
func (c *Client) HealthCheck() (*http.Response, error) {
	return c.makeRequest(context.Background(), "GET", "/health", nil, nil, AuthTypeNone)
}

// Get API info
func (c *Client) GetAPIInfo() (*http.Response, error) {
	return c.makeRequest(context.Background(), "GET", "/", nil, nil, AuthTypeNone)
}

// Helper method to parse error response
func ParseErrorResponse(resp *http.Response) (*ErrorResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return nil, fmt.Errorf("failed to parse error response: %w", err)
	}

	return &errorResp, nil
}

// Helper method to parse auction item response
func ParseAuctionItemResponse(resp *http.Response) (*AuctionItemResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var auctionResp AuctionItemResponse
	if err := json.Unmarshal(body, &auctionResp); err != nil {
		return nil, fmt.Errorf("failed to parse auction item response: %w", err)
	}

	return &auctionResp, nil
}
