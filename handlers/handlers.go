package handlers

import (
	"ClubmineStoreService/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handlers struct {
	ListingsHandler
	PurchaseHandler
	ItemHandler
}

func New(s *services.Services) *Handlers {
	return &Handlers{
		ListingsHandler: &listingsHandler{s.Listing},
		PurchaseHandler: &purchaseHandler{s.Purchase, s.Listing},
		ItemHandler:     &itemHandler{s.Item},
	}
}

func SetAPI(e *echo.Echo, h *Handlers) {
	g := e.Group("")

	// Listings
	g.GET("/listings", h.ListingsHandler.GetListing)
	g.POST("/listings", h.ListingsHandler.CreateListing)
	g.DELETE("/listings", h.ListingsHandler.DeleteListing)

	// Purchase
	g.POST("/listings/buy", h.PurchaseHandler.Purchase)

	// Items
	g.GET("/items", h.ItemHandler.GetItem)
	g.POST("/items", h.ItemHandler.CreateItem)
	g.DELETE("/items", h.ItemHandler.DeleteItem)

}

func Echo() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	return e
}
