# BidRequest Table Information

## Table: `yahoo.bid_request`

### Struct Fields

| Field Name | Go Type | Database Column | Nullable | Description |
|------------|---------|-----------------|----------|-------------|
| ID | `uint` | `id` | No | Primary key (from `gorm.Model`) |
| CreatedAt | `time.Time` | `created_at` | Yes | Record creation timestamp (from `gorm.Model`) |
| UpdatedAt | `time.Time` | `updated_at` | Yes | Record last update timestamp (from `gorm.Model`) |
| DeletedAt | `gorm.DeletedAt` | `deleted_at` | Yes | Soft delete timestamp (from `gorm.Model`) |
| RequestType | `string` | `request_type` | Yes | Type of transaction request (e.g., "BID", "BUYOUT") |
| OrderID | `string` | `order_id` | Yes | Unique order identifier (YsRefID) |
| AuctionID | `string` | `auction_id` | Yes | Yahoo Auction item ID |
| MaxBid | `int64` | `max_bid` | Yes | Maximum bid price for the order |
| Quantity | `int32` | `quantity` | Yes | Number of items requested |
| Partial | `bool` | `partial` | Yes | Whether partial quantity fulfillment is allowed |
| Status | `string` | `status` | Yes | Order status (e.g., "CREATED", "WIN_BID", "FAILED") |
| ErrorMessage | `string` | `error_message` | Yes | Error message if the order failed |

### Status Values

Common status values used in the system:
- `CREATED` - Order has been created and is pending processing
- `WIN_BID` - Bid was successful
- `FAILED` - Order processing failed

### RequestType Values

- `BID` - Regular bid transaction
- `BUYOUT` - Buyout transaction (immediate purchase)

### Indexes

Based on the DDL structure, the following indexes exist:
- Primary key on `id`
- Index on `order_id` (for lookups by order ID)
- Index on `auction_id` (for lookups by auction ID)
- Index on `status` (for filtering by status)
- Index on `deleted_at` (for soft delete queries)
- Index on `created_at` (for time-based queries)
- Index on `updated_at` (for time-based queries)

### Relationships

- One `BidRequest` can have multiple `YahooTransaction` records (linked via `BidRequestID`)


