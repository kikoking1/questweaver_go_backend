package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Location  string  `json:"location"`
	Temp      float64 `json:"temp"`
	Condition string  `json:"condition"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
		log.Println("Falling back to system environment variables")
	}

	// Get Auth0 configuration from environment variables
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")

	if auth0Domain == "" || auth0Audience == "" {
		log.Fatal("AUTH0_DOMAIN and AUTH0_AUDIENCE environment variables must be set")
	}

	issuerURL := fmt.Sprintf("https://%s/", auth0Domain)

	// Parse issuer URL for JWKS provider
	issuerURLParsed, err := url.Parse(issuerURL)
	if err != nil {
		log.Fatalf("Failed to parse issuer URL: %v", err)
	}

	// Set up JWKS provider to fetch Auth0's public keys
	provider := jwks.NewCachingProvider(issuerURLParsed, 1*time.Hour)

	// Set up the token validator
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL,
		[]string{auth0Audience},
		validator.WithAllowedClockSkew(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to set up JWT validator: %v", err)
	}

	// Create JWT middleware with error handler
	ensureValidToken := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(jwtErrorHandler),
	)

	// Set up Gin router
	r := gin.Default()

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Weather endpoint (requires Auth0 token)
	// Convert Auth0 middleware to Gin middleware
	r.GET("/weather", auth0GinMiddleware(ensureValidToken), getWeather)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// auth0GinMiddleware adapts Auth0 JWT middleware to Gin
func auth0GinMiddleware(m *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nextCalled bool
		handler := m.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		// If next was not called, the middleware already wrote a response
		if !nextCalled {
			c.Abort()
		}
	}
}

// getWeather handles the weather endpoint
func getWeather(c *gin.Context) {
	weather := WeatherResponse{
		Location:  "San Francisco, CA",
		Temp:      72.5,
		Condition: "Sunny",
	}

	c.JSON(http.StatusOK, weather)
}

// jwtErrorHandler handles JWT validation errors
func jwtErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("JWT validation error: %v", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"Unauthorized"}`))
}
