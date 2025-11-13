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
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/buyandship/supply-svr/biz/model/yahoo"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
			baseURL = "http://staging.yahoo-bridge.internal.com" // TODO: change to actual url
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
	AuctionID       string `json:"AuctionID" example:"x12345"`
	Title           string `json:"Title" example:"商品名１"`
	CurrentPrice    int    `json:"CurrentPrice" example:"1300"`
	UnitOfBidPrice  string `json:"UnitOfBidPrice" example:"JPY"`
	IsCurrentWinner bool   `json:"IsCurrentWinner" example:"false"`
	IsBuyBid        bool   `json:"IsBuyBid" example:"false"`
	IsNewBid        bool   `json:"IsNewBid" example:"true"`
	UnderReserved   bool   `json:"UnderReserved" example:"false"`
	NextBidPrice    int    `json:"NextBidPrice" example:"1400"`
	AuctionUrl      string `json:"AuctionUrl" example:"https://auctions.yahooapis.jp/AuctionWebService/V2/auctionItem?auctionID=x12345678"`
	AuctionItemUrl  string `json:"AuctionItemUrl" example:"https://page.auctions.yahoo.co.jp/jp/auction/x12345678"`

	Signature string `json:"Signature,omitempty" example:"4mYveHoMr0fad9AS.Seqc6ys2BdqMyWTxA2VG_RJDbDyZjtIU5MX8k_xqg--"`
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

type PlaceBidPreviewResponse struct {
	ResultSet struct {
		Result                PlaceBidResult `json:"Result"`
		TotalResultsAvailable int            `json:"totalResultsAvailable"`
		TotalResultsReturned  int            `json:"totalResultsReturned"`
		FirstResultPosition   int            `json:"firstResultPosition"`
	} `json:"ResultSet"`
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
	TransactionID    string                  `json:"transaction_id"`
	RequestGroupID   string                  `json:"request_group_id"`
	RetryCount       int                     `json:"retry_count"`
	YsRefID          string                  `json:"ys_ref_id"`
	YahooAccountID   string                  `json:"yahoo_account_id"`
	AuctionID        string                  `json:"auction_id"`
	CurrentPrice     int64                   `json:"current_price"`
	TransactionType  string                  `json:"transaction_type"`
	Status           string                  `json:"status"`
	APIEndpoint      string                  `json:"api_endpoint"`
	HTTPStatus       int                     `json:"http_status"`
	ProcessingTimeMS int                     `json:"processing_time_ms"`
	ReqPrice         int64                   `json:"req_price"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
	RequestData      supply.YahooPlaceBidReq `json:"request_data"`
	ResponseData     PlaceBidResponse        `json:"response_data"`
}

type GetTransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type BidStatus struct {
	HasBid       bool `json:"HasBid"`
	MyHighestBid int  `json:"MyHighestBid"`
	IsWinning    bool `json:"IsWinning"`
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
	AucUserId            string                 `json:"AucUserId" example:"sample_seller_123"`
	Rating               Rating                 `json:"Rating"`
	AucUserIdItemListURL string                 `json:"AucUserIdItemListURL" example:"https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=sample_seller_123"`
	AucUserIdRatingURL   string                 `json:"AucUserIdRatingURL" example:"https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=sample_seller_123"`
	DisplayName          string                 `json:"DisplayName" example:"サンプルセラー"`
	StoreName            string                 `json:"StoreName,omitempty" example:"サンプルストア"`
	IconUrl128           string                 `json:"IconUrl128" example:"https://s.yimg.jp/images/auct/profile/icon/128/sample_seller_123.jpg"`
	IconUrl256           string                 `json:"IconUrl256" example:"https://s.yimg.jp/images/auct/profile/icon/256/sample_seller_123.jpg"`
	IconUrl512           string                 `json:"IconUrl512" example:"https://s.yimg.jp/images/auct/profile/icon/512/sample_seller_123.jpg"`
	IsStore              bool                   `json:"IsStore" example:"true"`
	ShoppingSellerId     string                 `json:"ShoppingSellerId,omitempty" example:"store_12345"`
	Performance          map[string]interface{} `json:"Performance,omitempty"`
}

// ImageInfo represents image information
type ImageInfo struct {
	URL    string `json:"url" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890abc.jpg"`
	Width  int    `json:"width" example:"600"`
	Height int    `json:"height" example:"450"`
	Alt    string `json:"alt" example:"Sample User Image"`
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
	AucUserId            string `json:"AucUserId"`
	Rating               Rating `json:"Rating"`
	AucUserIdItemListURL string `json:"AucUserIdItemListURL"`
	AucUserIdRatingURL   string `json:"AucUserIdRatingURL"`
	DisplayName          string `json:"DisplayName"`
	IconUrl128           string `json:"IconUrl128"`
	IconUrl256           string `json:"IconUrl256"`
	IconUrl512           string `json:"IconUrl512"`
	IsStore              bool   `json:"IsStore"`
}

// HighestBidders represents highest bidders information
type HighestBidders struct {
	TotalHighestBidders int      `json:"totalHighestBidders"`
	Bidder              []Bidder `json:"Bidder"`
	IsMore              bool     `json:"IsMore"`
}

// ItemStatus represents item status
type ItemStatus struct {
	Condition string `json:"Condition"`
	Comment   string `json:"Comment"`
}

// ItemReturnable represents return policy
type ItemReturnable struct {
	Allowed bool   `json:"Allowed"`
	Comment string `json:"Comment"`
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
	Name   string `json:"name"`
	BankID string `json:"bank_id"`
}

// BankPayment represents bank payment information
type BankPayment struct {
	TotalBankMethodAvailable int          `json:"totalBankMethodAvailable"`
	Method                   []BankMethod `json:"Method"`
}

