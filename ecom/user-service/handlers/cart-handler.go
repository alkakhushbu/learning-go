package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"user-service/internal/carts"
	"user-service/internal/consul"
	"user-service/pkg/ctxmanage"
	"user-service/pkg/logkey"

	"github.com/gin-gonic/gin"
)

const (
	OPEN   = "OPEN"
	CLOSED = "CLOSED"
)

func (h *Handler) CreateOrUpdateCart(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	//get auth claims from the context
	claims, err := ctxmanage.GetAuthClaimsFromContext(c.Request.Context())
	if err != nil {
		slog.Error(
			"missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised request"})
		return
	}

	userId := claims.Subject
	var newCartItem carts.NewCartItem
	if err := c.ShouldBindJSON(&newCartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart request payload"})
		return
	}

	ctx := c.Request.Context()

	type ProductServiceResponse struct {
		ProductID string `json:"product_id"`
		Stock     int    `json:"stock"`
		PriceID   string `json:"price_id"`
	}
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
		httpQuery := fmt.Sprintf("http://%s:%d/products/stock/%s", address, port, newCartItem.ProductID)
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

	// comare the stock of product with the quantity in cart request
	productStock := <-productChan
	if newCartItem.Quantity > productStock.Stock {
		slog.Error("Not enough products in stock")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Not enough products in stock"})
		return
	}

	//Get Cart for the user
	cart, err := h.c.GetCart(ctx, userId)
	if err != nil {
		//create a new cart if cart does not exist
		cart, err = h.c.InsertCart(ctx, userId)
		if err != nil {
			slog.Error("Error in creating cart", slog.Any("Error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in creating cart"})
			return
		}
	}

	//add items to the cart
	err = h.c.AddItemsToCart(ctx, cart.ID, newCartItem)
	if err != nil {
		slog.Error("Error in adding items to cart", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in adding items to cart"})
		return
	}

	slog.Info("Cart added:", slog.Any("Cart", cart))
	c.JSON(http.StatusOK, gin.H{"message": "Items added successfully"})
}
