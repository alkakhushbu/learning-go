package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"order-service/internal/auth"
	"order-service/internal/consul"
	"order-service/pkg/ctxmanage"
	"order-service/pkg/logkey"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

func (h *Handler) ProductCheckout(c *gin.Context) {

	// Get the traceId from the request for tracking logs
	traceId := ctxmanage.GetTraceIdOfRequest(c)
	claims, ok := c.Request.Context().Value(auth.ClaimsKey).(auth.Claims)
	if !ok {
		slog.Error("claims not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusUnauthorized})
		return
	}

	type UserServiceResponse struct {
		StripCustomerId string `json:"stripe_customer_id"`
	}
	type ProductServiceResponse struct {
		ProductID string `json:"product_id"`
		Stock     int    `json:"stock"`
		PriceID   string `json:"price_id"`
	}

	productID := c.Param("productID")
	if productID == "" {
		slog.Error("missing product id", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	// Create channels for user-service goroutine results
	userChan := make(chan UserServiceResponse, 1) // For customer ID

	if h.client == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "consul client is not initialized"})
	}

	// call user-service endpoint here
	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "users")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/users/stripe", address, port)
		slog.Info("httpQuery: "+httpQuery, slog.String(logkey.TraceID, traceId))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpQuery, nil)
		if err != nil {
			slog.Error("error creating request", slog.String(logkey.TraceID, traceId), slog.Any("error", err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		authorizationHeader := c.Request.Header.Get("Authorization")
		req.Header.Set("Authorization", authorizationHeader)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("error fetching user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching stripe id from user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}

		defer resp.Body.Close()

		var userServiceResponse UserServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&userServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		// Print the customer Id if fetched successfully
		slog.Info("successfully fetched stripe customer id", slog.String(logkey.TraceID, traceId))
		userChan <- userServiceResponse
	}()

	// Create channels for product-service goroutine results
	productChan := make(chan ProductServiceResponse, 1) // For stock and price information

	// call product-service endpoint here
	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "products")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			productChan <- ProductServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/products/stock/%s", address, port, productID)
		resp, err := http.Get(httpQuery)
		if err != nil {
			slog.Error("error fetching product service", slog.String(logkey.TraceID, traceId))
			productChan <- ProductServiceResponse{}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching product information", slog.String(logkey.TraceID, traceId))
			productChan <- ProductServiceResponse{}
			return
		}
		var productServiceResponse ProductServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&productServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- ProductServiceResponse{}
			return
		}
		productChan <- productServiceResponse
	}()

	userServiceResponse := <-userChan
	if userServiceResponse.StripCustomerId == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching stripe customer id"})
		return
	}
	stockPriceData := <-productChan
	priceID := stockPriceData.PriceID
	stock := stockPriceData.Stock
	if stock <= 0 || priceID == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching product information"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"customerId": userServiceResponse.StripCustomerId, "price_id": priceID, "stock": stock})

	// Step 1: Retrieve the Stripe secret key from the environment variables
	sKey := os.Getenv("STRIPE_TEST_KEY")
	if sKey == "" {
		slog.Error("Stripe secret key not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Stripe secret key not found"})
	}

	// Step 2: Assign the Stripe API key to the Stripe library's internal configuration
	stripe.Key = sKey
	orderId := uuid.NewString()
	// Proceed to create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		Customer:                 stripe.String(userServiceResponse.StripCustomerId),
		SubmitType:               stripe.String("pay"),
		Currency:                 stripe.String(string(stripe.CurrencyINR)),
		BillingAddressCollection: stripe.String("auto"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1), // Adjust quantity as needed
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://example.com/success"),
		//ExpiresAt:
		CancelURL: stripe.String("https://example.com/cancel"),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_id":   orderId,
				"user_id":    claims.Subject, // userID in jwt token
				"product_id": productID,
			},
		},
	}

	//make api-call to create new payment session on stripe
	sessionStripe, err := session.New(params)
	if err != nil {
		slog.Error("error creating Stripe checkout session", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Stripe checkout session"})
		return
	}

	// Log success operation
	slog.Info("successfully initiated Stripe checkout session", slog.String("Trace ID", traceId), slog.String("ProductID", productID), slog.String("CheckoutSessionID", sessionStripe.ID))

	userId := claims.Subject
	ctx := c.Request.Context()
	err = h.dbConf.CreateOrder(ctx, orderId, userId, productID, sessionStripe.AmountTotal)
	if err != nil {
		slog.Error("error creating order", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})

}

