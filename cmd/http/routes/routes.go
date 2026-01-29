package routes

import (
	"microservice-products-catalog/cmd/http/dependencies"
	"net/http"
)

// TODO [technical debate] handle different versions

func SetupProductRoutes(mux *http.ServeMux, dep dependencies.Dependencies) {
	mux.HandleFunc("/api/products", EnableProductsCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dep.ReaderHandler.HandleGetProducts(w, r)

		case http.MethodPost:
			dep.WriterHandler.HandleCreateProduct(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/products/", EnableProductsCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dep.ReaderHandler.HandleGetProductByID(w, r)

		case http.MethodDelete:
			dep.WriterHandler.HandleDeleteProduct(w, r)

		case http.MethodPut:
			dep.WriterHandler.HandleUpdateProduct(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func SetupOrderRoutes(mux *http.ServeMux, dep dependencies.Dependencies) {
	mux.HandleFunc("/api/orders", EnableProductsCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dep.ReaderHandler.HandleGetOrders(w, r)

		case http.MethodPost:
			dep.WriterHandler.HandleCreateOrder(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}
