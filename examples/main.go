package main

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	handlersRegistry = make([]string, 0, 2)
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// test routes
	e.GET("/", registerHandler(NewHelloFromRootHandler()))
	e.POST("/the-great-post", registerHandler(NewHelloFromTheGreatPostHandler()))

	e.GET("/routes", getRoutesHandler)

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

type routeInfoResponse struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	HandlerName string `json:"handler_name"`
}

type routesResponse struct {
	Routes   []routeInfoResponse `json:"routes"`
	Handlers []string            `json:"handlers"`
}

func getRoutesHandler(c echo.Context) error {
	resp := routesResponse{
		Routes:   make([]routeInfoResponse, 0, len(c.Echo().Routes())),
		Handlers: handlersRegistry,
	}
	for _, r := range c.Echo().Routes() {
		resp.Routes = append(resp.Routes, routeInfoResponse{
			Method:      r.Method,
			Path:        r.Path,
			HandlerName: r.Name,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func registerHandler(handler func(echo.Context) error) echo.HandlerFunc {
	name := funcName(handler)
	handlersRegistry = append(handlersRegistry, name)
	return handler
}

func funcName(fn any) string {
	ptr := reflect.ValueOf(fn).Pointer()
	return runtime.FuncForPC(ptr).Name()
}
