package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"user-service/internal/auth"
	"user-service/internal/stores/kafka"
	"user-service/internal/users"
	"user-service/middleware"
	"user-service/pkg/ctxmanage"
	"user-service/pkg/logkey"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Signup(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	//1 KB = 1024 bytes
	// Check if the size of the request body is more than 5KB
	if c.Request.ContentLength > 5*1024 {
		// Log error for payload exceeding size limit
		slog.Error("request body limit breached", slog.String(logkey.TraceID, traceId), slog.Int64("Size Received", c.Request.ContentLength))

		// Return a 400 Bad Request status code along with an error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "payload exceeding size limit"})
		return
	}

	var newUser users.NewUser
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		// Log error and associate it with a trace id for easy correlation
		slog.Error("json validation error", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))

		// Respond with a 400 Bad Request status code and error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	err = h.validate.Struct(newUser)
	if err != nil {
		slog.Error("validation failed", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.u.InsertUser(ctx, newUser)
	if err != nil {
		slog.Error("error in creating the user", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User Creation Failed"})
		return
	}

	go func() {
		data, err := json.Marshal(user)
		if err != nil {
			slog.Error("error in marshaling user", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			return
		}
		key := []byte(user.ID)
		err = h.k.ProduceMessage(kafka.TopicAccountCreated, key, data)
		if err != nil {
			slog.Error("error in producing message", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			return
		}
	}()
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Login(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	//Validate request payload size
	if c.Request.ContentLength > 5*1024 {
		// Log error for payload exceeding size limit
		slog.Error("request body limit breached", slog.String(logkey.TraceID, traceId), slog.Int64("Size Received", c.Request.ContentLength))

		// Return a 400 Bad Request status code along with an error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "payload exceeding size limit"})
		return
	}

	//validate json binding of request payload
	var loginUser users.LoginUser
	err := c.ShouldBindJSON(&loginUser)
	if err != nil {
		// Log error and associate it with a trace id for easy correlation
		slog.Error("json validation error", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))

		// Respond with a 400 Bad Request status code and error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	//validate the fields of request payload
	err = h.validate.Struct(loginUser)
	if err != nil {
		slog.Error("validation failed", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	//validate user credentials, match user email and password
	user, err := h.u.ValidateUser(c.Request.Context(), loginUser)
	if err != nil {
		slog.Error("User validation failed", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
		return
	}

	//generate claims for token
	var claims = auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-service",
			Subject:   user.ID,                                              // userId
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(50 * time.Minute)), // after 50 minutes, this token expires
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Roles: user.Roles,
	}

	// generate token
	tokenStr, err := h.a.GenerateToken(claims)
	if err != nil {
		slog.Error("Token Generation failed", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token Generation failed"})
		return
	}

	//token generation successful, return token
	slog.Info("Successful login, token generated")
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func (h *Handler) AuthCheck(c *gin.Context) {
	// we are using type assertion here for converting any type to auth.Claims
	// always use ok when using type assertion so that there is no panic
	claims, ok := c.Request.Context().Value(middleware.ClaimsKey).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid auth token"})
		return
	}
	userId := claims.Subject
	c.JSON(200, gin.H{"Authentication Check": "You are authenticated " + userId})
}
