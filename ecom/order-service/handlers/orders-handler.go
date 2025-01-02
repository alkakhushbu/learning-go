package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"order-service/internal/consul"
	"order-service/pkg/ctxmanage"
	"order-service/pkg/logkey"
	"time"
)

func (h* Handler) Checkout(c *gin.Context) {
	//TODO: create a struct to handle response from the userservice
	//TODO: Create a function that returns service address and port
	//TODO: Make a request to user-service to fetch the stripe customer id
	// 	and unmarshal that into the struct created in step 1
	//TODO: authorizationHeader := c.Request.Header.Get("Authorization")
	//    req.Header.Set("Authorization", authorizationHeader)
	// Print the customer Id if fetched successfully

	// Get the traceId from the request for tracking logs
	traceId := ctxmanage.GetTraceIdOfRequest(c)

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

	// Create channels for goroutine results
	userChan := make(chan UserServiceResponse, 1) // For customer ID

	if h.client == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "consul client is not initialized"})
	}

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

	productChan := make(chan ProductServiceResponse, 1) // For stock and price information
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
	c.JSON(http.StatusOK, gin.H{"customerId": userServiceResponse.StripCustomerId, "price_id": priceID, "stock": stock})
	//c.JSON(http.StatusOK, gin.H{"stripe_customer_id": userServiceResponse.StripCustomerId})
}


