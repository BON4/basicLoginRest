package users

import "github.com/labstack/echo/v4"

type Handlers interface {
	Find() 	  echo.HandlerFunc
}