func (h *Handler) CartCheckout(c *gin.Context) {
	/**
	1. Make DB call and get priceId, product, quantity etc.
	2. Check product stock on product service, update order event accordingly
	3. Create order create event, update DB, and publish to kafka {orderId, cartId, []carts_items}
		3.1 Consumed by cart service => update cart
	4. Make stripe API call for payment

	5. Make payment
	6. Payment success

	7. stripe callback on Webhook API
	8. Create order paid event, update DB, and publish to kafka {orderId, cartId, []carts_items}
		8.1 Consume by cart service => update cart
		8.2 Consume by product service => update stock
	**/
	traceId := ctxmanage.GetTraceIdOfRequest(c)
	claims, ok := c.Request.Context().Value(auth.ClaimsKey).(auth.Claims)
	if !ok {
		slog.Error("claims not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusUnauthorized})
		return
	}

	type UserServiceResponse struct {
		StripCustomerId string `json:"stripe_customer_id"`
	}
	type ProductServiceResponse struct {
		ProductID string `json:"product_id"`
		Stock     int    `json:"stock"`
		PriceID   string `json:"price_id"`
	}

	productID := c.Param("productID")

	//get auth claims from the context
	ctx := c.Request.Context()
	claims, err := ctxmanage.GetAuthClaimsFromContext(ctx)
	if err != nil {
		slog.Error(
			"missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised request"})
		return
	}

	userId := claims.Subject

	cartItems, err := h.c.GetAllCartItems(ctx, userId)
	if err != nil {
		slog.Error("error fetching cart details", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching cart details"})
		return
	}

	var productIds []string
	productQuantityMap := make(map[string]int64)
	for _, item := range cartItems.CartItems {
		productIds = append(productIds, item.ProductID)
		productQuantityMap[item.ProductID] = int64(item.Quantity)
	}

	// Create channels for user-service goroutine results
	userChan := make(chan UserServiceResponse, 1) // For customer ID

	if h.client == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "consul client is not initialized"})
	}

	// call user-service endpoint here
	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "users")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/users/stripe", address, port)
		slog.Info("httpQuery: "+httpQuery, slog.String(logkey.TraceID, traceId))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpQuery, nil)
		if err != nil {
			slog.Error("error creating request", slog.String(logkey.TraceID, traceId), slog.Any("error", err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		authorizationHeader := c.Request.Header.Get("Authorization")
		req.Header.Set("Authorization", authorizationHeader)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("error fetching user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching stripe id from user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}

		defer resp.Body.Close()

		var userServiceResponse UserServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&userServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		// Print the customer Id if fetched successfully
		slog.Info("successfully fetched stripe customer id", slog.String(logkey.TraceID, traceId))
		userChan <- userServiceResponse
	}()

	// Create channels for product-service goroutine results
	productChan := make(chan []ProductServiceResponse, 1) // For stock and price information

	// call product-service endpoint here
	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "products")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			productChan <- nil
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/products/stock", address, port)
		productIdMap := map[string][]string{
			"productIds": productIds,
		}
		jsonData, err := json.Marshal(productIdMap)
		if err != nil {
			slog.Error("error marshalling cart items", slog.String(logkey.TraceID, traceId))
			productChan <- nil
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpQuery, bytes.NewBuffer(jsonData))
		if err != nil {
			slog.Error("error creating request", slog.String(logkey.TraceID, traceId), slog.Any("error", err.Error()))
			productChan <- nil
			return
		}

		authorizationHeader := c.Request.Header.Get("Authorization")
		req.Header.Set("Authorization", authorizationHeader)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("error", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- nil
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching product information", slog.String(logkey.TraceID, traceId))
			productChan <- nil
			return
		}
		var productInfos []ProductServiceResponse
		bytesResp, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bytesResp), "****************************")
		err = json.Unmarshal(bytesResp, &productInfos)
		// err = json.NewDecoder(resp.Body).Decode(&productInfos)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- nil
			return
		}
		productChan <- productInfos
	}()

	userServiceResponse := <-userChan
	if userServiceResponse.StripCustomerId == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching stripe customer id"})
		return
	}
	stockPrices := <-productChan
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, stockPriceData := range stockPrices {
		priceID := stockPriceData.PriceID
		stock := stockPriceData.Stock
		if stock <= 0 || priceID == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching product information"})
			return
		}
		lineItems = append(lineItems,
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(productQuantityMap[stockPriceData.ProductID]), // Adjust quantity as needed
			},
		)
	}
	slog.Info("Line Item", slog.Any("Line Item", lineItems))
	slog.Info("stockPrices", slog.Any("stockPrices", stockPrices))

	//c.JSON(http.StatusOK, gin.H{"customerId": userServiceResponse.StripCustomerId, "price_id": priceID, "stock": stock})

	cartItemsJson, err := json.Marshal(cartItems.CartItems)
	if err != nil {
		slog.Error("Error in marshaling cartItems", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error in marshaling cartItems"})
		return
	}
	// Step 1: Retrieve the Stripe secret key from the environment variables
	sKey := os.Getenv("STRIPE_TEST_KEY")
	if sKey == "" {
		slog.Error("Stripe secret key not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Stripe secret key not found"})
	}

	// Step 2: Assign the Stripe API key to the Stripe library's internal configuration
	stripe.Key = sKey
	orderId := uuid.NewString()
	// Proceed to create Stripe checkout session

	params := &stripe.CheckoutSessionParams{
		Customer:                 stripe.String(userServiceResponse.StripCustomerId),
		SubmitType:               stripe.String("pay"),
		Currency:                 stripe.String(string(stripe.CurrencyINR)),
		BillingAddressCollection: stripe.String("auto"),
		LineItems:                lineItems,
		Mode:                     stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:               stripe.String("https://example.com/success"),
		//ExpiresAt:
		CancelURL: stripe.String("https://example.com/cancel"),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_id":   orderId,
				"user_id":    claims.Subject, // userID in jwt token
				"product_id": productID,
				"cart_id":    cartItems.CartID,
				"cart_items": string(cartItemsJson),
			},
		},
	}

	//make api-call to create new payment session on stripe
	sessionStripe, err := session.New(params)
	if err != nil {
		slog.Error("error creating Stripe checkout session", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Stripe checkout session"})
		return
	}

	// Log success operation
	slog.Info("successfully initiated Stripe checkout session", slog.String("Trace ID", traceId), slog.String("ProductID", productID), slog.String("CheckoutSessionID", sessionStripe.ID))

	err = h.dbConf.CreateOrderForCart(ctx, orderId, userId, cartItems.CartID, sessionStripe.AmountTotal)
	if err != nil {
		slog.Error("error creating order", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order"})
		return
	}

	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
}

// func (h *Handler) produceEvents(orderId, userId, productId string) {
// 	orderEvent := struct {
// 		OrderId   string    `json:"order_id"` // UUID
// 		ProductId string    `json:"product_id"`
// 		Quantity  int       `json:"quantity"`
// 		UserId    string    `json:"user_id"`
// 		CreatedAt time.Time `json:"created_at"` // Timestamp of creation
// 	}{OrderId: orderId, UserId: userId, ProductId: productId, Quantity: 1}
// 	jsonData, _ := json.Marshal(orderEvent)
// 	key := []byte(orderId)
// 	err := h.kafkaConf.ProduceMessage(kafka.TopicOrderPaid, key, jsonData)
// 	if err != nil {
// 		slog.Error("Failed to produce message", slog.Any("error", err.Error()))
// 		return
// 	}
// 	slog.Info("Message produced", slog.Any("data", string(jsonData)))
// }

/*
*
                  +---------------------+
                  |       START         |
                  |   API Call Starts   |
                  +---------------------+
                            |
                            v
                +----------------------+
                |   Check STRIPE Key   |
                | (Environment Config) |
                +----------------------+
                            |
           STRIPE Key FOUND | STRIPE Key MISSING
                  |                   v
                  v       +----------------------------+
    +------------------+  | Respond with Error:         |
    | Extract User     |  | "Stripe test key not found" |
    | Claims & TraceID |  +----------------------------+
    +------------------+
                  |
                  v
          +-------------------+
          | Extract ProductID |
          |   From Request    |
          +-------------------+
                  |
        ProductID FOUND | 				ProductID MISSING
                  |                   			v
                  v       					+--------------------------------+
	+----------------------------------+	| Respond with Error: 			 |
    | Create Channels for Concurrent  	|	| "Product ID Missing"			 |
	|     Service Calls           		|	+--------------------------------+
    +----------------------------------+
                  |
                  v
   +---------------------------------------+
   | Start Parallel Service Calls          |
   | 1. Call User Service (Stripe ID)      |
   | 2. Call Product Service (Stock/Price) |
   +---------------------------------------+
                  |
          +-------------------+  +-------------------+
          | Wait for User ID  |  | Wait for Product   |
          +-------------------+  | Details            |
                  |               +-------------------+
       Stripe ID FOUND |   Product Details FOUND
                  |                    |
                  v                    v
       +---------------------------------------+
       | Validate Results:                    |
       | - Valid Stripe Customer ID           |
       | - Valid Product Details (Stock > 0,  |
       |   PriceID Exists)                    |
       +---------------------------------------+
                  |
          Validation PASSED | Validation FAILED
                  |                   	v
                  v      				----------------------------+
       +----------------------------+ 	Respond with Error    |
       | Create Stripe Checkout     | 	"Invalid Inputs"      |
       | Session with User & Product|	-----------------------+
       +----------------------------+
                  |// Create the order in the orders table with a pending status
                  v
     +--------------------------------+
     |  Respond with Checkout URL     |
     |  (Stripe Session Created)      |
     +--------------------------------+
                  |
                  v
           +-------------+
           |     END     |
           +-------------+
*/
