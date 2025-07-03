package api

import (
	"shop-dashboard/internal/handlers"
	mw "shop-dashboard/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewHandler()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Health check
	r.Get("/health", h.HealthCheck)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.With(mw.AuthMiddleware).Route("/admin", func(r chi.Router) {
			r.Get("/products", h.GetProducts)
			r.Post("/create/product", h.CreateProduct)
			r.Post("/edit/product/{id}", h.EditProductData)
			r.Post("/create/productUniques", h.GenerateProductUniques)
			r.Delete("/delete/product/{id}", h.SoftDeleteProduct)
			r.Post("/republish/product/{id}", h.RepublishProduct)
			r.Get("/product/{id}", h.GetProduct)

			r.Route("/find", func(r chi.Router) {
				r.Get("/vendors", h.FindVendors)
				r.Get("/categories", h.FindCategories)
				r.Get("/tags", h.FindTags)
			})
		})

		r.Post("/create/vendorSeName", h.GenerateVendorSeName)
		r.Post("/register/vendor", h.RegisterVendorHandler)
	})

	return r
}
