package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"saaa-api/internal/core/config"
	"saaa-api/internal/core/context"
	"saaa-api/internal/core/utils"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Request request from client handle log
func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") && r.Method != http.MethodGet {
			buffer, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.Errorf("[Request] read body request error: %s", err)
				render.Status(r, config.RR.Internal.BadRequest.HTTPStatusCode())
				render.JSON(w, r, config.RR.Internal.BadRequest.WithLocale(r))
				return
			}

			rc := ioutil.NopCloser(bytes.NewBuffer(buffer))
			r.Body = rc

			err = utils.JSONDuplicate(json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buffer))), nil)
			if err != nil {
				logrus.Errorf("[Request] check duplicate json body request error: %s", err)
				render.Status(r, config.RR.JSONDuplicateOrInvalidFormat.HTTPStatusCode())
				render.JSON(w, r, config.RR.JSONDuplicateOrInvalidFormat.WithLocale(r))
				return
			}
		}

		// intercept writing response and store it to write log.
		buff := bytes.Buffer{}
		newResponseWriter := newResponseWriter(w, &buff)
		w = newResponseWriter
		next.ServeHTTP(w, r)

		hostname, _ := os.Hostname()

		// write response log
		logs := logrus.Fields{
			"host":            hostname,
			"method":          r.Method,
			"path":            r.URL.Path,
			"Accept-Language": context.GetLanguage(r),
			"clientIP":        GetIPAddress(r),
			"User-Agent":      r.Header.Get("User-Agent"),
			"body-size":       fmt.Sprintf("%f KB", float64(newResponseWriter.body.Len())/1024.0),
			"statusCode":      fmt.Sprintf("%d %s", newResponseWriter.Status(), http.StatusText(newResponseWriter.Status())),
			"processTime":     time.Since(start),
		}

		if user, ok := context.GetUser(r); ok {
			logs["user_id"] = user.ID
		}

		if parameters, ok := context.GetParameters(r); ok {
			formValue := reflect.ValueOf(parameters)
			if formValue.Kind() == reflect.Ptr {
				formValue = formValue.Elem()
			}

			for i := 0; i < formValue.NumField(); i++ {
				valueField := formValue.Field(i)
				if fieldName := formValue.Type().Field(i).Name; isAboutPassword(fieldName) {
					paramValue := formValue.FieldByName(fieldName)
					password, ok := valueField.Interface().(string)
					if ok {
						paramValue.Set(reflect.ValueOf(utils.WrapPassword(password)))
					}
				}
			}
			parametersByte, _ := json.Marshal(parameters)
			logs["parameters"] = string(parametersByte)
		}

		if r.URL.Path != "/api/healthz" {
			if !strings.HasPrefix(r.URL.Path, "/api/v1/swagger") && !strings.HasPrefix(r.URL.Path, "/api/v1/robot") {
				logrus.WithFields(logs).Info(fmt.Sprintf("[%s][%s] response: %+v", r.Method, r.URL.Path, newResponseWriter.body.String()))
			}
		}
	})
}

func isAboutPassword(fieldName string) bool {
	return fieldName == "Password" ||
		fieldName == "CurrentPassword" ||
		fieldName == "NewPassword"
}

type responseWriterModel struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func newResponseWriter(rw http.ResponseWriter, buffer *bytes.Buffer) *responseWriterModel {
	nrw := &responseWriterModel{
		ResponseWriter: rw,
		body:           buffer,
	}
	return nrw
}

// Write write data
func (b responseWriterModel) Write(data []byte) (int, error) {
	b.body.Write(data)
	return b.ResponseWriter.Write(data)
}

// WriteHeader write header
func (rw *responseWriterModel) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

// Status get status from responseWriterModel
func (responseWriter *responseWriterModel) Status() int {
	return responseWriter.status
}

// IPFromRequest get ip address
func IPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// GetIPAddress get ip address from request
func GetIPAddress(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if userIP, ok := IPFromRequest(r); ok == nil {
		lastIP := userIP.String()
		if clientIP == "" {
			clientIP = lastIP
		} else {
			clientIP += fmt.Sprintf(", %s", lastIP)
		}
	}

	return clientIP
}
