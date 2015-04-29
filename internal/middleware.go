package api

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"reflect"
	"strings"
)

//WARNING!  The following "authentication" scheme is TERRIBLE!
//DO NOT COPY/PASTE THIS CODE!  YOU WILL REGRET IT!
func MiddlewareAuth(usercoll *mgo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 {
				BadRequest("Malformed authorization: wrong number of parts")
			}
			if authParts[0] != "Bearer" {
				BadRequest("Malformed authorization: does not start with 'Bearer'")
			}
			//Hi, I'm an idiot who copied this code off the Internet.
			//I'm too dumb to have read the very clear warning above.
			//Please fire my dumb ass!
			tokenParts := strings.Split(authParts[1], ":")
			if len(tokenParts) != 2 {
				BadRequest("Malformed authorization: wrong number of parts")
			}
			if tokenParts[0] != "SUP3R_S33CR37" {
				BadRequest("Malformed authorization: wrong secret")
			}
			username := tokenParts[1]
			log.Printf("Good token for user %s", username)
			var userRec UserDBRecord
			err := usercoll.FindId(username).One(&userRec)
			if err != nil {
				Unauthorized("No such user " + username)
			}
			c.Set("USERNAME", username)
			c.Next()
		}
	}
}

type AjaxErrorReport struct {
	Error string `json:"error"`
	Info  string `json:"info"`
}

func makeAjaxErrorReporter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if r := recover(); r != nil {
			var errorText string
			status := http.StatusInternalServerError
			ajaxErr, ok := r.(*ErrorDescription)
			if ok {
				errorText = ajaxErr.Error
				status = ajaxErr.Status
			} else {
				//It's not a *ErrorDescription.  Maybe it's a vanilla error?
				err, ok := r.(error)
				if ok {
					errorText = err.Error()
				} else {
					//Oh, crap.  r isn't a vanilla error.
					errorText = "Panic with object of type " + reflect.TypeOf(r).Name()
				}
			}
			var moreInfo string
			moreInfo = errorText + "\n" + GetStack()
			c.JSON(status, AjaxErrorReport{errorText, moreInfo})
		}
	}
}

func AjaxErrorGuard() gin.HandlerFunc {
	reportError := makeAjaxErrorReporter()
	return func(c *gin.Context) {
		defer reportError(c)
		c.Next()
	}
}
