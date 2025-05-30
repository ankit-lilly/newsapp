package handlers

import (
	"github.com/ankit-lilly/newsapp/internal/templates"
	"github.com/ankit-lilly/newsapp/internal/templates/components/ui"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	*BaseHandler
}

func (h *ErrorHandler) CustomHTTPErrorHandler(err error, c echo.Context) {
	c.Echo().Logger.Error("Something went wrong: ", c.Request().URL, err)

	if !c.Response().Committed {
		h.Render(c, RenderProps{
			Title:            "Error",
			Component:        ui.ErrorBlock(err.Error()),
			WrapperComponent: templates.Index,
			CacheStrategy:    "no-cache",
		})

	}

}
