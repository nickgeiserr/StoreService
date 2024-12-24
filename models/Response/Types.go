package Response

import (
	"ClubmineStoreService/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type InternalServerErrorWithType struct {
	Message string
	Data    interface{} // Additional data
}

type RequestCompleteWithType struct {
	Message string
	Data    interface{}
}

func (e *InternalServerErrorWithType) Error() string {
	return e.Message
}

func InternalServerError(c echo.Context, err error) error {
	logger.Error("Internal Server Error : ", zap.Error(err))
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal Server Error."})
}

func InternalServerErrorWithTypeResponse(c echo.Context, data interface{}, error error) error {
	logger.Error("Internal Server Error : ", zap.Error(error))
	err := &InternalServerErrorWithType{
		Message: "Internal Server Error.",
		Data:    data,
	}
	return c.JSON(http.StatusInternalServerError, err)
}

func RequestCompleteWithTypeResponse(c echo.Context, data interface{}) error {
	payload := &RequestCompleteWithType{
		Message: "Success.",
		Data:    data,
	}
	return c.JSON(http.StatusOK, payload)
}

func RequestComplete(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "complete"})
}

func BadData(c echo.Context, err error) error {
	if err != nil {
		logger.Error("Bad Data Error : ", zap.Error(err))
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Bad Data."})
}

func NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not found."})
}

func PaymentRequired(c echo.Context) error {
	return c.JSON(http.StatusPaymentRequired, map[string]interface{}{"message": "Payment Required."})
}

func UnauthorizedRequest(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Unauthorized."})
}

func MethodNotAllowed(c echo.Context) error {
	return c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{"message": "not allowed."})
}

func ForbiddenRequest(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "forbidden"})
}

func ConflictResponse(c echo.Context) error {
	return c.JSON(http.StatusConflict, map[string]interface{}{"message": "conflict"})
}
