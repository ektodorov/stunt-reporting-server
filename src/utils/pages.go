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

func ServeLogin(responseWriter http.ResponseWriter, message string) {
	loginTemplate, err := template.ParseFiles(STR_templates_login_html);
	if err != nil {
		log.Printf("ServeLogin, Error=%s", err.Error());
	}
	
	msg := new(objects.MessageHolder);
	msg.Message = message;
	err = loginTemplate.Execute(responseWriter, msg);
	if err != nil {
		log.Printf("ServeLogin, Error=%s", err.Error());
	}
}

func ServeRegister(responseWriter http.ResponseWriter, message string) {
	registerTemplate, err := template.ParseFiles(STR_templates_register_html);
	if err != nil {
		log.Printf("ServeRegister, Error=%s", err.Error());
	}
	
	msg := new(objects.MessageHolder);
	msg.Message = message;
	err = registerTemplate.Execute(responseWriter, msg);
	if err != nil {
		log.Printf("ServeRegister, Error=%s", err.Error());
	}
}

func ServeContent(responseWriter http.ResponseWriter, userName string) {
	pageTemplate, err := template.ParseFiles(STR_templates_Content_html);
	if err != nil {
		log.Printf("ServeContent, Error=%s", err.Error());
	}
		
	user := new(objects.User);
	user.Email = userName;
	err = pageTemplate.Execute(responseWriter, user);
	if err != nil {
		log.Printf("ServeContent, Error=%s", err.Error());
	}
}

func AddCookie(responseWriter http.ResponseWriter, token string) {
	cookie := new(http.Cookie)
	cookie.Name = API_KEY_token
	cookie.Value = token
	cookie.Domain = "localhost"
	cookie.MaxAge = TOKEN_VALIDITY_SECONDS
	cookie.Path = "/"
	http.SetCookie(responseWriter, cookie)
}