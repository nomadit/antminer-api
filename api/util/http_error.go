package util

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

// CheckErrorInHTTP is util for http
func CheckErrorInHTTP(c *gin.Context, err error) error {
	//if raven.ProjectID() == "" {
	//	raven.SetDSN("")
	//}

	if err != nil {
		raven.CaptureError(err, nil)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	return err
}
