// Package formjson provides Middleware for converting posted x-www-form-urlencoded data into json
package formjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"

	"github.com/mar1n3r0/go-api-boilerplate/pkg/errors"
	"github.com/mar1n3r0/go-api-boilerplate/pkg/http/response"
	"github.com/vardius/gorouter/v4"
)

type FormError struct {
	Error   error
	Message string
}

//Provides "x-www-form-urlencoded" to "json" conversion middleware for gorouter
func FormJson() gorouter.MiddlewareFunc {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if mediatype == "application/x-www-form-urlencoded" {
				defer r.Body.Close()
				// get body
				buf, _ := ioutil.ReadAll(r.Body)

				params, err := url.ParseQuery(string(buf))
				if err != nil {
					log.Fatal(err)
					return
				}
				// map body form data
				jsonMap := map[string]string{}
				for key, val := range params {
					if len(val[0]) > 0 {
						jsonMap[key] = val[0]
					}
				}

				//marshal json
				jsonString, err := json.Marshal(jsonMap)
				if err != nil {
					//error marshalling, skip to handler
					conversionError(r, w)
					return
				}

				//write new body
				r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(jsonString))))

				//convert content-type header
				r.Header.Set("Content-Type", "application/json")
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

func conversionError(r *http.Request, w http.ResponseWriter) {
	response.RespondJSONError(r.Context(), w, errors.New(errors.INTERNAL, "Error converting form data"))
	return
}
