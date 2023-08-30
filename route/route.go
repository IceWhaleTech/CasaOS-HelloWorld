package route

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"

	"github.com/IceWhaleTech/CasaOS-HelloWorld/codegen"
	"github.com/IceWhaleTech/CasaOS-HelloWorld/service"
)

type HelloWorld struct{}

var (
	_swagger *openapi3.T

	APIPath string
)

func init() {
	swagger, err := codegen.GetSwagger()
	if err != nil {
		panic(err)
	}

	_swagger = swagger

	u, err := url.Parse(_swagger.Servers[0].URL)
	if err != nil {
		panic(err)
	}

	APIPath = strings.TrimRight(u.Path, "/")
}

func GetRouter() http.Handler {
	hello := NewHelloWorld()

	e := echo.New()

	e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.POST, echo.GET, echo.OPTIONS, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderXCSRFToken, echo.HeaderContentType, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderAccessControlAllowMethods, echo.HeaderConnection, echo.HeaderOrigin, echo.HeaderXRequestedWith},
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders},
		MaxAge:           172800,
		AllowCredentials: true,
	}))

	e.Use(echo_middleware.Gzip())

	e.Use(echo_middleware.Logger())

	e.Use(middleware.OapiRequestValidatorWithOptions(_swagger, &middleware.Options{
		Options: openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc},
	}))

	codegen.RegisterHandlersWithBaseURL(e, hello, APIPath)

	return e
}

func NewHelloWorld() codegen.ServerInterface {
	return &HelloWorld{}
}

func (h *HelloWorld) Ping(ctx echo.Context) error {
	ping := service.HelloWorld.Ping()

	return ctx.JSON(http.StatusOK, codegen.PongOK{
		Data: &ping,
	})
}
