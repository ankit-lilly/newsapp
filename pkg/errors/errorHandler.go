package errors

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {

	log.Println(err)
	code := http.StatusInternalServerError
	message := http.StatusText(http.StatusInternalServerError)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	} else if IsApiError(err) {
		apiErr := err.(*ApiError)
		code = apiErr.Code
		message = apiErr.Message
	} else if IsNotFoundError(err) {
		code = http.StatusNotFound
		message = err.Error()
		c.Redirect(http.StatusFound, "/404")
	} else {
		c.Logger().Error(err)
	}

	c.Logger().Error(err)

	errorPage := fmt.Sprintf("pkg/views/%d.html", code)
	if err := c.File(errorPage); err != nil {
		// If the error page doesn't exist, send a JSON response
		if !c.Response().Committed {
			c.JSON(code, echo.Map{"error": message})
		}
	}
}
