package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Struct Definitions

type Response struct {
	Data    map[string][]CartProductDetails `json:"data"`
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
}

type CartProductDetails struct {
	ID                  int    `json:"id"`
	CartID              int    `json:"cart_id"`
	ProductID           int    `json:"product_id"`
	Quantity            int    `json:"quantity"`
	Name                string `json:"name"`
	BrandName           string `json:"brand_name"`
	CategoryName        string `json:"category_name"`
	Price               int    `json:"price"` // Changed to int
	OfferPrice          int    `json:"offer_price"`
	OfferDiscount       string `json:"offer_discount"`
	SelectedSizeVariant string `json:"selectedSizeVariant"`
}

// Mock HTTP Handler
type HttpHandler struct{}

func (h *HttpHandler) Get(ctx *ContextStruct, url string, headers Headers) ([]byte, error) {
	// Mock JSON response (as per your initial data)
	jsonData := `{
        "data": {
            "9663735": [
                {
                    "id": 1642,
                    "cart_id": 26,
                    "product_id": 345341,
                    "quantity": 1,
                    "name": "Garnier Bright Complete VITAMIN C Facewash (150 g)",
                    "brand_name": "Garnier",
                    "category_name": "Face Wash",
                    "price": 299,
                    "offer_price": 299,
                    "offer_discount": "0",
                    "selectedSizeVariant": "150g"
                }
            ],
            "9663744": [
                {
                    "id": 622,
                    "cart_id": 180,
                    "product_id": 338359,
                    "quantity": 1,
                    "name": "Good Vibes Ubtan De-tan Glow Light Day Cream with Power of Serum (50 g)",
                    "brand_name": "Good Vibes",
                    "category_name": "Day Cream",
                    "price": 275,
                    "offer_price": 275,
                    "offer_discount": "0",
                    "selectedSizeVariant": "50 g"
                }
            ]
        },
        "success": true,
        "message": "Fetched data successfully"
    }`

	return []byte(jsonData), nil
}

var http_handler = &HttpHandler{}

// Context and Headers Structs
type ContextStruct struct {
	SpanLabelPrefix string
}

type Headers struct {
	IdleConnTimeout time.Duration
}

// Logger Mock
var logger = struct {
	Log func(level, message, context, file, packageName string, err error)
}{
	Log: func(level, message, context, file, packageName string, err error) {
		log.Printf("[%s] %s | Context: %s | File: %s | Package: %s | Error: %v\n",
			level, message, context, file, packageName, err)
	},
}

// GetCartDetails Function (Revised)
func GetCartDetails(ctx *ContextStruct, userIDs []int32) (map[int32][]CartProductDetails, error) {
	cartDetailsMap := make(map[int32][]CartProductDetails)

	// Mock URL for demonstration
	URL := "http://mock-catalogue-url/neo/cart-items/v1"

	ctx.SpanLabelPrefix = "Catalogue.Products"
	responseBytes, err := http_handler.Get(ctx, URL, Headers{IdleConnTimeout: 3 * time.Second})
	ctx.SpanLabelPrefix = ""

	if err != nil {
		logger.Log("Error", "HTTP GET request failed", "GetCartDetails", "main.go", "main", err)
		return nil, err
	}

	if len(responseBytes) == 0 {
		return nil, errors.New("empty response received")
	}

	var respData Response
	err = json.Unmarshal(responseBytes, &respData)
	if err != nil {
		logger.Log("Error", "Error while unmarshaling Catalogue response", "GetCartDetails", "main.go", "main", err)
		return nil, err
	}

	if !respData.Success {
		return nil, errors.New(respData.Message)
	}

	// Convert userID to string to match the keys in the Data map
	for _, userID := range userIDs {
		userKey := strconv.Itoa(int(userID))
		products, exists := respData.Data[userKey]
		if !exists {
			// Log warning and continue processing other userIDs
			logger.Log("Warning", fmt.Sprintf("No cart items found for userID %d", userID), "GetCartDetails", "main.go", "main", nil)
			continue
		}

		fmt.Printf("Processing Key: %s\n", userKey)
		cartDetailsMap[userID] = products
	}

	if len(cartDetailsMap) == 0 {
		return nil, errors.New("no cart items found for the provided userIDs")
	}

	return cartDetailsMap, nil
}

func main() {
	ctx := &ContextStruct{}
	userIDs := []int32{9663744, 9663735, 12345} // Example userID

	fmt.Println(userIDs)

	cartDetailsMap, err := GetCartDetails(ctx, userIDs)
	if err != nil {
		log.Fatalf("Failed to get cart details: %v", err)
	}

	// Print the cart details
	for _, userID := range userIDs {
		products, exists := cartDetailsMap[userID]
		if !exists {
			fmt.Printf("\nNo cart items found for userID: %d\n", userID)
			continue
		}

		fmt.Printf("\nCart Details for UserID: %d\n", userID)
		for _, product := range products {
			fmt.Printf("Product ID: %d\n", product.ProductID)
		}
	}
}
