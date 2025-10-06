package api

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
)

type Registry struct {
	echo       *echo.Echo
	operations []Operation
}

func NewRegistry(e *echo.Echo) *Registry {
	return &Registry{
		echo:       e,
		operations: make([]Operation, 0, 10),
	}
}

func (r *Registry) Register(handlerFunc echo.HandlerFunc, opts ...OperationOption) echo.HandlerFunc {
	pc := reflect.ValueOf(handlerFunc).Pointer()
	handlerName := runtime.FuncForPC(pc).Name()
	op := Operation{
		handlerName: handlerName,
	}

	// path and method will be set when the route is added to Echo and fetched later

	for _, opt := range opts {
		opt(&op)
	}

	r.operations = append(r.operations, op)

	return handlerFunc
}

func (r *Registry) Operations() []Operation {
	routes := r.echo.Routes()
	for i, op := range r.operations {
		for _, route := range routes {
			if route.Name == op.handlerName {
				r.operations[i].path = route.Path
				r.operations[i].method = route.Method
				r.operations[i].id = fmt.Sprintf("%s-%s", strings.ToLower(route.Method), strings.ReplaceAll(strings.Trim(route.Path, "/"), "/", "-"))
				break
			}
		}
	}
	return r.operations
}

type Operation struct {
	id          string
	method      string
	path        string
	handlerName string
	description string
	summary     string
	tags        []string
	requests    map[string]Request
	responses   map[int]Response
}

func (o Operation) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Operation ID: %s\n", o.id))
	b.WriteString(fmt.Sprintf("Method: %s\n", o.method))
	b.WriteString(fmt.Sprintf("Path: %s\n", o.path))
	b.WriteString(fmt.Sprintf("Handler: %s\n", o.handlerName))
	if o.summary != "" {
		b.WriteString(fmt.Sprintf("Summary: %s\n", o.summary))
	}
	if o.description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", o.description))
	}
	if len(o.tags) > 0 {
		b.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(o.tags, ", ")))
	}
	if len(o.requests) > 0 {
		b.WriteString("Requests:\n")
		for contentType, req := range o.requests {
			b.WriteString(fmt.Sprintf("  - Content-Type: %s, Body Type: %s\n", contentType, req.bodyType.String()))
		}
	}
	if len(o.responses) > 0 {
		b.WriteString("Responses:\n")
		for statusCode, resp := range o.responses {
			b.WriteString(fmt.Sprintf("  - Status Code: %d, Body Type: %s, Content-Type: %s, Description: %s\n",
				statusCode, resp.bodyType.String(), resp.contentType, resp.description))
		}
	}
	return b.String()
}

type OperationOption func(*Operation)

func WithID(id string) OperationOption {
	return func(o *Operation) {
		o.id = id
	}
}

func WithDescription(description string) OperationOption {
	return func(o *Operation) {
		o.description = description
	}
}

func WithSummary(summary string) OperationOption {
	return func(o *Operation) {
		o.summary = summary
	}
}

func WithTags(tags ...string) OperationOption {
	return func(o *Operation) {
		o.tags = tags
	}
}

func WithResponse[T any](statusCode int, description string, contentType string) OperationOption {
	return func(o *Operation) {
		var t T
		if o.responses == nil {
			o.responses = make(map[int]Response)
		}
		o.responses[statusCode] = Response{
			bodyType:    reflect.TypeOf(t),
			description: description,
			contentType: contentType,
		}
	}
}

func WithRequest[T any](contentType string) OperationOption {
	return func(o *Operation) {
		var t T
		if o.requests == nil {
			o.requests = make(map[string]Request)
		}
		o.requests[contentType] = Request{
			bodyType: reflect.TypeOf(t),
		}
	}
}

type Response struct {
	bodyType    reflect.Type
	description string
	contentType string
}

func (r Response) String() string {
	return fmt.Sprintf("%s (%s): %s", r.bodyType.String(), r.contentType, r.description)
}

type Request struct {
	bodyType reflect.Type
}

func (r Request) String() string {
	return r.bodyType.String()
}
