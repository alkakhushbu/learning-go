package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/internal/carts"
	"order-service/internal/consul"
	"order-service/pkg/ctxmanage"
	"order-service/pkg/logkey"

	"github.com/gin-gonic/gin"
)

const (
	OPEN   = "OPEN"
	CLOSED = "CLOSED"
)

func (h *Handler) AddToCart(c *gin.Context) {
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
		slog.Error("Invalid cart request payload", slog.Any("Error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart request payload"})
		return
	}

	//validate cartItem
	err = h.validate.Struct(newCartItem)
	if err != nil {
		slog.Error("Invalid format of cart request payload", slog.Any("Error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format of cart request payload"})
		return
	}
	ctx := c.Request.Context()

	// Create channels for product-service goroutine results
	productChan := make(chan carts.ProductServiceResponse, 1) // For stock and price information

	// call product-service endpoint here
	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "products")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			productChan <- carts.ProductServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/products/stock/%s", address, port, newCartItem.ProductID)
		resp, err := http.Get(httpQuery)
		if err != nil {
			slog.Error("error fetching product service", slog.String(logkey.TraceID, traceId))
			productChan <- carts.ProductServiceResponse{}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching product information", slog.String(logkey.TraceID, traceId))
			productChan <- carts.ProductServiceResponse{}
			return
		}
		var productServiceResponse carts.ProductServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&productServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- carts.ProductServiceResponse{}
			return
		}
		productChan <- productServiceResponse
	}()

	// comare the stock of product with the quantity in cart request
	productStock := <-productChan
	stock := productStock.Stock
	// priceId := productStock.PriceID
	if newCartItem.Quantity > stock {
		slog.Error("Not enough products in stock")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Not enough products in stock"})
		return
	}

	//Get Cart for the user
	cart, err := h.c.GetCart(ctx, userId)
	if err != nil {
		//create a new cart if cart does not exist
		fmt.Println("Create a new cart")
		cart, err = h.c.InsertCart(ctx, userId)
		if err != nil {
			slog.Error("Error in creating cart", slog.Any("Error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in creating cart"})
			return
		}
	}

	fmt.Println("Got cart id from user")
	//add items to the cart
	fmt.Println("Add items to cart..............")
	err = h.c.AddItemsToCart(ctx, cart.ID, newCartItem, productStock)
	if err != nil {
		if errors.Is(err, carts.ErrNotEnoughStock) {
			slog.Info(carts.ErrNotEnoughStock.Error())
			c.JSON(http.StatusOK, gin.H{"message": carts.ErrNotEnoughStock.Error()})
			return
		}
		slog.Error("Error in adding items to cart", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in adding items to cart"})
		return
	}

	slog.Info("Cart added:", slog.Any("Cart", cart))
	c.JSON(http.StatusOK, gin.H{"message": "Items added successfully"})
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
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
		slog.Error("Invalid cart request payload", slog.Any("Error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart request payload"})
		return
	}

	//validate cartItem
	err = h.validate.Struct(newCartItem)
	if err != nil {
		slog.Error("Invalid format of cart request payload", slog.Any("Error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format of cart request payload"})
		return
	}
	ctx := c.Request.Context()

	//Get Cart for the user
	cart, err := h.c.GetCart(ctx, userId)
	if err != nil {
		slog.Error("Cart does not exist for the user", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": "error in removing items from cart"})
		return
	}

	//remove items from the cart
	err = h.c.RemoveItemsFromCart(ctx, cart.ID, newCartItem)
	if err != nil {
		if errors.Is(err, carts.ErrItemNotInCart) {
			c.JSON(http.StatusOK, gin.H{"error": carts.ErrItemNotInCart.Error()})
			return
		}
		slog.Error("Error in removing items from cart", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in removing items from cart"})
		return
	}

	slog.Info("Items renoved:", slog.Any("Cart", cart))
	c.JSON(http.StatusOK, gin.H{"message": "Items removed successfully"})
}

func (h *Handler) GetAllCartItems(c *gin.Context) {
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
	ctx := c.Request.Context()
	cartResponse, err := h.c.GetAllCartItems(ctx, userId)
	if err != nil {
		if errors.Is(err, carts.ErrEmptyCart) {
			slog.Error("Error in fetching cart items", slog.Any("Error", err.Error()))
			c.JSON(http.StatusNoContent, gin.H{"message": err.Error()})
			return
		}
		slog.Error("Error in fetching cart items", slog.Any("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Error in fetching cart items"})
		return
	}
	c.JSON(http.StatusOK, cartResponse)
}
