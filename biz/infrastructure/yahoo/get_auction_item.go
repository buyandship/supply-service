package yahoo

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/buyandship/bns-golib/cache"
	globalConfig "github.com/buyandship/bns-golib/config"
	"github.com/buyandship/supply-service/biz/common/config"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/shopspring/decimal"
)

const (
	DefaultShippingFee = 1230
)

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

// Images represents collection of images
type Images struct {
	Image1  *AuctionImage `json:"Image1,omitempty"`
	Image2  *AuctionImage `json:"Image2,omitempty"`
	Image3  *AuctionImage `json:"Image3,omitempty"`
	Image4  *AuctionImage `json:"Image4,omitempty"`
	Image5  *AuctionImage `json:"Image5,omitempty"`
	Image6  *AuctionImage `json:"Image6,omitempty"`
	Image7  *AuctionImage `json:"Image7,omitempty"`
	Image8  *AuctionImage `json:"Image8,omitempty"`
	Image9  *AuctionImage `json:"Image9,omitempty"`
	Image10 *AuctionImage `json:"Image10,omitempty"`
}

func (i *Images) List() []string {
	if i.Image1 == nil {
		return []string{}
	}

	urls := []string{i.Image1.URL}
	if i.Image2 != nil {
		urls = append(urls, i.Image2.URL)
	}
	if i.Image3 != nil {
		urls = append(urls, i.Image3.URL)
	}
	if i.Image4 != nil {
		urls = append(urls, i.Image4.URL)
	}
	if i.Image5 != nil {
		urls = append(urls, i.Image5.URL)
	}
	if i.Image6 != nil {
		urls = append(urls, i.Image6.URL)
	}
	if i.Image7 != nil {
		urls = append(urls, i.Image7.URL)
	}
	if i.Image8 != nil {
		urls = append(urls, i.Image8.URL)
	}
	if i.Image9 != nil {
		urls = append(urls, i.Image9.URL)
	}
	if i.Image10 != nil {
		urls = append(urls, i.Image10.URL)
	}

	return urls
}

