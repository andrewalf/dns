package main

import (
	_ "dns/api/openapi"
	"dns/internal/location/handler"
	"dns/internal/location/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var sectorID int
var port int

// port and sectorID are passed as env variable for docker container
func init() {
	rawSectorID := os.Getenv("SECTOR_ID")
	if rawSectorID == "" {
		log.Fatal("sectorID must be specified")
	}
	var e error
	sectorID, e = strconv.Atoi(rawSectorID)
	if e != nil {
		log.Fatal("sectorID must be integer")
	}
	rawPort := os.Getenv("DNS_PORT")
	if rawPort == "" {
		log.Fatal("port must be specified")
	}
	port, e = strconv.Atoi(rawPort)
	if e != nil {
		log.Fatal("port must be integer")
	}
}

// @title Drone Navigation System
// @version 1.0
// @BasePath /api/v1
func main() {
	r := getRouter()
	fmt.Printf("DNS server running at port %v\n", port)
	// all requests are handled in separate goroutines so it;s ok use log.Fatal
	// it will be executed when server will be stopped (server is in a loop)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func getRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// calculator strategy is interface and any implementation can be injected
			// also we can create a pool of calculators, for example: inject CalcOne strategy
			// for api/v1/ handler, CalcTwo for /api/v2/ handler
			h := handler.NewGetLocationHandler(service.NewDefaultCalculator(sectorID))
			r.Method(http.MethodGet, "/location", h)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