// EasyPayment represents easy payment information
type EasyPayment struct {
	SafeKeepingPayment string `json:"SafeKeepingPayment"`
	IsCreditCard       bool   `json:"IsCreditCard"`
	AllowInstallment   bool   `json:"AllowInstallment"`
	IsPayPay           bool   `json:"IsPayPay"`
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
	Type                       string `json:"type"`
	Index                      int    `json:"Index"`
	Name                       string `json:"Name"`
	ServiceCode                string `json:"ServiceCode"`
	IsOfficialDelivery         bool   `json:"IsOfficialDelivery"`
	IsPrivacyDeliveryAvailable bool   `json:"IsPrivacyDeliveryAvailable"`
	SinglePrice                int    `json:"SinglePrice"`
	PriceURL                   string `json:"PriceURL,omitempty"`
	DeliveryFeeSize            string `json:"DeliveryFeeSize,omitempty"`
}

// Shipping represents shipping information
type Shipping struct {
	TotalShippingMethodAvailable int              `json:"totalShippingMethodAvailable"`
	LowestIndex                  int              `json:"LowestIndex"`
	Method                       []ShippingMethod `json:"Method"`
}

// BaggageInfo represents baggage information
type BaggageInfo struct {
	Size        string `json:"Size"`
	SizeIndex   int    `json:"SizeIndex"`
	Weight      string `json:"Weight"`
	WeightIndex int    `json:"WeightIndex"`
}

// CharityOption represents charity option
type CharityOption struct {
	Proportion int `json:"Proportion"`
}

// ItemSpec represents item specifications
type ItemSpec struct {
	Size    string `json:"Size,omitempty"`
	Segment string `json:"Segment,omitempty"`
}

// CarRegist represents car registration information
type CarRegist struct {
	Model string `json:"Model"`
}

// CarOptions represents car options
type CarOptions struct {
	Item []string `json:"Item"`
}

// Car represents car auction information
type Car struct {
	TotalCosts              int        `json:"TotalCosts"`
	TaxinTotalCosts         int        `json:"TaxinTotalCosts"`
	TotalPrice              int        `json:"TotalPrice"`
	TaxinTotalPrice         int        `json:"TaxinTotalPrice"`
	TotalBidorbuyPrice      int        `json:"TotalBidorbuyPrice"`
	TaxinTotalBidorbuyPrice int        `json:"TaxinTotalBidorbuyPrice"`
	OverheadCosts           int        `json:"OverheadCosts"`
	TaxinOverheadCosts      int        `json:"TaxinOverheadCosts"`
	LegalCosts              int        `json:"LegalCosts"`
	ContactTelNumber        string     `json:"ContactTelNumber"`
	ContactReceptionTime    string     `json:"ContactReceptionTime"`
	ContactUrl              string     `json:"ContactUrl"`
	Regist                  CarRegist  `json:"Regist"`
	Options                 CarOptions `json:"Options"`
	TotalAmountComment      string     `json:"TotalAmountComment"`
}

// ExternalFleaMarketInfo represents external flea market information
type ExternalFleaMarketInfo struct {
	IsWinner bool `json:"IsWinner"`
}

// ShoppingSpec represents shopping specification
type ShoppingSpec struct {
	ID      int `json:"ID"`
	ValueID int `json:"ValueID"`
}

// ShoppingSpecs represents shopping specifications
type ShoppingSpecs struct {
	TotalShoppingSpecs int            `json:"totalShoppingSpecs"`
	Spec               []ShoppingSpec `json:"Spec"`
}

// ItemTagList represents item tag list
type ItemTagList struct {
	TotalItemTagList int      `json:"totalItemTagList"`
	Tag              []string `json:"Tag"`
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
	PageView                        int    `json:"PageView"`
	WatchListNum                    int    `json:"WatchListNum"`
	ReportedViolationNum            int    `json:"ReportedViolationNum"`
	AnsweredQAndANum                int    `json:"AnsweredQAndANum"`
	UnansweredQAndANum              int    `json:"UnansweredQAndANum"`
	OfferNum                        int    `json:"OfferNum"`
	UnansweredOfferNum              int    `json:"UnansweredOfferNum"`
	AffiliateRatio                  int    `json:"AffiliateRatio"`
	PageViewFromAff                 int    `json:"PageViewFromAff"`
	WatchListNumFromAff             int    `json:"WatchListNumFromAff"`
	BidsFromAff                     int    `json:"BidsFromAff"`
	IsWon                           bool   `json:"IsWon"`
	IsFirstSubmit                   bool   `json:"IsFirstSubmit"`
	Duration                        int    `json:"Duration"`
	FirstAutoResubmitAvailableCount int    `json:"FirstAutoResubmitAvailableCount"`
	AutoResubmitAvailableCount      int    `json:"AutoResubmitAvailableCount"`
	FeaturedDpd                     string `json:"FeaturedDpd"`
	ResubmitPriceDownRatio          int    `json:"ResubmitPriceDownRatio"`
	IsNoResubmit                    bool   `json:"IsNoResubmit"`
	BidQuantityLimit                int    `json:"BidQuantityLimit"`
}

// WinnerRating represents winner rating information
type WinnerRating struct {
	Point       int  `json:"Point"`
	IsSuspended bool `json:"IsSuspended"`
	IsDeleted   bool `json:"IsDeleted"`
	IsNotRated  bool `json:"IsNotRated,omitempty"`
}

// WinnerShoppingInfo represents winner shopping information
type WinnerShoppingInfo struct {
	OrderId string `json:"OrderId"`
}

// Winner represents winner information
type Winner struct {
	AucUserId          string              `json:"AucUserId"`
	Rating             WinnerRating        `json:"Rating"`
	IsRemovable        bool                `json:"IsRemovable"`
	RemovableLimitTime int64               `json:"RemovableLimitTime"`
	WonQuantity        int                 `json:"WonQuantity"`
	LastBidQuantity    int                 `json:"LastBidQuantity"`
	WonPrice           int                 `json:"WonPrice"`
	TaxinWonPrice      float64             `json:"TaxinWonPrice"`
	LastBidTime        int64               `json:"LastBidTime"`
	BuyTime            string              `json:"BuyTime"`
	IsFnaviBundledDeal bool                `json:"IsFnaviBundledDeal"`
	ShoppingInfo       *WinnerShoppingInfo `json:"ShoppingInfo,omitempty"`
	DisplayName        string              `json:"DisplayName"`
	IconUrl128         string              `json:"IconUrl128"`
	IconUrl256         string              `json:"IconUrl256"`
	IconUrl512         string              `json:"IconUrl512"`
	IsStore            bool                `json:"IsStore"`
}