// Thumbnails represents thumbnail URLs
type Thumbnails struct {
	Thumbnail1  string `json:"Thumbnail1,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890abc.jpg"`
	Thumbnail2  string `json:"Thumbnail2,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890def.jpg"`
	Thumbnail3  string `json:"Thumbnail3,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890ghi.jpg"`
	Thumbnail4  string `json:"Thumbnail4,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890jkl.jpg"`
	Thumbnail5  string `json:"Thumbnail5,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890mno.jpg"`
	Thumbnail6  string `json:"Thumbnail6,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890pqr.jpg"`
	Thumbnail7  string `json:"Thumbnail7,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890stu.jpg"`
	Thumbnail8  string `json:"Thumbnail8,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890vwx.jpg"`
	Thumbnail9  string `json:"Thumbnail9,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890yz.jpg"`
	Thumbnail10 string `json:"Thumbnail10,omitempty" example:"https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890123.jpg"`
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

// AuctionItemDetail represents detailed Yahoo Auction item information
type AuctionItemDetail struct {
	AuctionID                  string                  `json:"AuctionID,omitempty" example:"x123456789"`
	CategoryID                 int                     `json:"CategoryID" example:"22216"`
	CategoryFarm               int                     `json:"CategoryFarm,omitempty" example:"2"`
	CategoryIdPath             string                  `json:"CategoryIdPath,omitempty" example:"0,2084005403,22216"`
	CategoryPath               string                  `json:"CategoryPath,omitempty" example:"オークション > 音楽 > CD > R&B、ソウル"`
	Title                      string                  `json:"Title,omitempty" example:"【新品未開封】サンプルCD アルバム 限定版"`
	SeoKeywords                string                  `json:"SeoKeywords,omitempty" example:"CD,R&B,ソウル,新品,未開封"`
	Seller                     SellerInfo              `json:"Seller,omitempty"`
	ShoppingItemCode           string                  `json:"ShoppingItemCode,omitempty" example:"shopping_item_abc123"`
	AuctionItemUrl             string                  `json:"AuctionItemUrl,omitempty" example:"https://page.auctions.yahoo.co.jp/jp/auction/x123456789"`
	Img                        Images                  `json:"Img,omitempty"`
	ImgColor                   string                  `json:"ImgColor,omitempty" example:"red"`
	Thumbnails                 Thumbnails              `json:"Thumbnails,omitempty"`
	Initprice                  float64                 `json:"Initprice,omitempty" example:"1000"`
	LastInitprice              float64                 `json:"LastInitprice,omitempty" example:"1200"`
	Price                      float64                 `json:"Price,omitempty" example:"2480"`
	TaxinStartPrice            float64                 `json:"TaxinStartPrice,omitempty" example:"1080"`
	TaxinPrice                 float64                 `json:"TaxinPrice,omitempty" example:"2678.4"`
	TaxinBidorbuy              float64                 `json:"TaxinBidorbuy,omitempty" example:"5400"`
	Bidorbuy                   float64                 `json:"Bidorbuy,omitempty" example:"5000"`
	TaxRate                    int                     `json:"TaxRate,omitempty" example:"8"`
	Quantity                   int                     `json:"Quantity,omitempty" example:"2"`
	AvailableQuantity          int                     `json:"AvailableQuantity" example:"1"`
	WatchListNum               int                     `json:"WatchListNum,omitempty" example:"42"`
	Bids                       int                     `json:"Bids,omitempty" example:"5"`
	HighestBidders             *HighestBidders         `json:"HighestBidders,omitempty"`
	YPoint                     int                     `json:"YPoint,omitempty" example:"10"`
	ItemStatus                 ItemStatus              `json:"ItemStatus,omitempty"`
	ItemReturnable             ItemReturnable          `json:"ItemReturnable,omitempty"`
	StartTime                  string                  `json:"StartTime,omitempty" example:"2025-01-15T10:00:00+09:00"`
	EndTime                    string                  `json:"EndTime,omitempty" example:"2025-02-15T23:59:59+09:00"`
	IsBidCreditRestrictions    bool                    `json:"IsBidCreditRestrictions,omitempty" example:"true"`
	IsBidderRestrictions       bool                    `json:"IsBidderRestrictions,omitempty" example:"true"`
	IsBidderRatioRestrictions  bool                    `json:"isBidderRatioRestrictions,omitempty" example:"false"`
	IsEarlyClosing             bool                    `json:"IsEarlyClosing,omitempty" example:"false"`
	IsAutomaticExtension       bool                    `json:"IsAutomaticExtension,omitempty" example:"true"`
	IsOffer                    bool                    `json:"IsOffer,omitempty" example:"true"`
	IsCharity                  bool                    `json:"IsCharity,omitempty" example:"false"`
	Option                     Option                  `json:"Option,omitempty"`
	Description                string                  `json:"Description,omitempty" example:"<![CDATA[新品未開封のCDアルバムです。限定版となります。<br>送料無料でお届けします。]]>"`
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
	ShoppingItem               *ShoppingItem           `json:"ShoppingItem,omitempty"`
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

	ShippingFee int64 `json:"ShippingFee" example:"1230"`
}

func (a *AuctionItemDetail) GetBuyoutPriceString() string {
	if a.TaxinBidorbuy > 0 {
		return decimal.NewFromFloat(a.TaxinBidorbuy).StringFixed(0)
	}
	return decimal.NewFromFloat(a.Bidorbuy).StringFixed(0)
}

func (a *AuctionItemDetail) GetBuyoutPriceWithShippingFee() string {
	buyoutPrice := a.GetBuyoutPriceString()
	shippingFee := decimal.NewFromInt(a.ShippingFee)
	buyoutPriceDecimal, err := decimal.NewFromString(buyoutPrice)
	if err != nil {
		return ""
	}
	return buyoutPriceDecimal.Add(shippingFee).StringFixed(0)
}

func (a *AuctionItemDetail) GetBidPriceString() string {
	return decimal.NewFromFloat(a.Price).StringFixed(0)
}

func (a *AuctionItemDetail) GetBidPriceWithShippingFee() string {
	bidPrice := a.GetBidPriceString()
	shippingFee := decimal.NewFromInt(a.ShippingFee)
	bidPriceDecimal, err := decimal.NewFromString(bidPrice)
	if err != nil {
		return ""
	}
	return bidPriceDecimal.Add(shippingFee).StringFixed(0)
}

func (a *AuctionItemDetail) GetDescription() string {
	if a.Description == "" {
		return ""
	}

	// Remove CDATA wrapper if present
	desc := a.Description
	if strings.HasPrefix(desc, "<![CDATA[") && strings.HasSuffix(desc, "]]>") {
		desc = strings.TrimPrefix(desc, "<![CDATA[")
		desc = strings.TrimSuffix(desc, "]]>")
	}

	// Parse HTML and extract text
	doc, err := html.Parse(strings.NewReader(desc))
	if err != nil {
		// If parsing fails, return original description
		return a.Description
	}

	var result strings.Builder
	var lastChar byte
	var extractText func(*html.Node, bool, *html.Node)
	extractText = func(n *html.Node, addSpace bool, prevSibling *html.Node) {
		switch n.Type {
		case html.TextNode:
			// Replace literal \n with actual newlines
			text := strings.ReplaceAll(n.Data, "\\n", "\n")
			text = strings.TrimSpace(text)
			if text != "" {
				if addSpace && result.Len() > 0 && lastChar != '\n' {
					result.WriteString(" ")
					lastChar = ' '
				}
				result.WriteString(text)
				if len(text) > 0 {
					lastChar = text[len(text)-1]
				}
			}
		case html.ElementNode:
			// Handle line breaks and block elements
			switch n.Data {
			case "br", "BR":
				if result.Len() > 0 {
					// Check if previous sibling was also a BR tag for paragraph break
					isConsecutiveBR := prevSibling != nil &&
						prevSibling.Type == html.ElementNode &&
						(prevSibling.Data == "br" || prevSibling.Data == "BR")

					if isConsecutiveBR {
						// Consecutive BR tags create a paragraph break (double newline)
						result.WriteString("\n")
					}
					result.WriteString("\n")
					lastChar = '\n'
				}
			case "p", "P", "div", "DIV":
				if result.Len() > 0 {
					result.WriteString("\n")
					lastChar = '\n'
				}
			}
		}
		var prev *html.Node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// Add space between text nodes, but not after block elements
			shouldAddSpace := n.Type != html.ElementNode || (n.Data != "br" && n.Data != "BR" && n.Data != "p" && n.Data != "P" && n.Data != "div" && n.Data != "DIV")
			extractText(c, shouldAddSpace, prev)
			prev = c
		}
	}
	extractText(doc, false, nil)

	plainText := result.String()

	// Normalize whitespace: replace multiple spaces/tabs with single space, but preserve newlines
	spaceRegex := regexp.MustCompile(`[ \t]+`)
	plainText = spaceRegex.ReplaceAllString(plainText, " ")

	// Normalize multiple consecutive newlines (3+) to double newline (paragraph break)
	// This preserves intentional paragraph breaks from <BR><BR> while removing excessive newlines
	newlineRegex := regexp.MustCompile(`\n{3,}`)
	plainText = newlineRegex.ReplaceAllString(plainText, "\n\n")

	// Trim leading/trailing whitespace
	plainText = strings.TrimSpace(plainText)

	return plainText
}

