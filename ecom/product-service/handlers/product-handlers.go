package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"product-service/internal/products"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	slog.Info("Inside CreateProduct")
	if c.Request.ContentLength > 5*1024 {
		// Log error for payload exceeding size limit
		slog.Error("json validation error", slog.String("Error", "request body limit breached"))
		// Return a 400 Bad Request status code along with an error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "payload exceeding size limit"})
		return
	}

	var newProduct products.NewProduct
	err := c.ShouldBindBodyWithJSON(&newProduct)
	if err != nil {
		slog.Error("json validation error", slog.String("Error", err.Error()))

		// Respond with a 400 Bad Request status code and error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	err = h.validate.Struct(newProduct)
	if err != nil {
		slog.Error("json validation error", slog.String("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "one or more json validation failed"})
		return
	}

	price, err := ValidatePrice(newProduct.Price)
	if err != nil {
		slog.Error("invalid price error", slog.String("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.conf.InsertProduct(c.Request.Context(), newProduct)
	if err != nil {
		slog.Error("product insertion failed error", slog.String("Error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error creating new product"})
		return
	}

	// create product on stripe and update product service db as well
	// create a new context for registering product on stripe and keeping a copy of it locally
	// the request context gets cancelled with response return
	go h.conf.CreateProductStripe(context.Background(), price, product.ID, product.Name)
	c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProductInfo(c *gin.Context) {
	productId := c.Param("productId")

	productInfo, err := h.conf.GetProductInfo(c.Request.Context(), productId)
	if err != nil {
		slog.Error("Error in fetching product Info", slog.String("productId", productId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error in fetching product Info"})
		return
	}
	c.JSON(http.StatusOK, productInfo)
}

func (h *Handler) GetProductInfos(c *gin.Context) {
	var request struct {
		ProductIds []string `json:"productIds"`
	}

	//validate Productinfo
	err := h.validate.Struct(request)
	if err != nil {
		slog.Error("Invalid format of cart request payload", slog.Any("Error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format of product request payload"})
		return
	}
	if err := c.BindJSON(&request); err != nil {
		slog.Error("Error in binding JSON payload", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	productInfos, err := h.conf.GetProductInfos(c.Request.Context(), request.ProductIds)
	if err != nil {
		slog.Error("Error in fetching product Infos", slog.Any("productIds", request.ProductIds))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error in fetching product Infos"})
		return
	}

	c.JSON(http.StatusOK, productInfos)
}

/*
1. Trim extra spaces
2. split the price based by dot(.)
3. The size of the split slice should be either 1 or 2
4. Rupee part : priceSlice[0], paisa part: priceSlice[1]
5. Convert rupee part into int and multiply by 100
6. Size of paisa part should be either 1 or 2
6.1 if paisa part size is 1, append a zero at the end and covert to integer
7. Add paisa to rupee
*/
func ValidatePrice(priceStr string) (uint64, error) {
	//trim extra space from price
	priceStr = strings.Trim(priceStr, " ")

	//split the price based by dot(.)
	prices := strings.Split(priceStr, ".")
	var rupee, paisa uint64
	if len(prices) == 0 || len(prices) > 2 {
		return 0, fmt.Errorf("invalid price, please provide price in valid format")
	}

	rupee, err := strconv.ParseUint(prices[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price, please provide price in valid format")
	}

	if len(prices) == 2 {
		// check size of paisa part for edge cases like 11.001 or 11.011
		// we want upto 2 digits of paisa in INR price
		if len(prices[1]) > 2 {
			return 0, fmt.Errorf("invalid price, please provide price in valid format")
		}
		paisa, err = strconv.ParseUint(prices[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid price, please provide price in valid format")
		}

		// append 0 if paisa part has only one digit
		// e.g INR 99.2 => Convert it to 9900 + 20 = 9920
		// convert 99.09 into 9900
		if len(prices[1]) == 1 {
			paisa *= 10
		}
	}
	return rupee*100 + paisa, nil
}
