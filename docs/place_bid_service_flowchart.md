# PlaceBidService Flowchart

```mermaid
flowchart TD
    Start([Start: PlaceBidService]) --> Validation
    Validation --> |Error|ValidationError
    Validation --> |Success|GetAuctionItem
    GetAuctionItem -->|Error|AuctionNotFoundError
    GetAuctionItem -->|Success| IsAuctionValid{Is Auction valid?}
    
    IsAuctionValid -->|No| InvalidAuctionError
    IsAuctionValid -->|Yes| InsertOrder
    
    InsertOrder -->|Error| SystemError
    InsertOrder -->|Success| CheckOrderStatus
    
    CheckOrderStatus -->|Error| OrderExistsError
    CheckOrderStatus -->|Success| PlaceBidPreview
    
    PlaceBidPreview -->|Error| UpdateBidRequestToPreviewError
    UpdateBidRequestToPreviewError --> PlaceBidPreviewError
    PlaceBidPreview -->|Success| IsShopItem{Is Shop item?}
    IsShopItem--> |Yes|ShpAucItem
    IsShopItem --> |No|PlaceBid
    ShpAucItem --> |Success| UpdateBidRequestToWinBid
    ShpAucItem --> |Error| UpdateBidRequestToShpAucItemError
    PlaceBid --> |Success| UpdateBidRequestToWinBid
    PlaceBid --> |Error| UpdateBidRequestToPlaceBidError
    UpdateBidRequestToWinBid --> Success

    UpdateBidRequestToPlaceBidError --> PlaceBidError
    UpdateBidRequestToShpAucItemError --> PlaceBidError

    ValidationError --> End
    AuctionNotFoundError --> End
    InvalidAuctionError --> End
    SystemError --> End
    OrderExistsError --> End
    PlaceBidPreviewError --> End
    PlaceBidError --> End
    Success --> End

    
    style ValidationError fill:#FF6B6B
    style AuctionNotFoundError fill:#FF6B6B
    style InvalidAuctionError fill:#FF6B6B
    style SystemError fill:#FF6B6B
    style Start fill:#90EE90
    style End fill:#FFB6C1
    style OrderExistsError fill:#FF6B6B
    style PlaceBidPreviewError fill:#FF6B6B
    style PlaceBidError fill:#FF6B6B
    style Success fill:#51CF66
```

## Flow Description

1. **Validation Phase**: Validates all required fields (YsRefID, TransactionType, AuctionID, Price, Quantity)
2. **Auction Item Retrieval**: Fetches the auction item and performs basic validations (status, price, quantity)
3. **Order Creation**: Creates and inserts the bid request order into the database with status "CREATED"
4. **Bid Preview**: Calls the preview API to get the signature required for placing the bid
5. **Transaction Type Handling**:
   - **BID**: Currently not supported, returns error
   - **BUYOUT**: Places the bid and updates order status to "WIN_BID" on success
6. **Error Handling**: Updates order status to "FAILED" with error message when operations fail

## Notes

- Several validation checks are marked as TODO and don't currently return errors
- The function only supports BUYOUT transactions currently
- Database operations are performed at multiple points to track order status
- All errors are logged before being returned

