package main

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
)

func PanicRecover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (passBackError error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10) // 4 KB
					length := runtime.Stack(stack, true)
					c.Logger().Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					passBackError = err
					c.Error(err)
				}
			}()

			passBackError = next(c)
			return passBackError
		}
	}
}

func main() {
	e := echo.New()

	// Register the PanicRecover middleware
	e.Use(PanicRecover())

	// Sample route that triggers a panic
	e.GET("/panic", func(c echo.Context) error {
		panic("something went wrong!")
	})

	// Start the server
	e.Start(":4000")
}
