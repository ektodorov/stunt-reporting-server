package utils

import (
	"html/template"
	"log"
	"net/http"
	"objects"
)

func ServeError(aResponseWriter http.ResponseWriter, aMessage string, aTemplateFilePath string) {
	templateError, err := template.ParseFiles(aTemplateFilePath)
	if err != nil {
		log.Printf("ServeError, Error parsing template, filePath=%s, error=%s", aTemplateFilePath, err.Error())
		return
	}
	
	msg := new(objects.Message)
	msg.Message = aMessage
	errorExecute := templateError.Execute(aResponseWriter, msg)
	if errorExecute != nil {
		log.Printf("ServeError, Error executing template, filePath=%s, error=%s", STR_template_page_error_html, errorExecute.Error())
	}
}

func ServeResult(aResponseWriter http.ResponseWriter, aResult *objects.Result, aTemplateFilePath string) {
	templateResult, err := template.ParseFiles(aTemplateFilePath)
	if err != nil {
		log.Printf("ServeError, Error parsing template, filePath=%s, error=%s", STR_template_page_error_html, err.Error())
		return
	}
	
	 errorExecute := templateResult.Execute(aResponseWriter, aResult)
	 if(errorExecute != nil) {
	 	log.Printf("ServeResult, Error executing template, filePath=%s, error=%s", aTemplateFilePath, errorExecute.Error())
	 }
}