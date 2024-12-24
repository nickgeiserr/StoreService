package handlers

import (
	"ClubmineStoreService/models"
	"ClubmineStoreService/models/RequestTypes"
	"ClubmineStoreService/models/Response"
	"ClubmineStoreService/services"
	"errors"
	"github.com/labstack/echo/v4"
)

type (
	ItemHandler interface {
		GetItem(c echo.Context) error
		CreateItem(c echo.Context) error
		DeleteItem(c echo.Context) error
	}

	itemHandler struct {
		services.ItemService
	}
)

func (i itemHandler) GetItem(c echo.Context) error {
	fetchRequest := new(RequestTypes.GetItemRequest)
	fetchRequest.ItemID = c.QueryParam("item_id")

	if fetchRequest.ItemID == "" {
		return Response.BadData(c, nil)
	}

	r, err := i.ItemService.GetItemFromId(fetchRequest.ItemID)
	if err != nil {
		if errors.Is(err, models.NoRowsError{}) {
			return Response.NotFound(c)
		}
		return Response.InternalServerError(c, err)
	}
	return Response.RequestCompleteWithTypeResponse(c, r)

}

func (i itemHandler) CreateItem(c echo.Context) error {
	createRequest := new(RequestTypes.CreateItemRequest)
	if err := c.Bind(createRequest); err != nil {
		return Response.InternalServerError(c, err)
	}

	if createRequest.ItemID == "" || createRequest.ItemStack == "" || createRequest.CreatorID == "" {
		return Response.BadData(c, nil)
	}

	err := i.ItemService.CreateItem(createRequest)
	if err != nil {
		if errors.Is(err, models.ForeignKeyConstraint{}) {
			return Response.BadData(c, nil)
		}
		if errors.Is(err, models.AlreadyExistsInDatabaseError{}) {
			return Response.ConflictResponse(c)
		}
		return Response.InternalServerError(c, err)
	}

	return Response.RequestCompleteWithTypeResponse(c, models.Item{ItemID: createRequest.ItemID, ItemStack: createRequest.ItemStack})
}

func (i itemHandler) DeleteItem(c echo.Context) error {
	deleteItemRequest := new(RequestTypes.DeleteItemRequest)

	deleteItemRequest.ItemID = c.QueryParam("item_id")
	if deleteItemRequest.ItemID == "" {
		return Response.BadData(c, nil)
	}

	err := i.ItemService.DeleteItem(deleteItemRequest.ItemID)
	if err != nil {
		return Response.InternalServerError(c, err)
	}

	return Response.RequestComplete(c)
}
