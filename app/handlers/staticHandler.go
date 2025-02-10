package handlers

import (
	"forum/app/config"
	"forum/app/models"
	"net/http"
	"os"
)

func Static(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		config.Logger.Println("Attempt to access static files with method: ", req.Method, " Rejected.")
		models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "405 - Method Not Allowed.")
		return
	}

	fileInfo, err := os.Stat(req.URL.Path[1:])
	if err != nil {
		if os.IsNotExist(err) {
			config.Logger.Println("File: ", req.URL.Path[1:], " not found in static folder.")
			models.SendErrorResponse(resp, http.StatusNotFound, "404 - Page Not Found")
			return
		}
		config.Logger.Println("File: ", req.URL.Path[1:], " not found in static folder.")
		models.SendErrorResponse(resp, http.StatusNotFound, "404 - Page Not Found")
		return
	}
	if fileInfo.IsDir() {
		config.Logger.Println("Attempt to access Forbidden Folder: ", fileInfo)
		models.SendErrorResponse(resp, http.StatusForbidden, "Access Forbidden")
		return
	}
	http.ServeFile(resp, req, req.URL.Path[1:])
}
