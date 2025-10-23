package yahoo

import (
	"fmt"
	"io"
	"log"
)

// ExampleUsage demonstrates how to use the Yahoo Auction Bridge client
func ExampleUsage() {
	// Initialize client
	client := NewClient(
		"https://yahoo-auction-bridge.example.com", // Base URL
		"your-api-key",    // API Key
		"your-secret-key", // Secret Key
	)

	// Example 1: Health Check
	fmt.Println("=== Health Check ===")
	resp, err := client.HealthCheck()
	if err != nil {
		log.Printf("Health check failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Health check response: %s\n", string(body))

	// Example 2: Get Auction Item (Public API)
	fmt.Println("\n=== Get Auction Item ===")
	auctionReq := AuctionItemRequest{
		AuctionID: "x12345",
		AppID:     "your-app-id",
	}

	resp, err = client.GetAuctionItem(auctionReq)
	if err != nil {
		log.Printf("Get auction item failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		auctionItem, err := ParseAuctionItemResponse(resp)
		if err != nil {
			log.Printf("Failed to parse auction item: %v", err)
		} else {
			fmt.Printf("Auction Title: %s\n", auctionItem.ResultSet.Result.Title)
			fmt.Printf("Current Price: %d\n", auctionItem.ResultSet.Result.CurrentPrice)
		}
	} else {
		errorResp, _ := ParseErrorResponse(resp)
		fmt.Printf("Error: %+v\n", errorResp)
	}

	// Example 3: OAuth Authorization Flow
	fmt.Println("\n=== OAuth Authorization ===")
	authReq := OAuthAuthorizeRequest{
		YahooAccountID: "account123",
	}

	resp, err = client.Authorize(authReq)
	if err != nil {
		log.Printf("Authorization failed: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Authorization response status: %d\n", resp.StatusCode)
	fmt.Printf("Location header: %s\n", resp.Header.Get("Location"))

	// Example 4: Place Bid Preview
	fmt.Println("\n=== Place Bid Preview ===")
	previewReq := PlaceBidPreviewRequest{
		YahooAccountID:  "account123",
		YsRefID:         "YS-REF-001",
		TransactionType: "BID",
		AuctionID:       "x12345",
		Price:           1000,
		Quantity:        1,
		Partial:         false,
	}

	resp, err = client.PlaceBidPreview(previewReq)
	if err != nil {
		log.Printf("Bid preview failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Bid preview response: %s\n", string(body))

	// Example 5: Place Bid (requires signature from preview)
	fmt.Println("\n=== Place Bid ===")
	bidReq := PlaceBidRequest{
		YahooAccountID:  "account123",
		YsRefID:         "YS-REF-001",
		TransactionType: "BID",
		AuctionID:       "x12345",
		Price:           1000,
		Signature:       "signature-from-preview", // This should come from preview response
		Quantity:        1,
		Partial:         false,
	}

	resp, err = client.PlaceBid(bidReq)
	if err != nil {
		log.Printf("Place bid failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Place bid response: %s\n", string(body))

	// Example 6: Search Transactions
	fmt.Println("\n=== Search Transactions ===")
	searchReq := TransactionSearchRequest{
		YahooAccountID: "account123",
		StartDate:      "2024-01-01",
		EndDate:        "2024-12-31",
		Status:         "completed",
		Limit:          10,
		Offset:         0,
	}

	resp, err = client.SearchTransactions(searchReq)
	if err != nil {
		log.Printf("Search transactions failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Search transactions response: %s\n", string(body))

	// Example 7: Get Specific Transaction
	fmt.Println("\n=== Get Transaction ===")
	resp, err = client.GetTransaction("transaction123", "account123")
	if err != nil {
		log.Printf("Get transaction failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get transaction response: %s\n", string(body))

	// Example 8: Export Transactions CSV
	fmt.Println("\n=== Export Transactions CSV ===")
	exportReq := TransactionSearchRequest{
		YahooAccountID: "account123",
		StartDate:      "2024-01-01",
		EndDate:        "2024-12-31",
		Status:         "completed",
	}

	resp, err = client.ExportTransactionsCSV(exportReq)
	if err != nil {
		log.Printf("Export transactions failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Export CSV response (first 200 chars): %.200s...\n", string(body))
}

// ExampleBiddingFlow demonstrates the complete bidding process
func ExampleBiddingFlow() {
	client := NewClient(
		"https://yahoo-auction-bridge.example.com",
		"your-api-key",
		"your-secret-key",
	)

	// Step 1: Get auction item information
	fmt.Println("=== Step 1: Get Auction Item ===")
	auctionReq := AuctionItemRequest{
		AuctionID: "x12345",
	}

	resp, err := client.GetAuctionItem(auctionReq)
	if err != nil {
		log.Printf("Failed to get auction item: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Failed to get auction item, status: %d", resp.StatusCode)
		return
	}

	auctionItem, err := ParseAuctionItemResponse(resp)
	if err != nil {
		log.Printf("Failed to parse auction item: %v", err)
		return
	}

	fmt.Printf("Auction: %s\n", auctionItem.ResultSet.Result.Title)
	fmt.Printf("Current Price: %d\n", auctionItem.ResultSet.Result.CurrentPrice)
	fmt.Printf("Status: %s\n", auctionItem.ResultSet.Result.ItemStatus)

	// Step 2: Get bid preview with signature
	fmt.Println("\n=== Step 2: Get Bid Preview ===")
	previewReq := PlaceBidPreviewRequest{
		YahooAccountID:  "account123",
		YsRefID:         "YS-REF-001",
		TransactionType: "BID",
		AuctionID:       "x12345",
		Price:           auctionItem.ResultSet.Result.CurrentPrice + 100, // Bid 100 yen more
		Quantity:        1,
		Partial:         false,
	}

	resp, err = client.PlaceBidPreview(previewReq)
	if err != nil {
		log.Printf("Failed to get bid preview: %v", err)
		return
	}
	defer resp.Body.Close()

	// Parse preview response to get signature
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Bid preview response: %s\n", string(body))

	// In a real implementation, you would parse the JSON response to extract the signature
	signature := "extracted-signature-from-response"

	// Step 3: Place the actual bid
	fmt.Println("\n=== Step 3: Place Bid ===")
	bidReq := PlaceBidRequest{
		YahooAccountID:  "account123",
		YsRefID:         "YS-REF-001",
		TransactionType: "BID",
		AuctionID:       "x12345",
		Price:           previewReq.Price,
		Signature:       signature,
		Quantity:        1,
		Partial:         false,
	}

	resp, err = client.PlaceBid(bidReq)
	if err != nil {
		log.Printf("Failed to place bid: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Bid placed successfully: %s\n", string(body))
}
