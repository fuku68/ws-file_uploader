package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"

	"echo-upload/config"
	"echo-upload/uploader"
)

func (h *Handler) wsupload(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		u := uploader.NewS3Uploader(config.S3_BUCKET)
		u.UploadStart("test.txt")

		c.Logger().Info("start")
		for {
			// Write
			err := websocket.Message.Send(ws, "uplaoding...")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if errors.Is(err, io.EOF) {
				u.UploadEnd()
			} else if err != nil {
				c.Logger().Error(err)
			} else if len(msg) == 0 {
				u.UploadEnd()
				c.Logger().Info("end")
				break
			} else {
				err = u.Upload([]byte(msg))
				if err != nil {
					c.Logger().Error(err)
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