// WinnersInfo represents winners information
type WinnersInfo struct {
	WinnersNum int      `json:"WinnersNum"`
	Winner     []Winner `json:"Winner"`
}

// ReserveRating represents reserve rating information
type ReserveRating struct {
	Point       int  `json:"Point"`
	IsSuspended bool `json:"IsSuspended"`
	IsDeleted   bool `json:"IsDeleted"`
}

// Reserve represents reserve information
type Reserve struct {
	AucUserId         string        `json:"AucUserId"`
	Rating            ReserveRating `json:"Rating"`
	LastBidQuantity   int           `json:"LastBidQuantity"`
	LastBidPrice      int           `json:"LastBidPrice"`
	TaxinLastBidPrice float64       `json:"TaxinLastBidPrice"`
	LastBidTime       int64         `json:"LastBidTime"`
	DisplayName       string        `json:"DisplayName"`
	IconUrl128        string        `json:"IconUrl128"`
	IconUrl256        string        `json:"IconUrl256"`
	IconUrl512        string        `json:"IconUrl512"`
	IsStore           bool          `json:"IsStore"`
}

// ReservesInfo represents reserves information
type ReservesInfo struct {
	ReservesNum int       `json:"ReservesNum"`
	Reserve     []Reserve `json:"Reserve"`
}

// Cancel represents cancel information
type Cancel struct {
	// Add fields if needed when sample data is available
}

// CancelsInfo represents cancels information
type CancelsInfo struct {
	CancelsNum int      `json:"CancelsNum"`
	Cancel     []Cancel `json:"Cancel"`
}

// LastBid represents last bid information
type LastBid struct {
	Price              int     `json:"Price"`
	TaxinPrice         float64 `json:"TaxinPrice"`
	Quantity           int     `json:"Quantity"`
	Partial            bool    `json:"Partial"`
	IsFnaviBundledDeal bool    `json:"IsFnaviBundledDeal"`
}

// NextBid represents next bid information
type NextBid struct {
	Price         int `json:"Price"`
	LimitQuantity int `json:"LimitQuantity"`
	UnitPrice     int `json:"UnitPrice"`
}

// BidInfo represents bid information
type BidInfo struct {
	IsHighestBidder bool    `json:"IsHighestBidder"`
	IsWinner        bool    `json:"IsWinner"`
	IsDeletedWinner bool    `json:"IsDeletedWinner"`
	IsNextWinner    bool    `json:"IsNextWinner"`
	LastBid         LastBid `json:"LastBid"`
	NextBid         NextBid `json:"NextBid"`
}

// OfferInfo represents offer information
type OfferInfo struct {
	OfferCondition      int `json:"OfferCondition"`
	SellerOfferredPrice int `json:"SellerOfferredPrice"`
	BidderOfferredPrice int `json:"BidderOfferredPrice"`
	RemainingOfferNum   int `json:"RemainingOfferNum"`
}

// EasyPaymentDetail represents easy payment detail
type EasyPaymentDetail struct {
	AucUserId  string `json:"AucUserId"`
	Status     string `json:"Status"`
	LimitTime  int64  `json:"LimitTime"`
	UpdateTime int64  `json:"UpdateTime"`
}

// EasyPaymentInfo represents easy payment information
type EasyPaymentInfo struct {
	EasyPayment EasyPaymentDetail `json:"EasyPayment"`
}

// StorePayment represents store payment information
type StorePayment struct {
	TotalStorePaymentMethodAvailable int      `json:"totalStorePaymentMethodAvailable"`
	Method                           []string `json:"Method"`
	UpdateTime                       int64    `json:"UpdateTime"`
}

