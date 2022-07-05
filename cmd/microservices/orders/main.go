package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting the orders microservices")

	ctx := cmd.Context()

	r, closefn := createOrderMicroservice()
	defer closefn()

	server := &http.Server{Addr: os.Getenv("SHOP_ORDER_BIND_ADDR"), Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()

	log.Println("closing order microservice")

	if err := server.Close(); err != nil {
		panic(err)
	}
}

func createOrderMicroservice() (router *chi.Mux, closefn func()) {
	cmd.WaitForService(os.Getenv("SHOP_RABBITMQ_ADDR"))

	shopHTTPClient := orders_infra_product.NewHttpClient(os.Getenv("SHOP_PRODUCTS_SERVICE_ADDR"))

	r := cmd.CreateRouter()

	orders_public_http.AddRoutes(r, ordersService, ordersRepo)
	orders_privete_http.AddRoutes(r, ordersService, ordersRepo)

	return r, func() {}
}
