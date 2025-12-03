# YahooTransaction Table Information

## Table: `yahoo.transaction`

### Struct Fields

| Field Name | Go Type | Database Column | Nullable | Description |
|------------|---------|-----------------|----------|-------------|
| ID | `uint` | `id` | No | Primary key (from `gorm.Model`) |
| CreatedAt | `time.Time` | `created_at` | Yes | Record creation timestamp (from `gorm.Model`) |
| UpdatedAt | `time.Time` | `updated_at` | Yes | Record last update timestamp (from `gorm.Model`) |
| DeletedAt | `gorm.DeletedAt` | `deleted_at` | Yes | Soft delete timestamp (from `gorm.Model`) |
| BidRequestID | `string` | `bid_request_id` | Yes | Reference to the associated bid request (OrderID) |
| Price | `int64` | `price` | Yes | Transaction price |
| Status | `string` | `status` | Yes | Transaction status |
| ErrorMessage | `string` | `error_message` | Yes | Error message if the transaction failed |

### Status Values

Common status values used in the system:
- `CREATED` - Transaction has been created
- `WIN_BID` - Transaction was successful
- `FAILED` - Transaction processing failed

### Indexes

Based on the DDL structure, the following indexes exist:
- Primary key on `id`
- Index on `bid_request_id` (for lookups by bid request ID)
- Index on `status` (for filtering by status)
- Index on `deleted_at` (for soft delete queries)
- Index on `created_at` (for time-based queries)
- Index on `updated_at` (for time-based queries)

### Relationships

- Many `YahooTransaction` records can belong to one `BidRequest` (linked via `BidRequestID` = `BidRequest.OrderID`)
- This table tracks transaction history for bid requests

### Usage Notes

- Each bid request can have multiple transaction records as the status changes
- The `BidRequestID` field links to the `order_id` in the `yahoo.bid_request` table
- Transaction records are created when bid requests are processed or updated


