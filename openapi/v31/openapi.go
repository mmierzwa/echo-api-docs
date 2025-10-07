package v31

const (
	openAPIVersion = "3.1.0"
)

type dataType string

const (
	typeString  dataType = "string"
	typeNumber  dataType = "number"
	typeInteger dataType = "integer"
	typeBoolean dataType = "boolean"
	typeArray   dataType = "array"
	typeObject  dataType = "object"
)

type format string

const (
	formatInt32    format = "int32"
	formatInt64    format = "int64"
	formatFloat    format = "float"
	formatDouble   format = "double"
	formatByte     format = "byte"
	formatBinary   format = "binary"
	formatDate     format = "date"
	formatDateTime format = "date-time"
	formatPassword format = "password"
	formatEmail    format = "email"
	formatUUID     format = "uuid"
)

type document struct {
	OpenAPI    string              `json:"openapi" yaml:"openapi"`
	Info       info                `json:"info" yaml:"info"`
	Servers    []server            `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths      map[string]pathItem `json:"paths" yaml:"paths"`
	Components *components         `json:"components,omitempty" yaml:"components,omitempty"`
	Tags       []string            `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func newDocument(title, version string) *document {
	return &document{
		OpenAPI: openAPIVersion,
		Info:    newInfo(title, version),
		Paths:   make(map[string]pathItem),
	}
}

func (d *document) addServer(url, description string) {
	d.Servers = append(d.Servers, server{
		URL:         url,
		Description: description,
	})
}

func (d *document) addPath(path string, item pathItem) {
	d.Paths[path] = item
}

func (d *document) setComponents(components *components) {
	d.Components = components
}

func (d *document) addTag(tag string) {
	for _, t := range d.Tags {
		if t == tag {
			return
		}
	}
	d.Tags = append(d.Tags, tag)
}

type info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *license `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version,omitempty" yaml:"version,omitempty"`
}

func newInfo(title, version string) info {
	return info{
		Title:   title,
		Version: version,
	}
}

func (d *info) setDescription(description string) {
	d.Description = description
}

func (d *info) setTermsOfService(terms string) {
	d.TermsOfService = terms
}

func (d *info) setContact(name, url, email string) {
	d.Contact = &contact{
		Name:  name,
		URL:   url,
		Email: email,
	}
}

func (d *info) setLicense(name, identifier, url string) {
	d.License = &license{
		Name:       name,
		Identifier: identifier,
		URL:        url,
	}
}

type contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type license struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
}

type server struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type pathItem struct {
	Summary     string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Servers     []server    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Get         *operation  `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *operation  `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *operation  `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *operation  `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *operation  `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *operation  `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *operation  `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *operation  `json:"trace,omitempty" yaml:"trace,omitempty"`
}

type operation struct {
	Tags        []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []parameter           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *requestBody          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   responses             `json:"responses" yaml:"responses"`
	Deprecated  bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security    []map[string][]string `json:"security,omitempty" yaml:"security,omitempty"`
	Servers     []server              `json:"servers,omitempty" yaml:"servers,omitempty"`
}

type requestBody struct {
	Description string           `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]media `json:"content" yaml:"content"`
	Required    bool             `json:"required,omitempty" yaml:"required,omitempty"`
}

type media struct {
	Schema  *schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example any     `json:"example,omitempty" yaml:"example,omitempty"`
}

type responses map[string]response

type response struct {
	Description string            `json:"description" yaml:"description"`
	Content     map[string]media  `json:"content,omitempty" yaml:"content,omitempty"`
	Headers     map[string]header `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type header struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type schema struct {
	Type                 dataType          `json:"type,omitempty" yaml:"type,omitempty"`
	Properties           map[string]schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items                *schema           `json:"items,omitempty" yaml:"items,omitempty"`
	Format               format            `json:"format,omitempty" yaml:"format,omitempty"`
	Ref                  string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Enum                 []string          `json:"enum,omitempty" yaml:"enum,omitempty"`
	Required             []string          `json:"required,omitempty" yaml:"required,omitempty"`
	Nullable             bool              `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	OneOf                []schema          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []schema          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOf                []schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Example              any               `json:"example,omitempty" yaml:"example,omitempty"`
	Examples             []any             `json:"examples,omitempty" yaml:"example,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	AdditionalProperties *schema           `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
}

type parameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type components struct {
	Schemas         map[string]schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Parameters      map[string]parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	SecuritySchemes map[string]securityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Responses       map[string]response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Headers         map[string]header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	RequestBodies   map[string]requestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Examples        map[string]any            `json:"examples,omitempty" yaml:"examples,omitempty"`
	Links           map[string]link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

type securityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
	Flows            *oauthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
}

type oauthFlows struct {
	Implicit          *oauthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *oauthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *oauthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *oauthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

type oauthFlow struct {
	AuthorizationURL string   `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string   `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string   `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           []string `json:"scopes" yaml:"scopes"`
}

type link struct {
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody any               `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Server      *server           `json:"server,omitempty" yaml:"server,omitempty"`
}

type callback map[string]pathItem
