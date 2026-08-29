package handlers

import (
	"net/http"

	"github.com/ankit-lilly/newsmaxxing/internal/templates"
	"github.com/ankit-lilly/newsmaxxing/internal/templates/components/ui"
	"github.com/labstack/echo/v5"
)

type ErrorHandler struct {
	*BaseHandler
}

func (h *ErrorHandler) CustomHTTPErrorHandler(c *echo.Context, err error) {
	c.Logger().Error("request failed", "url", c.Request().URL.String(), "error", err)

	response, unwrapErr := echo.UnwrapResponse(c.Response())
	if unwrapErr == nil && response.Committed {
		return
	}
	statusCode := echo.StatusCode(err)
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	if renderErr := h.Render(c, RenderProps{
		Title:            "Error",
		Component:        ui.ErrorBlock(err.Error()),
		WrapperComponent: templates.Index,
		CacheStrategy:    "no-cache",
		StatusCode:       statusCode,
	}); renderErr != nil {
		c.Logger().Error("failed to render error response", "error", renderErr)
	}
}
