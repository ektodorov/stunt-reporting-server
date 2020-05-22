package utils

import (
	"html/template"
	"log"
	"net/http"
	"objects"
)

//ServeError serves error HTML page to web portal clients.
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

//ServeResult serves result response to report clients.
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

//ServeLogin serves login page to web portal clients.
func ServeLogin(responseWriter http.ResponseWriter, message string) {
	loginTemplate, err := template.ParseFiles(STR_templates_login_html);
	if err != nil {
		log.Printf("ServeLogin, Error=%s", err.Error());
	}
	
	msg := new(objects.Message);
	msg.Message = message;
	err = loginTemplate.Execute(responseWriter, msg);
	if err != nil {
		log.Printf("ServeLogin, Error=%s", err.Error());
	}
}

//ServeRegister serves register page to web portal clients.
func ServeRegister(responseWriter http.ResponseWriter, message string) {
	registerTemplate, err := template.ParseFiles(STR_templates_register_html)
	if err != nil {
		log.Printf("ServeRegister, Error=%s", err.Error())
	}
	
	msg := new(objects.Message)
	msg.Message = message
	err = registerTemplate.Execute(responseWriter, msg)
	if err != nil {
		log.Printf("ServeRegister, Error=%s", err.Error())
	}
}

//ServeAddApiKey serves add API key page to web portal clients.
func ServeAddApiKey(responseWriter http.ResponseWriter) {
	addApiKeyTemplate, err := template.ParseFiles(STR_template_add_apikey_html)
	if err != nil {
		log.Printf("ServeAddApiKey, Error=%s", err.Error())
	}
	err = addApiKeyTemplate.Execute(responseWriter, nil)
	if err != nil {
		log.Printf("ServeAddApiKey, Error=%s", err.Error())
	}
}
