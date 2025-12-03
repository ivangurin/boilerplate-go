package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"boilerplate/internal/model"
	"boilerplate/internal/pkg/metadata"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	// Записываем в буфер, но НЕ в ResponseWriter сразу
	return w.body.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	// Записываем в буфер, но НЕ в ResponseWriter сразу
	return w.body.WriteString(s)
}

func (m *middleware) Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bw := &bodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}

		ctx.Writer = bw

		ctx.Next()

		// Теперь обрабатываем ответ
		if ctx.Writer.Status() > http.StatusBadRequest {
			resp := &model.HandlerError{}
			err := json.Unmarshal(bw.body.Bytes(), resp)
			if err != nil {
				resp.Error = bw.body.String()
			}

			errorText := resp.Error
			if errorText == "" {
				errorText = resp.Error
			}

			userID, _ := metadata.GetUserID(ctx)

			m.logger.ErrorKV(ctx, ctx.Request.URL.Path,
				"status", fmt.Sprint(ctx.Writer.Status()),
				"method", ctx.Request.Method,
				"path", ctx.Request.URL.Path,
				"user_id", fmt.Sprint(userID),
				"error", errorText)

			if ctx.Writer.Status() >= http.StatusInternalServerError {
				resp = &model.HandlerError{
					Error: "Внутренняя ошибка",
				}

				newBody, err := json.Marshal(resp)
				if err != nil {
					m.logger.ErrorKV(ctx, "Ошибка сериализации ответа",
						"error", err.Error())
					return
				}

				_, err = bw.ResponseWriter.Write(newBody)
				if err != nil {
					m.logger.ErrorKV(ctx, "Ошибка записи ответа",
						"error", err.Error())
				}
			} else {
				_, err = bw.ResponseWriter.Write(bw.body.Bytes())
				if err != nil {
					m.logger.ErrorKV(ctx, "Ошибка записи ответа",
						"error", err.Error())
				}
			}
		} else {
			_, err := bw.ResponseWriter.Write(bw.body.Bytes())
			if err != nil {
				m.logger.ErrorKV(ctx, "Ошибка записи ответа",
					"error", err.Error())
			}
		}
	}
}