// AuctionItemDetail represents detailed Yahoo Auction item information
type AuctionItemDetail struct {
	AuctionID                  string                  `json:"AuctionID" example:"x123456789"`
	CategoryID                 string                  `json:"CategoryID" example:"22216"`
	CategoryFarm               int                     `json:"CategoryFarm" example:"2"`
	CategoryIdPath             string                  `json:"CategoryIdPath" example:"0,2084005403,22216"`
	CategoryPath               string                  `json:"CategoryPath" example:"オークション > 音楽 > CD > R&B、ソウル"`
	Title                      string                  `json:"Title" example:"【新品未開封】サンプルCD アルバム 限定版"`
	SeoKeywords                string                  `json:"SeoKeywords,omitempty" example:"CD,R&B,ソウル,新品,未開封"`
	Seller                     SellerInfo              `json:"Seller"`
	ShoppingItemCode           string                  `json:"ShoppingItemCode,omitempty" example:"shopping_item_abc123"`
	AuctionItemUrl             string                  `json:"AuctionItemUrl" example:"https://page.auctions.yahoo.co.jp/jp/auction/x123456789"`
	Img                        Images                  `json:"Img"`
	ImgColor                   string                  `json:"ImgColor,omitempty" example:"red"`
	Thumbnails                 Thumbnails              `json:"Thumbnails,omitempty"`
	Initprice                  int                     `json:"Initprice" example:"1000"`
	LastInitprice              int                     `json:"LastInitprice,omitempty" example:"1200"`
	Price                      int                     `json:"Price" example:"2480"`
	TaxinStartPrice            float64                 `json:"TaxinStartPrice,omitempty" example:"1080"`
	TaxinPrice                 float64                 `json:"TaxinPrice,omitempty" example:"2678.4"`
	TaxinBidorbuy              float64                 `json:"TaxinBidorbuy,omitempty" example:"5400"`
	Bidorbuy                   int                     `json:"Bidorbuy" example:"5000"`
	TaxRate                    int                     `json:"TaxRate,omitempty" example:"8"`
	Quantity                   int                     `json:"Quantity" example:"2"`
	AvailableQuantity          int                     `json:"AvailableQuantity,omitempty" example:"1"`
	WatchListNum               int                     `json:"WatchListNum,omitempty" example:"42"`
	Bids                       int                     `json:"Bids" example:"5"`
	HighestBidders             *HighestBidders         `json:"HighestBidders,omitempty"`
	YPoint                     int                     `json:"YPoint,omitempty" example:"10"`
	ItemStatus                 ItemStatus              `json:"ItemStatus"`
	ItemReturnable             ItemReturnable          `json:"ItemReturnable,omitempty"`
	StartTime                  string                  `json:"StartTime" example:"2025-01-15T10:00:00+09:00"`
	EndTime                    string                  `json:"EndTime" example:"2025-02-15T23:59:59+09:00"`
	IsBidCreditRestrictions    bool                    `json:"IsBidCreditRestrictions,omitempty" example:"true"`
	IsBidderRestrictions       bool                    `json:"IsBidderRestrictions,omitempty" example:"true"`
	IsBidderRatioRestrictions  bool                    `json:"isBidderRatioRestrictions,omitempty" example:"false"`
	IsEarlyClosing             bool                    `json:"IsEarlyClosing,omitempty" example:"false"`
	IsAutomaticExtension       bool                    `json:"IsAutomaticExtension,omitempty" example:"true"`
	IsOffer                    bool                    `json:"IsOffer,omitempty" example:"true"`
	IsCharity                  bool                    `json:"IsCharity,omitempty" example:"false"`
	Option                     Option                  `json:"Option,omitempty"`
	Description                string                  `json:"Description" example:"<![CDATA[新品未開封のCDアルバムです。限定版となります。<br>送料無料でお届けします。]]>"`
	ItemDescriptionURL         string                  `json:"ItemDescriptionURL,omitempty" example:"https://pageX.auctions.yahoo.co.jp/jp/show/description?aID=x123456789&plainview=1"`
	DescriptionInputType       string                  `json:"Description_input_type,omitempty" example:"html"`
	Payment                    Payment                 `json:"Payment,omitempty"`
	BlindBusiness              string                  `json:"BlindBusiness,omitempty" example:"impossible"`
	SevenElevenReceive         string                  `json:"SevenElevenReceive,omitempty" example:"impossible"`
	ChargeForShipping          string                  `json:"ChargeForShipping,omitempty" example:"seller"`
	Location                   string                  `json:"Location,omitempty" example:"東京都"`
	IsWorldwide                bool                    `json:"IsWorldwide,omitempty" example:"true"`
	ShipTime                   string                  `json:"ShipTime,omitempty" example:"after"`
	ShippingInput              string                  `json:"ShippingInput,omitempty" example:"now"`
	IsYahunekoPack             bool                    `json:"IsYahunekoPack,omitempty" example:"true"`
	IsJPOfficialDelivery       bool                    `json:"IsJPOfficialDelivery,omitempty" example:"true"`
	IsPrivacyDeliveryAvailable bool                    `json:"IsPrivacyDeliveryAvailable,omitempty" example:"true"`
	ShipSchedule               int                     `json:"ShipSchedule,omitempty" example:"1"`
	ManualStartTime            string                  `json:"ManualStartTime,omitempty" example:"2025-01-15T10:00:00+09:00"`
	Shipping                   *Shipping               `json:"Shipping,omitempty"`
	BaggageInfo                BaggageInfo             `json:"BaggageInfo,omitempty"`
	IsAdult                    bool                    `json:"IsAdult,omitempty" example:"false"`
	IsCreature                 bool                    `json:"IsCreature,omitempty" example:"false"`
	IsSpecificCategory         bool                    `json:"IsSpecificCategory,omitempty" example:"false"`
	IsCharityCategory          bool                    `json:"IsCharityCategory,omitempty" example:"false"`
	CharityOption              *CharityOption          `json:"CharityOption,omitempty"`
	AnsweredQAndANum           int                     `json:"AnsweredQAndANum,omitempty" example:"3"`
	Status                     string                  `json:"Status" example:"open"`
	CpaRate                    int                     `json:"CpaRate,omitempty" example:"5"`
	BiddingViaCpa              bool                    `json:"BiddingViaCpa,omitempty" example:"true"`
	BrandLineIDPath            string                  `json:"BrandLineIDPath,omitempty" example:"brand123|line456"`
	BrandLineNamePath          string                  `json:"BrandLineNamePath,omitempty" example:"サンプルブランド|サンプルライン"`
	ItemSpec                   ItemSpec                `json:"ItemSpec,omitempty"`
	CatalogId                  string                  `json:"CatalogId,omitempty" example:"catalog_12345"`
	ProductName                string                  `json:"ProductName,omitempty" example:"サンプルCD アルバム"`
	Car                        *Car                    `json:"Car,omitempty"`
	OfferNum                   int                     `json:"OfferNum,omitempty" example:"2"`
	HasOfferAccept             bool                    `json:"HasOfferAccept,omitempty" example:"false"`
	ArticleNumber              string                  `json:"ArticleNumber,omitempty" example:"1234567890123"`
	IsDsk                      bool                    `json:"IsDsk,omitempty" example:"true"`
	CategoryInsuranceType      int                     `json:"CategoryInsuranceType,omitempty" example:"1"`
	ExternalFleaMarketInfo     *ExternalFleaMarketInfo `json:"ExternalFleaMarketInfo,omitempty"`
	ShoppingSpecs              *ShoppingSpecs          `json:"ShoppingSpecs,omitempty"`
	ItemTagList                *ItemTagList            `json:"ItemTagList,omitempty"`
	ShoppingItem               ShoppingItem            `json:"ShoppingItem,omitempty"`
	IsWatched                  bool                    `json:"IsWatched,omitempty" example:"true"`
	NotifyID                   string                  `json:"NotifyID,omitempty" example:"notify_abc123"`
	StoreSearchKeywords        string                  `json:"StoreSearchKeywords,omitempty" example:"CD,音楽,限定版"`
	SellingInfo                *SellingInfo            `json:"SellingInfo,omitempty"`
	AucUserIdContactUrl        string                  `json:"AucUserIdContactUrl,omitempty" example:"https://auctions.yahoo.co.jp/jp/show/contact?aID=x123456789"`
	WinnersInfo                *WinnersInfo            `json:"WinnersInfo,omitempty"`
	ReservesInfo               *ReservesInfo           `json:"ReservesInfo,omitempty"`
	CancelsInfo                *CancelsInfo            `json:"CancelsInfo,omitempty"`
	BidInfo                    *BidInfo                `json:"BidInfo,omitempty"`
	OfferInfo                  *OfferInfo              `json:"OfferInfo,omitempty"`
	EasyPaymentInfo            *EasyPaymentInfo        `json:"EasyPaymentInfo,omitempty"`
	StorePayment               *StorePayment           `json:"StorePayment,omitempty"`
}