func (a *AuctionItemDetail) ToBidAuctionItem() *yahoo.BidAuctionItem {
	auctionID := base64.StdEncoding.EncodeToString([]byte(a.AuctionID))
	status := base64.StdEncoding.EncodeToString([]byte(a.Status))
	name := base64.StdEncoding.EncodeToString([]byte(a.Title))
	currentPrice := base64.StdEncoding.EncodeToString([]byte(strconv.FormatFloat(a.Price, 'f', -1, 64)))
	buyoutPrice := base64.StdEncoding.EncodeToString([]byte(strconv.FormatFloat(a.Bidorbuy, 'f', -1, 64)))
	// itemType := base64.StdEncoding.EncodeToString([]byte(a.ItemType))
	description := base64.StdEncoding.EncodeToString([]byte(a.Description)) // TODO: parse HTML to plain text
	startTime := base64.StdEncoding.EncodeToString([]byte(a.StartTime))
	endTime := base64.StdEncoding.EncodeToString([]byte(a.EndTime))
	itemCategoryID := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(a.CategoryID)))
	itemCondition := base64.StdEncoding.EncodeToString([]byte(a.ItemStatus.Condition))
	var itemBrand string
	if a.ItemTagList != nil {
		itemBrand = base64.StdEncoding.EncodeToString([]byte(strings.Join(a.ItemTagList.Tag, ",")))
	}
	itemWatchListNum := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(a.WatchListNum)))
	sellerID := base64.StdEncoding.EncodeToString([]byte(a.Seller.AucUserId))
	sellerRating := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(a.Seller.Rating.Point)))

	return &yahoo.BidAuctionItem{
		AuctionID:        auctionID,
		Status:           status,
		Name:             name,
		CurrentPrice:     currentPrice,
		BuyoutPrice:      buyoutPrice,
		Description:      description,
		StartTime:        startTime,
		EndTime:          endTime,
		ItemCategoryID:   itemCategoryID,
		ItemCondition:    itemCondition,
		ItemWatchListNum: itemWatchListNum,
		SellerID:         sellerID,
		SellerRating:     sellerRating,
		ItemBrand:        itemBrand,
	}
}

