package handlers

import (
	"ClubmineStoreService/models"
	"ClubmineStoreService/models/Response"
	"ClubmineStoreService/services"
	"errors"
	"github.com/labstack/echo/v4"
)

type (
	PurchaseHandler interface {
		Purchase(c echo.Context) error
	}

	purchaseHandler struct {
		services.PurchaseService
		services.ListingsService
	}
)

func (p purchaseHandler) Purchase(c echo.Context) error {
	transactionAttempt := new(models.TransactionAttempt)
	if err := c.Bind(transactionAttempt); err != nil {
		return Response.InternalServerError(c, err)
	}
	listing, err := p.ListingsService.GetListing(transactionAttempt.ListingID)
	if err != nil {
		return Response.InternalServerError(c, err)
	}
	if (listing != models.ListingResponse{}) {
		complete, err := p.PurchaseService.Purchase(listing, transactionAttempt)
		if complete == false && err != nil {
			if errors.Is(err, models.UserDoesNotExist{}) {
				return Response.MethodNotAllowed(c)
			}
			if errors.Is(err, models.UserBidOnOwnBid{}) {
				return Response.ForbiddenRequest(c)
			}
			if errors.Is(err, models.SmallBidError{}) {
				return Response.BadData(c, nil)
			}
			return Response.InternalServerError(c, err)
		}
		if complete == false && err == nil {
			return Response.PaymentRequired(c)
		}
		return Response.RequestComplete(c)
	}

	return Response.NotFound(c)
}
