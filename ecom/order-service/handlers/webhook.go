package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/internal/orders"
	"order-service/internal/stores/kafka"
	"order-service/pkg/logkey"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

func (h *Handler) Webhook(c *gin.Context) {
	traceId := uuid.NewString()
	fmt.Println(traceId)

	const MaxBodyBytes = int64(65536)

	// Limit the request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	var event stripe.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		slog.Error("Failed to bind JSON", slog.Any("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(event.Type, "********")
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			slog.Error("Failed to unmarshal JSON", slog.Any("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		slog.Info("Payment Intent Succeeded", slog.Any("paymentIntent ID", paymentIntent.ID))
		orderId := paymentIntent.Metadata["order_id"]
		productID := paymentIntent.Metadata["product_id"]
		userID := paymentIntent.Metadata["user_id"]
		cartId := paymentIntent.Metadata["cart_id"]
		slog.Info("Metadata received", slog.String(logkey.TraceID, traceId), slog.String("OrderID", orderId), slog.String("UserID", userID), slog.String("ProductID", productID))
		if productID != "" {
			go func() {
				jsonData, err := json.Marshal(kafka.OrderPaidEvent{
					OrderId:   orderId,
					ProductId: productID,
					Quantity:  1,
					CreatedAt: time.Now().UTC(),
				})

				if err != nil {
					slog.Error("Failed to marshal JSON", slog.Any("error", err.Error()))
					return
				}
				key := []byte(orderId)
				err = h.kafkaConf.ProduceMessage(kafka.TopicOrderPaid, key, jsonData)
				if err != nil {
					slog.Error("Failed to produce message", slog.Any("error", err.Error()))
					return
				}
				slog.Info("Message produced", slog.Any("data", string(jsonData)))
			}()
		} else {
			cartItemsString := paymentIntent.Metadata["cart_items"]

			var cartItems []kafka.CartItem
			err := json.Unmarshal([]byte(cartItemsString), &cartItems)
			if err != nil {
				slog.Error("Error in marshaling cartItems", slog.Any("error", err.Error()))
			}

			jsonData, err := json.Marshal(kafka.OrderPaidEvent{
				OrderId:   orderId,
				ProductId: productID,
				CartItems: cartItems,
				CartId:    cartId,
				CreatedAt: time.Now().UTC(),
			})
			key := []byte(orderId)
			err = h.kafkaConf.ProduceMessage(kafka.TopicOrderPaid, key, jsonData)
			if err != nil {
				slog.Error("Failed to produce message", slog.Any("error", err.Error()))
				return
			}
			slog.Info("Message produced", slog.Any("data", string(jsonData)))

		}
		ctx := c.Request.Context()
		err = h.dbConf.UpdateOrder(ctx, orderId, orders.StatusPaid, paymentIntent.ID)
		if err != nil {
			slog.Error("Failed to update order", slog.Any("error", err.Error()))
			return
		}

		err = h.c.UpdateCartStatus(ctx, cartId, "CLOSED")
		if err != nil {
			slog.Error("Failed to update cart", slog.Any("error", err.Error()))
			return
		}
		c.Status(http.StatusOK)
	}
}
