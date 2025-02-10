package config

import (
	"os"
	"text/template"
)


var Templates *template.Template

func InitTemplates(pattern string) {
	Logger.Println("Attempting to parse templates with pattern: ", pattern)

	Templates = template.New("forum")

	var err error
	Templates, err = Templates.ParseGlob(pattern)
	if err != nil {
		Logger.Println("Error parsing templates: ", err)
		Logger.Println("Detailed info: Failed to load templates using pattern: ", pattern)
		os.Exit(1)
	}

	Logger.Println("Templates loaded successfully")
}
