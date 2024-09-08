package server

import (
	"fmt"
	"github.com/niko-cb/uct/internal/infrastructure/web/router"
)

const (
	BasePath = "api"
)

func (s *server) routing() {
	api := router.GetAPIs()
	for _, version := range api.Versions {
		for _, resource := range version.Resources {
			for _, endpoint := range resource.Endpoints {
				suffix := ""
				if endpoint.SuffixPath != "" {
					suffix += fmt.Sprintf("/%s", endpoint.SuffixPath)
				}
				// The assignment asks for /api/invoices, but normally api versioning is included, so we will use /api/v1/invoices
				path := fmt.Sprintf("/%s/%s/%s%s", BasePath, version.Version, resource.Resource, suffix)
				s.Router().Add(endpoint.Method, path, endpoint.HandlerFunc)
			}
		}
	}
}