// AuctionItemRequest represents a request for auction item information
type AuctionItemRequest struct {
	AuctionID string `json:"auctionID"`
	AppID     string `json:"appid,omitempty"`
}

// Response models
type AuctionItemResponse struct {
	ResultSet struct {
		TotalResultsAvailable int               `json:"@totalResultsAvailable"`
		TotalResultsReturned  int               `json:"@totalResultsReturned"`
		FirstResultPosition   int               `json:"@firstResultPosition"`
		Result                AuctionItemDetail `json:"Result"`
	} `json:"ResultSet"`
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
			TotalResultsAvailable int               `json:"@totalResultsAvailable"`
			TotalResultsReturned  int               `json:"@totalResultsReturned"`
			FirstResultPosition   int               `json:"@firstResultPosition"`
			Result                AuctionItemDetail `json:"Result"`
		}{
			Result: AuctionItemDetail{
				AuctionID:      "x123456789",
				CategoryID:     22216,
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
				ImgColor:         "red",
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
				YPoint:            10,
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
				ShoppingItem: &ShoppingItem{
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
func (c *Client) GetAuctionItemAuth(ctx context.Context, req AuctionItemRequest) (*AuctionItemResponse, error) {

	var yahooAccountID string
	switch globalConfig.GlobalAppConfig.Env {
	case "dev":
		yahooAccountID = config.DevYahoo02AccountID
	case "prod":
		yahooAccountID = config.ProdYahoo02AccountID
	}

	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	params.Set("yahoo_account_id", yahooAccountID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/auctionItemAuth", params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				{
					return nil, bizErr.BizError{
						Status:  http.StatusNotFound,
						ErrCode: http.StatusNotFound,
						ErrMsg:  "auction item not found",
					}
				}
			}
		}
		hlog.CtxErrorf(ctx, "get auction item:[%s] auth error: %v", req.AuctionID, err)
		return nil, bizErr.BizError{
			Status:  http.StatusInternalServerError,
			ErrCode: http.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}

	var auctionItemAuthResponse AuctionItemResponse
	if err := c.parseResponse(resp, &auctionItemAuthResponse); err != nil {
		return nil, err
	}

	auctionItemAuthResponse.ResultSet.Result.Description = auctionItemAuthResponse.ResultSet.Result.GetDescription()

	calculateShippingFee(ctx, &auctionItemAuthResponse)

	return &auctionItemAuthResponse, nil
}

func calculateShippingFee(ctx context.Context, auctionItemAuthResponse *AuctionItemResponse) error {
	if auctionItemAuthResponse == nil {
		return nil
	}
	if auctionItemAuthResponse.ResultSet.Result.ChargeForShipping == "seller" {
		auctionItemAuthResponse.ResultSet.Result.ShippingFee = 0
		return nil
	}
	shipping := auctionItemAuthResponse.ResultSet.Result.Shipping
	if shipping == nil {
		auctionItemAuthResponse.ResultSet.Result.ShippingFee = DefaultShippingFee
		return nil
	}

	// load the shipping fee from the database
	// get from redis
	shippingFeeSetting := make(map[int]map[string]map[string]float64)
	if err := cache.GetRedisClient().Get(ctx, "supplysrv:yahoo:shipping_fee",
		&shippingFeeSetting); err != nil {
		// load from database
		shippingFees, err := db.GetHandler().GetShippingFee(ctx)
		if err != nil {
			return err
		}

		for _, shippingFee := range shippingFees {
			if _, ok := shippingFeeSetting[shippingFee.ServiceCode]; !ok {
				shippingFeeSetting[shippingFee.ServiceCode] = make(map[string]map[string]float64)
			}
			if _, ok := shippingFeeSetting[shippingFee.ServiceCode][shippingFee.From]; !ok {
				shippingFeeSetting[shippingFee.ServiceCode][shippingFee.From] = make(map[string]float64)
			}
			shippingFeeSetting[shippingFee.ServiceCode][shippingFee.From][shippingFee.Size] = shippingFee.Fee
		}
		// set to redis
		if err := cache.GetRedisClient().Set(ctx, "supplysrv:yahoo:shipping_fee", shippingFeeSetting, 60*time.Minute); err != nil {
			return err
		}
	}

	// get the lowest shipping fee
	lowestShippingFee := 1000000.0
	for _, method := range shipping.Method {
		if method.SinglePrice != 0 {
			if method.SinglePrice < lowestShippingFee {
				lowestShippingFee = method.SinglePrice
			}
			continue
		}

		switch method.ServiceCode {
		case 112:
			if lowestShippingFee > 230 {
				lowestShippingFee = 230
			}
			continue
		case 113:
			fee := shippingFeeSetting[method.ServiceCode][auctionItemAuthResponse.ResultSet.Result.Location]["0"]
			if fee == 0 {
				fee = DefaultShippingFee
			}
			if fee < lowestShippingFee {
				lowestShippingFee = fee
			}
			continue
		case 114:
			fee := shippingFeeSetting[method.ServiceCode][auctionItemAuthResponse.ResultSet.Result.Location][method.DeliveryFeeSize]
			if fee == 0 {
				fee = DefaultShippingFee
			}
			if fee < lowestShippingFee {
				lowestShippingFee = fee
			}
			continue
		case 115:
			fee := 0.0
			switch method.DeliveryFeeSize {
			case "0", "":
				fee = 230
			case "20":
				fee = 180
			case "50":
				fee = 440
			}
			if fee < lowestShippingFee {
				lowestShippingFee = fee
			}
			continue
		case 116:
			fee := shippingFeeSetting[method.ServiceCode][auctionItemAuthResponse.ResultSet.Result.Location][method.DeliveryFeeSize]
			if fee == 0 {
				fee = DefaultShippingFee
			}
			if fee < lowestShippingFee {
				lowestShippingFee = fee
			}
			continue
		}
	}

	if lowestShippingFee > 100000 {
		auctionItemAuthResponse.ResultSet.Result.ShippingFee = DefaultShippingFee
	} else {
		auctionItemAuthResponse.ResultSet.Result.ShippingFee = int64(lowestShippingFee)
	}

	return nil
}
