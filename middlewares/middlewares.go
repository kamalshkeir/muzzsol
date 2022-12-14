package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kamalshkeir/muzzsol/services"
	"github.com/labstack/echo/v4"
)

var CookieSessionName = "session"

type ContextKey string
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := ""
		// get session cookie
		cook, err := c.Cookie(CookieSessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized,map[string]any{
				"error":"Unauthorised access",
			})
			return nil
		}
		if token == "" {
			token = cook.Value
		}
		dec,_ := services.Decrypt(token)
		m := map[string]any{}
		err = json.Unmarshal([]byte(dec),&m)
		if err != nil {
			c.JSON(http.StatusUnauthorized,map[string]any{
				"error":"Unauthorised access",
			})
			return nil
		}
		if _,ok := m["id"];!ok {
			c.JSON(http.StatusUnauthorized,map[string]any{
				"error":"Unauthorised access",
			})
			return nil
		} else {
			// authorized, use context to pass session to handlers
			ctx := c.Request().Context()
			req := c.Request().WithContext(context.WithValue(ctx,ContextKey(CookieSessionName),m))
			*c.Request()=*req
			return next(c)
		}
		
	}
}

func DeleteCookie(c echo.Context, name string) {
	cook, err := c.Cookie(name)
	if err != nil {
		return
	}
	cook.Expires = time.Unix(0, 0)
	cook.MaxAge = -1
	cook.Value = ""
	c.SetCookie(cook)
}