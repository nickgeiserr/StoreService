package handlers

import (
	"ClubmineStoreService/logger"
	"ClubmineStoreService/models"
	"ClubmineStoreService/models/RequestTypes"
	"ClubmineStoreService/models/Response"
	"ClubmineStoreService/services"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

type (
	ListingsHandler interface {
		CreateListing(c echo.Context) error
		GetListing(c echo.Context) error
		DeleteListing(c echo.Context) error
	}

	listingsHandler struct {
		services.ListingsService
	}
)

func (l listingsHandler) CreateListing(c echo.Context) error {
	listingRequest := new(RequestTypes.ListingConstructionRequest)
	if err := c.Bind(listingRequest); err != nil {
		return Response.InternalServerError(c, err)
	}

	if listingRequest.ItemID == "" || listingRequest.ListerID == "" || listingRequest.ItemStack == "" {
		return Response.BadData(c, nil)
	}
	listing := models.Listing{
		ListingID: uuid.New().String(),
		ListerID:  listingRequest.ListerID,
		ItemID:    listingRequest.ItemID,
		ItemStack: listingRequest.ItemStack,
		Price:     listingRequest.Price,
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05-07:00"),
	}

	r, err := l.ListingsService.CreateListing(listing)
	if r == false || err != nil {
		if errors.Is(err, models.ForeignKeyConstraint{}) || errors.Is(err, models.ItemDoesNotExist{}) {
			return Response.BadData(c, nil)
		}
		return Response.InternalServerError(c, err)
	}
	responseListing, err := l.ListingsService.GetListing(listing.ListingID)
	if err != nil {
		return Response.InternalServerErrorWithTypeResponse(c, []models.ListingResponse{}, err)
	}
	return Response.RequestCompleteWithTypeResponse(c, responseListing)
}

func (l listingsHandler) GetListing(c echo.Context) error {
	fetchRequest := new(RequestTypes.GetListingRequest)

	fetchRequest.ListingID = c.QueryParam("listing_id")
	page := c.QueryParam("page")
	fetchRequest.ListerID = c.QueryParam("lister_id")

	var pageInt int
	if page != "" {
		num, err := strconv.Atoi(page)
		if err != nil {
			return Response.BadData(c, err)
		}
		pageInt = num
	} else {
		pageInt = 0
	}

	var listings []models.ListingResponse
	var err error

	switch {
	case fetchRequest.ListerID != "":
		listings, err = l.ListingsService.GetAllWithFilter(fetchRequest.ListerID, pageInt, "lister")
	case fetchRequest.ListingID != "":
		listing, err := l.ListingsService.GetListing(fetchRequest.ListingID)
		if err != nil {
			return Response.InternalServerErrorWithTypeResponse(c, []models.ListingResponse{}, err)
		}
		if (listing != models.ListingResponse{}) {
			return Response.RequestCompleteWithTypeResponse(c, listing)
		}
		return Response.NotFound(c)
	default:
		listings, err = l.ListingsService.GetAllListings(pageInt)
	}

	if err != nil {
		return Response.InternalServerErrorWithTypeResponse(c, []models.ListingResponse{}, err)
	}

	if listings == nil || len(listings) == 0 {
		return Response.RequestCompleteWithTypeResponse(c, []models.ListingResponse{})
	}

	return Response.RequestCompleteWithTypeResponse(c, listings)
}

func (l listingsHandler) DeleteListing(c echo.Context) error {
	deleteListingRequest := new(RequestTypes.DeleteListingRequest)

	deleteListingRequest.ListingID = c.QueryParam("listing_id")
	deleteListingRequest.ID = c.Request().Header.Get("PlayerUID")

	if deleteListingRequest.ListingID == "" || deleteListingRequest.ID == "" {
		return Response.BadData(c, nil)
	}

	listing, err := l.ListingsService.GetListing(deleteListingRequest.ListingID)
	logger.Debug(listing.CreatedAt)
	if err != nil {
		return Response.InternalServerError(c, err)
	}
	if (listing != models.ListingResponse{}) {
		if listing.ListerID != deleteListingRequest.ID {
			return Response.MethodNotAllowed(c)
		}
		deleteListingRequest = nil
		err := l.ListingsService.DeleteListing(listing.ListingID)
		if err != nil {
			return Response.InternalServerError(c, err)
		}
		return Response.RequestComplete(c)
	}

	return Response.NotFound(c)
}
