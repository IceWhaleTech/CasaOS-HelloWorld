//go:generate bash -c "mkdir -p codegen && go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.4 -generate types,server,spec -package codegen api/openapi.yaml > codegen/hello_world_api.go"

package main

import (
	"net"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/model"

	"github.com/IceWhaleTech/CasaOS-HelloWorld/config"
	"github.com/IceWhaleTech/CasaOS-HelloWorld/route"
	"github.com/IceWhaleTech/CasaOS-HelloWorld/service"
)

func main() {
	service.Initialize(config.CommonInfo.RuntimePath)

	// setup listener
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		panic(err)
	}

	// initialize routers and register at gateway
	if err := service.Gateway.CreateRoute(&model.Route{
		Path:   route.APIPath,
		Target: "http://" + listener.Addr().String(),
	}); err != nil {
		panic(err)
	}

	s := &http.Server{
		Handler:           route.GetRouter(),
		ReadHeaderTimeout: 5 * time.Second, // fix G112: Potential slowloris attack (see https://github.com/securego/gosec)
	}

	err = s.Serve(listener) // not using http.serve() to fix G114: Use of net/http serve function that has no support for setting timeouts (see https://github.com/securego/gosec)
	if err != nil {
		panic(err)
	}
}
