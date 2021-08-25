package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	Register() echo.HandlerFunc
	Login()    echo.HandlerFunc
	Logout()   echo.HandlerFunc
	Update()   echo.HandlerFunc
	Delete()   echo.HandlerFunc
	GetByID()  echo.HandlerFunc
}

