/*
	Package is to support getting an IP address from a client's browser and returing IP Address
	It provides functionality comparable to http://jsonip.appspot.com
*/

package getipaddress

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/GetClientIPAddress", errorHandler(handlerGetClientIPAddress))
}

func handlerGetClientIPAddress(w http.ResponseWriter, r *http.Request) error {
	//Verify that request has an origin handler
	if r.Header.Get("Origin") == "" {
		return newErrorHttp(http.StatusBadRequest, "Cross domain request require Origin header")
	}

	//Verify that the request method is a GET
	if r.Method != "GET" {
		return newErrorHttp(http.StatusBadRequest, "Cross domain request only supports GET")
	}

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "no-cache")

	//Get IP address from request and set it in return data
	retData := &getClientIPAddressReturn{IP: r.RemoteAddr}
	encodedRetData, err := json.Marshal(retData)
	if err != nil {
		return newErrorHttp(http.StatusInternalServerError, err.Error())
	}

	//Return the response to the user
	_, err = w.Write(encodedRetData)
	if err != nil {
		return newErrorHttp(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func errorHandler(f httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					http.Error(w, fmt.Sprintf("%v", rec), http.StatusInternalServerError)
				}
			}
		}()
		err := f(w, r)
		if err != nil {
			switch err := err.(type) {
			case *errorHttp:
				http.Error(w, err.message, err.code)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

type getClientIPAddressReturn struct {
	IP string
}

type errorHttp struct {
	code    int
	message string
}

func (f *errorHttp) Error() string {
	return fmt.Sprintf("%d - %s", f.code, f.message)
}

func newErrorHttp(code int, message string) error {
	return &errorHttp{code: code, message: message}
}

type httpHandler func(http.ResponseWriter, *http.Request) error

//Run on dev command: E:\>E:\Go\go_appengine\dev_appserver.py E:\Go\GoExternal\src\appengine\getipaddress
//Update comand: E:\>E:\Go\go_appengine\appcfg.py update E:\Go\GoExternal\src\appengine\getipaddress