// Response models
type AuctionItemResponse struct {
	ResultSet struct {
		TotalResultsAvailable int               `json:"totalResultsAvailable"`
		TotalResultsReturned  int               `json:"totalResultsReturned"`
		FirstResultPosition   int               `json:"firstResultPosition"`
		Result                AuctionItemDetail `json:"Result"`
	} `json:"ResultSet"`
}

type Seller struct {
	ID          string  `json:"Id"`
	Rating      float64 `json:"Rating"`
	IsSuspended bool    `json:"IsSuspended"`
	IsDeleted   bool    `json:"IsDeleted"`
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
	hlog.CtxDebugf(ctx, "makeRequest: %s %s %s %v", method, path, params.Encode(), body)
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
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request with retry mechanism
	var resp *http.Response
	operation := func() (*http.Response, error) {
		var err error
		resp, err = c.httpClient.Do(ctx, req)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				hlog.CtxErrorf(ctx, "http close error: %s", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			hlog.CtxWarnf(ctx, "http error, error_code: [%d], error_msg: [%s]",
				resp.StatusCode, string(respBody))
		}

		// TODO: retrable error

		if resp.StatusCode >= 500 {
			// non-retryable error
			return nil, backoff.Permanent(fmt.Errorf("server error: %d", resp.StatusCode))
		}

		return resp, nil
	}

	resp, err = backoff.Retry(ctx, operation, retry.GetDefaultRetryOpts()...)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to send request after retries: %v", err)
		return nil, fmt.Errorf("failed to send request after retries: %w", err)
	}

	return resp, nil
}

// API Methods
func (c *Client) parseResponse(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	hlog.Debugf("parseResponse: %s", string(body))

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

	hlog.CtxDebugf(ctx, "placeBidResponse: %+v", placeBidResponse)

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
				CurrentPrice:    req.Price,
				UnitOfBidPrice:  "JPY",
				IsCurrentWinner: false,
				IsBuyBid:        false,
				IsNewBid:        true,
				UnderReserved:   false,
				NextBidPrice:    req.Price + 100,
				AuctionUrl:      "https://auctions.yahooapis.jp/AuctionWebService/V2/auctionItem?auctionID=x12345678",
				AuctionItemUrl:  "https://page.auctions.yahoo.co.jp/jp/auction/x12345678",
			},
			TotalResultsAvailable: 1,
			TotalResultsReturned:  1,
			FirstResultPosition:   1,
		},
	}
	return &placeBidResponse, nil
}

// GetAuctionItem gets auction item information (public API)
func (c *Client) GetAuctionItem(ctx context.Context, req AuctionItemRequest) (*AuctionItemResponse, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/auctionItem", params, nil, AuthTypeNone)
	if err != nil {
		return nil, err
	}

	var auctionItemResponse AuctionItemResponse
	if err := c.parseResponse(resp, &auctionItemResponse); err != nil {
		return nil, err
	}

	return &auctionItemResponse, nil
}

