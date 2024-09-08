package router

import (
	"github.com/labstack/echo/v4"
)

type API struct {
	Versions []*Version
}

type Version struct {
	Version   string
	Resources []*Resource
}

type Resource struct {
	Resource  string
	Endpoints []*Endpoint
}

type Endpoint struct {
	Method      string
	SuffixPath  string
	HandlerFunc echo.HandlerFunc
}

func GetAPIs() *API {
	return &API{
		Versions: []*Version{
			{
				Version: "v1",
				Resources: []*Resource{
					invoice(),
				},
			},
		},
	}
}
