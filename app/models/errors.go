package models

import (
	"encoding/json"
	"forum/app/config"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
}

func SendErrorResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	if statusCode == 404 || statusCode == 405 || statusCode == 403 {
		config.Logger.Println("Rendering ", statusCode, " section in home template")
		data := struct {
			StatusCode int
			Message    string
		}{
			StatusCode: statusCode,
			Message:    message,
		}
		
		config.Templates.ExecuteTemplate(resp, "404.html", data)
		return
	} else {
		resp.Header().Set("Content-Type", "application/json")
		jsonResponse := ErrorResponse{
			StatusCode: statusCode,
			Message:    message}
		json.NewEncoder(resp).Encode(jsonResponse)
	}

}