func (c *Client) MockGetAuctionItem(ctx context.Context, req AuctionItemRequest) (*AuctionItemResponse, error) {
	auctionItemResponse := AuctionItemResponse{
		ResultSet: struct {
			TotalResultsAvailable int               `json:"totalResultsAvailable"`
			TotalResultsReturned  int               `json:"totalResultsReturned"`
			FirstResultPosition   int               `json:"firstResultPosition"`
			Result                AuctionItemDetail `json:"Result"`
		}{
			Result: AuctionItemDetail{
				AuctionID:      "x123456789",
				CategoryID:     "22216",
				CategoryFarm:   2,
				CategoryIdPath: "0,2084005403,22216",
				CategoryPath:   "オークション > 音楽 > CD > R&B、ソウル",
				Title:          "【新品未開封】サンプルCD アルバム 限定版",
				SeoKeywords:    "CD,R&B,ソウル,新品,未開封",
				Seller: SellerInfo{
					AucUserId: "sample_seller_123",
					Rating: Rating{
						Point:                   150,
						TotalGoodRating:         145,
						TotalNormalRating:       3,
						TotalBadRating:          2,
						SellerTotalGoodRating:   120,
						SellerTotalNormalRating: 2,
						SellerTotalBadRating:    1,
						IsSuspended:             false,
						IsDeleted:               false,
					},
					AucUserIdItemListURL: "https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=sample_seller_123",
					AucUserIdRatingURL:   "https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=sample_seller_123",
					DisplayName:          "サンプルセラー",
					StoreName:            "サンプルストア",
					IconUrl128:           "https://s.yimg.jp/images/auct/profile/icon/128/sample_seller_123.jpg",
					IconUrl256:           "https://s.yimg.jp/images/auct/profile/icon/256/sample_seller_123.jpg",
					IconUrl512:           "https://s.yimg.jp/images/auct/profile/icon/512/sample_seller_123.jpg",
					IsStore:              true,
					ShoppingSellerId:     "store_12345",
					Performance:          map[string]interface{}{},
				},
				ShoppingItemCode: "shopping_item_abc123",
				AuctionItemUrl:   "https://page.auctions.yahoo.co.jp/jp/auction/x123456789",
				Img: Images{
					Image1: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890abc.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像1",
					},
					Image2: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890def.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像2",
					},
					Image3: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890ghi.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像3",
					},
				},
				ImgColor: "red",
				Thumbnails: Thumbnails{
					Thumbnail1: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890abc.jpg",
					Thumbnail2: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890def.jpg",
					Thumbnail3: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890ghi.jpg",
				},
				Initprice:         1000,
				LastInitprice:     1200,
				Price:             2480,
				TaxinStartPrice:   1080,
				TaxinPrice:        2678.4,
				TaxinBidorbuy:     5400,
				Bidorbuy:          5000,
				TaxRate:           8,
				Quantity:          2,
				AvailableQuantity: 1,
				WatchListNum:      42,
				Bids:              5,
				HighestBidders: &HighestBidders{
					TotalHighestBidders: 2,
					Bidder: []Bidder{
						{
							AucUserId: "bidder_user_1",
							Rating: Rating{
								Point:       50,
								IsSuspended: false,
								IsDeleted:   false,
							},
							AucUserIdItemListURL: "https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=bidder_user_1",
							AucUserIdRatingURL:   "https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=bidder_user_1",
							DisplayName:          "入札者1",
							IconUrl128:           "https://s.yimg.jp/images/auct/profile/icon/128/bidder_user_1.jpg",
							IconUrl256:           "https://s.yimg.jp/images/auct/profile/icon/256/bidder_user_1.jpg",
							IconUrl512:           "https://s.yimg.jp/images/auct/profile/icon/512/bidder_user_1.jpg",
							IsStore:              false,
						},
						{
							AucUserId: "bidder_user_2",
							Rating: Rating{
								Point:       75,
								IsSuspended: false,
								IsDeleted:   false,
							},
							AucUserIdItemListURL: "https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=bidder_user_2",
							AucUserIdRatingURL:   "https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=bidder_user_2",
							DisplayName:          "入札者2",
							IconUrl128:           "https://s.yimg.jp/images/auct/profile/icon/128/bidder_user_2.jpg",
							IconUrl256:           "https://s.yimg.jp/images/auct/profile/icon/256/bidder_user_2.jpg",
							IconUrl512:           "https://s.yimg.jp/images/auct/profile/icon/512/bidder_user_2.jpg",
							IsStore:              false,
						},
					},
					IsMore: false,
				},
				YPoint: 10,
				ItemStatus: ItemStatus{
					Condition: "new",
					Comment:   "新品・未開封品です",
				},
				ItemReturnable: ItemReturnable{
					Allowed: true,
					Comment: "未開封のため返品可能です",
				},
				StartTime:                 "2025-01-15T10:00:00+09:00",
				EndTime:                   "2025-02-15T23:59:59+09:00",
				IsBidCreditRestrictions:   true,
				IsBidderRestrictions:      true,
				IsBidderRatioRestrictions: false,
				IsEarlyClosing:            false,
				IsAutomaticExtension:      true,
				IsOffer:                   true,
				IsCharity:                 false,
				Option: Option{
					StoreIcon:            "https://image.auctions.yahoo.co.jp/images/store.gif",
					FeaturedIcon:         "https://image.auctions.yahoo.co.jp/images/featured.gif",
					FreeshippingIcon:     "https://image.auctions.yahoo.co.jp/images/freeshipping.gif",
					NewItemIcon:          "https://image.auctions.yahoo.co.jp/images/newitem.gif",
					EasyPaymentIcon:      "https://img.yahoo.co.jp/images/pay/icon_s16.gif",
					IsTradingNaviAuction: true,
				},
				Description:          "<![CDATA[新品未開封のCDアルバムです。限定版となります。<br>送料無料でお届けします。]]>",
				ItemDescriptionURL:   "https://pageX.auctions.yahoo.co.jp/jp/show/description?aID=x123456789&plainview=1",
				DescriptionInputType: "html",
				Payment: Payment{
					YBank: map[string]interface{}{},
					EasyPayment: &EasyPayment{
						SafeKeepingPayment: "1.00",
						IsCreditCard:       true,
						AllowInstallment:   true,
						IsPayPay:           true,
					},
					Bank: &BankPayment{
						TotalBankMethodAvailable: 3,
						Method: []BankMethod{
							{Name: "三菱UFJ銀行", BankID: "0005"},
							{Name: "みずほ銀行", BankID: "0001"},
							{Name: "ゆうちょ銀行", BankID: "9900"},
						},
					},
					CashRegistration: "可能",
					PostalTransfer:   "可能",
					PostalOrder:      "可能",
					CashOnDelivery:   "可能",
					Other: &OtherPayment{
						TotalOtherMethodAvailable: 1,
						Method:                    []string{"手渡し"},
					},
				},
				BlindBusiness:              "impossible",
				SevenElevenReceive:         "impossible",
				ChargeForShipping:          "seller",
				Location:                   "東京都",
				IsWorldwide:                true,
				ShipTime:                   "after",
				ShippingInput:              "now",
				IsYahunekoPack:             true,
				IsJPOfficialDelivery:       true,
				IsPrivacyDeliveryAvailable: true,
				ShipSchedule:               1,
				ManualStartTime:            "2025-01-15T10:00:00+09:00",
				Shipping: &Shipping{
					TotalShippingMethodAvailable: 4,
					LowestIndex:                  0,
					Method: []ShippingMethod{
						{
							Type:                       "ship_name1",
							Index:                      0,
							Name:                       "ヤフネコ!（ネコポス）",
							ServiceCode:                "112",
							IsOfficialDelivery:         true,
							IsPrivacyDeliveryAvailable: true,
							SinglePrice:                210,
						},
						{
							Type:                       "ship_name2",
							Index:                      1,
							Name:                       "ヤフネコ!（宅急便コンパクト）",
							ServiceCode:                "113",
							IsOfficialDelivery:         true,
							IsPrivacyDeliveryAvailable: true,
							SinglePrice:                450,
							DeliveryFeeSize:            "80",
						},
						{
							Type:                       "ship_name3",
							Index:                      2,
							Name:                       "ゆうパケット（おてがる版）",
							ServiceCode:                "115",
							IsOfficialDelivery:         true,
							IsPrivacyDeliveryAvailable: true,
							SinglePrice:                230,
						},
						{
							Type:                       "ship_name4",
							Index:                      3,
							Name:                       "ゆうパック（おてがる版）",
							ServiceCode:                "116",
							IsOfficialDelivery:         true,
							IsPrivacyDeliveryAvailable: true,
							SinglePrice:                800,
							PriceURL:                   "https://yahoo.co.jp/shipping",
							DeliveryFeeSize:            "60",
						},
					},
				},
				BaggageInfo: BaggageInfo{
					Size:        "～70cm",
					SizeIndex:   1,
					Weight:      "～4kg",
					WeightIndex: 2,
				},
				IsAdult:            false,
				IsCreature:         false,
				IsSpecificCategory: false,
				IsCharityCategory:  false,
				CharityOption: &CharityOption{
					Proportion: 10,
				},
				AnsweredQAndANum:  3,
				Status:            "open",
				CpaRate:           5,
				BiddingViaCpa:     true,
				BrandLineIDPath:   "brand123|line456",
				BrandLineNamePath: "サンプルブランド|サンプルライン",
				ItemSpec: ItemSpec{
					Size:    "M",
					Segment: "メンズ",
				},
				CatalogId:   "catalog_12345",
				ProductName: "サンプルCD アルバム",
				Car: &Car{
					TotalCosts:              250000,
					TaxinTotalCosts:         270000,
					TotalPrice:              1250000,
					TaxinTotalPrice:         1350000,
					TotalBidorbuyPrice:      1500000,
					TaxinTotalBidorbuyPrice: 1620000,
					OverheadCosts:           150000,
					TaxinOverheadCosts:      162000,
					LegalCosts:              100000,
					ContactTelNumber:        "03-1234-5678",
					ContactReceptionTime:    "平日10:00-18:00",
					ContactUrl:              "https://example.com/contact",
					Regist: CarRegist{
						Model: "2020年式",
					},
					Options: CarOptions{
						Item: []string{"カーナビ", "ETC", "バックカメラ"},
					},
					TotalAmountComment: "諸費用込みの総額です",
				},
				OfferNum:              2,
				HasOfferAccept:        false,
				ArticleNumber:         "1234567890123",
				IsDsk:                 true,
				CategoryInsuranceType: 1,
				ExternalFleaMarketInfo: &ExternalFleaMarketInfo{
					IsWinner: false,
				},
				ShoppingSpecs: &ShoppingSpecs{
					TotalShoppingSpecs: 2,
					Spec: []ShoppingSpec{
						{ID: 100, ValueID: 1001},
						{ID: 200, ValueID: 2001},
					},
				},
				ItemTagList: &ItemTagList{
					TotalItemTagList: 2,
					Tag:              []string{"adidas", "Nike"},
				},
				ShoppingItem: ShoppingItem{
					PostageSetId:    12345,
					PostageId:       67890,
					LeadTimeId:      5000,
					ItemWeight:      500,
					IsOptionEnabled: true,
				},
				IsWatched:           true,
				NotifyID:            "notify_abc123",
				StoreSearchKeywords: "CD,音楽,限定版",
				SellingInfo: &SellingInfo{
					PageView:                        1250,
					WatchListNum:                    42,
					ReportedViolationNum:            0,
					AnsweredQAndANum:                3,
					UnansweredQAndANum:              1,
					OfferNum:                        2,
					UnansweredOfferNum:              0,
					AffiliateRatio:                  5,
					PageViewFromAff:                 180,
					WatchListNumFromAff:             8,
					BidsFromAff:                     2,
					IsWon:                           false,
					IsFirstSubmit:                   true,
					Duration:                        7,
					FirstAutoResubmitAvailableCount: 3,
					AutoResubmitAvailableCount:      3,
					FeaturedDpd:                     "500",
					ResubmitPriceDownRatio:          5,
					IsNoResubmit:                    false,
					BidQuantityLimit:                5,
				},
				AucUserIdContactUrl: "https://auctions.yahoo.co.jp/jp/show/contact?aID=x123456789",
				WinnersInfo: &WinnersInfo{
					WinnersNum: 1,
					Winner: []Winner{
						{
							AucUserId: "winner_user_1",
							Rating: WinnerRating{
								Point:       120,
								IsSuspended: false,
								IsDeleted:   false,
								IsNotRated:  false,
							},
							IsRemovable:        true,
							RemovableLimitTime: 1739836800,
							WonQuantity:        1,
							LastBidQuantity:    1,
							WonPrice:           2480,
							TaxinWonPrice:      2678.4,
							LastBidTime:        1707955199,
							BuyTime:            "2025-02-15T23:59:59+09:00",
							IsFnaviBundledDeal: false,
							ShoppingInfo: &WinnerShoppingInfo{
								OrderId: "order_12345",
							},
							DisplayName: "落札者1",
							IconUrl128:  "https://s.yimg.jp/images/auct/profile/icon/128/winner_user_1.jpg",
							IconUrl256:  "https://s.yimg.jp/images/auct/profile/icon/256/winner_user_1.jpg",
							IconUrl512:  "https://s.yimg.jp/images/auct/profile/icon/512/winner_user_1.jpg",
							IsStore:     false,
						},
					},
				},
				ReservesInfo: &ReservesInfo{
					ReservesNum: 1,
					Reserve: []Reserve{
						{
							AucUserId: "reserve_user_1",
							Rating: ReserveRating{
								Point:       85,
								IsSuspended: false,
								IsDeleted:   false,
							},
							LastBidQuantity:   1,
							LastBidPrice:      2450,
							TaxinLastBidPrice: 2646,
							LastBidTime:       1707955150,
							DisplayName:       "次点者1",
							IconUrl128:        "https://s.yimg.jp/images/auct/profile/icon/128/reserve_user_1.jpg",
							IconUrl256:        "https://s.yimg.jp/images/auct/profile/icon/256/reserve_user_1.jpg",
							IconUrl512:        "https://s.yimg.jp/images/auct/profile/icon/512/reserve_user_1.jpg",
							IsStore:           false,
						},
					},
				},
				CancelsInfo: &CancelsInfo{
					CancelsNum: 0,
					Cancel:     []Cancel{},
				},
				BidInfo: &BidInfo{
					IsHighestBidder: true,
					IsWinner:        false,
					IsDeletedWinner: false,
					IsNextWinner:    false,
					LastBid: LastBid{
						Price:              2480,
						TaxinPrice:         2678.4,
						Quantity:           1,
						Partial:            false,
						IsFnaviBundledDeal: false,
					},
					NextBid: NextBid{
						Price:         2580,
						LimitQuantity: 1,
						UnitPrice:     100,
					},
				},
				OfferInfo: &OfferInfo{
					OfferCondition:      1,
					SellerOfferredPrice: 4500,
					BidderOfferredPrice: 4000,
					RemainingOfferNum:   2,
				},
				EasyPaymentInfo: &EasyPaymentInfo{
					EasyPayment: EasyPaymentDetail{
						AucUserId:  "winner_user_1",
						Status:     "completed",
						LimitTime:  1710547199,
						UpdateTime: 1709337599,
					},
				},
				StorePayment: &StorePayment{
					TotalStorePaymentMethodAvailable: 2,
					Method:                           []string{"クレジットカード", "代金引換"},
					UpdateTime:                       1704614400,
				},
			},
			TotalResultsAvailable: 1,
			TotalResultsReturned:  1,
			FirstResultPosition:   1,
		},
	}
	return &auctionItemResponse, nil
}

