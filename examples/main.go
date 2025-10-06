package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mmierzwa/echo-api-docs/api"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	registry := api.NewRegistry(e)

	// test routes
	e.GET("/", registry.Register(echo.HandlerFunc(NewHelloFromRootHandler()),
		api.WithID("helloFromRoot"),
		api.WithDescription("A simple hello world endpoint at the root path"),
		api.WithTags("example", "hello"),
		api.WithResponse[helloFromRootResponse](http.StatusOK, "application/json", "Successful response"),
	))
	e.POST("/the-great-post", registry.Register(echo.HandlerFunc(NewHelloFromTheGreatPostHandler()),
		api.WithDescription("A great POST endpoint that greets you"),
		api.WithTags("example"),
		api.WithRequest[helloFromTheGreatPostRequest]("application/json"),
		api.WithResponse[helloFromTheGreatPostResponse](http.StatusOK, "Successful response", "application/json"),
		api.WithResponse[echo.HTTPError](http.StatusBadRequest, "Invalid request body", "application/json"),
	))

	fmt.Println("Registered API Operations:")
	for _, op := range registry.Operations() {
		fmt.Println(op)
	}

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

type (
	HelloFromRootHandler         echo.HandlerFunc
	HelloFromTheGreatPostHandler echo.HandlerFunc
)

type helloFromRootResponse struct {
	Message string `json:"message"`
}

func NewHelloFromRootHandler() HelloFromRootHandler {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, helloFromRootResponse{
			Message: "Hello from the root endpoint! You just passed",
		})
	}
}

type helloFromTheGreatPostRequest struct {
	Name        string    `json:"name"`
	Foodie      bool      `json:"foodie"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type helloFromTheGreatPostResponse struct {
	Message     string    `json:"message"`
	Name        string    `json:"name"`
	Foodie      bool      `json:"foodie"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func NewHelloFromTheGreatPostHandler() HelloFromTheGreatPostHandler {
	return func(c echo.Context) error {
		req := helloFromTheGreatPostRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").WithInternal(err)
		}

		return c.JSON(http.StatusOK, helloFromTheGreatPostResponse{
			Message: "Hello from the great post endpoint!",
		})
	}
}
