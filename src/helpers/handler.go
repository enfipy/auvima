package helpers

import (
	"net/http"

	"github.com/labstack/echo"
)

type Response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func Handle(logic func(ctx echo.Context) interface{}) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		res := new(Response)

		defer func() {
			if rec := recover(); rec != nil {
				err := rec.(error)
				msg := err.Error()

				res.Status = false
				res.Error = msg
				ctx.JSON(http.StatusInternalServerError, res)
			} else {
				res.Status = true
				ctx.JSON(http.StatusOK, res)
			}
		}()

		res.Data = logic(ctx)
		return nil
	}
}