// GetAuctionItemAuth gets authenticated auction item information
func (c *Client) GetAuctionItemAuth(ctx context.Context, req AuctionItemRequest, yahooAccountID string) (*AuctionItemResponse, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	params.Set("yahoo_account_id", yahooAccountID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/auctionItemAuth", params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	var auctionItemAuthResponse AuctionItemResponse
	if err := c.parseResponse(resp, &auctionItemAuthResponse); err != nil {
		return nil, err
	}

	return &auctionItemAuthResponse, nil
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

// GetTransaction gets specific transaction details
func (c *Client) GetTransaction(ctx context.Context, req *supply.YahooGetTransactionReq, yahooAccountID string) (*Transaction, error) {
	path := fmt.Sprintf("/api/v1/transactions/%s", req.TransactionID)
	params := url.Values{}
	params.Set("yahoo_account_id", yahooAccountID)

	resp, err := c.makeRequest(ctx, "GET", path, params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	if err := c.parseResponse(resp, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
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
		RequestData:     supply.YahooPlaceBidReq{},
		ResponseData: PlaceBidResponse{
			ResultSet: struct {
				Result                PlaceBidResult `json:"Result"`
				TotalResultsAvailable int            `json:"totalResultsAvailable"`
				TotalResultsReturned  int            `json:"totalResultsReturned"`
				FirstResultPosition   int            `json:"firstResultPosition"`
			}{
				Result: PlaceBidResult{
					AuctionID:       "x12345",
					Title:           "Mock Title",
					CurrentPrice:    1000,
					UnitOfBidPrice:  "JPY",
					IsCurrentWinner: false,
					IsBuyBid:        false,
					IsNewBid:        true,
					UnderReserved:   false,
					NextBidPrice:    1100,
					AuctionUrl:      "https://auctions.yahooapis.jp/AuctionWebService/V2/auctionItem?auctionID=x12345678",
					AuctionItemUrl:  "https://page.auctions.yahoo.co.jp/jp/auction/x12345678",
					Signature:       "abc123def456...",
				},
				TotalResultsAvailable: 1,
				TotalResultsReturned:  1,
				FirstResultPosition:   1,
			},
		},
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

type GetCategoryTreeResponse struct {
	ResultSet struct {
		Result                yahoo.Category `json:"Result"`
		TotalResultsAvailable int            `json:"totalResultsAvailable"`
		TotalResultsReturned  int            `json:"totalResultsReturned"`
		FirstResultPosition   int            `json:"firstResultPosition"`
	}
}

func (c *Client) GetCategoryTree(ctx context.Context, req *supply.YahooGetCategoryTreeReq) (*GetCategoryTreeResponse, error) {
	params := url.Values{}
	params.Set("category", req.Category)
	params.Set("adf", req.Adf)
	params.Set("is_fnavi_only", req.IsFnaviOnly)

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/categoryTree", params, nil, AuthTypeNone)
	if err != nil {
		return nil, err
	}

	var httpResp GetCategoryTreeResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}

func (c *Client) MockGetCategoryTree(ctx context.Context, req *supply.YahooGetCategoryTreeReq) (*GetCategoryTreeResponse, error) {
	httpResp := GetCategoryTreeResponse{
		ResultSet: struct {
			Result                yahoo.Category `json:"Result"`
			TotalResultsAvailable int            `json:"totalResultsAvailable"`
			TotalResultsReturned  int            `json:"totalResultsReturned"`
			FirstResultPosition   int            `json:"firstResultPosition"`
		}{
			Result: yahoo.Category{
				CategoryID:       1234567890,
				CategoryName:     "Mock Category",
				CategoryPath:     "Mock Category Path",
				CategoryIDPath:   "0,1234567890",
				ParentCategoryID: 0,
				IsLeaf:           false,
				Depth:            1,
				Order:            0,
				IsLink:           false,
				IsAdult:          false,
				ChildCategoryNum: 0,
				IsLeafToLink:     nil,
				ChildCategory: []*yahoo.Category{
					{
						CategoryID:     1234567891,
						CategoryName:   "Mock Child Category",
						CategoryPath:   "Mock Child Category Path",
						CategoryIDPath: "0,1234567891",
					},
				},
			},
			TotalResultsAvailable: 1,
			TotalResultsReturned:  1,
			FirstResultPosition:   1,
		},
	}
	return &httpResp, nil
}
