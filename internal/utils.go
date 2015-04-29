package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

func forceAuth(c *gin.Context) string {
	usernameI, err := c.Get("USERNAME")
	if err != nil {
		Unauthorized("Not authorized")
	} else {
		username, ok := usernameI.(string)
		if !ok {
			InternalError("USERNAME is of incorrect type!")
		}
		return username
	}
	panic(errors.New("This should not be possible!"))
}

func extractJSON(req *http.Request, payload interface{}) {
	defer req.Body.Close()
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		InternalError(err.Error())
	}
	log.Printf("Body: %s", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, payload)
	if err != nil {
		BadRequest(err.Error())
	}
}

func GetStack() string {
	chunk := 2048
	stackTrace := ""
	for {
		buf := make([]byte, chunk)
		buflen := runtime.Stack(buf, false)
		if buflen < chunk {
			stackTrace = string(buf[:buflen])
			break
		}
		chunk = chunk * 2
	}

	return stackTrace
}
