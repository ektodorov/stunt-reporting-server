package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"objects"
	"strconv"
	"html/template"
)

func HandlerRoot(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    //ServeError(aResponseWriter, STR_MSG_NOTFOUND, STR_template_page_error_html)
    
    aRequest.ParseForm();
	
	cookie, errCookie := aRequest.Cookie(API_KEY_token)
	if errCookie != nil {
		log.Printf("handleRoot, Error reading cookie, error=%s", errCookie.Error())
		ServeLogin(aResponseWriter, STR_MSG_login);
		return
	}
	isTokenValid, userId := DbIsTokenValid(cookie.Value, nil)
	if errCookie == nil && !isTokenValid {
		ServeLogin(aResponseWriter, STR_MSG_login);
		return
	}
	
	user, errorUser := DbGetUserLoad(userId, nil);
	if errorUser != nil {
		log.Printf("errorUser=%s", errorUser.Error())
	}
	log.Printf("cookie.value=%s", cookie.Value)
	
	//Check if the file in the url path exists
	templateFile, err := template.ParseFiles(aRequest.URL.Path[1:]);
	if err != nil {
		ServeError(aResponseWriter, STR_MSG_404, STR_template_page_error_html);
	} else {	
		AddCookie(aResponseWriter, cookie.Value)
		if aRequest.URL.Path[1:] == "templates/Content.html" && user.Email != STR_EMPTY {
			err = templateFile.Execute(aResponseWriter, user);
		} else {
			err = templateFile.Execute(aResponseWriter, nil);
		}
		if err != nil {
			log.Printf("handleRoot, Error=", err.Error());
		}
	}
}

func HandlerEcho(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	responseText := STR_EMPTY 
	
	aRequest.ParseForm()
	
	body := aRequest.Form
	log.Printf("aRequest.Form=%s", body)
	
	bytesBody, err := ioutil.ReadAll(aRequest.Body)
	if(err != nil) {
		log.Printf("Error reading body, err=%s", err.Error())
	} else {
		log.Printf("bytesBody=%s", string(bytesBody))
		responseText = string(bytesBody)
	}
	
	headers := aRequest.Header
	for key, value := range headers {
		log.Printf("header=%s\n", key)
		fmt.Fprintf(aResponseWriter, "Header=%s\n", key)		
		for idx, val := range value {
			log.Printf("idx=%d, value=%s", idx, val)
			fmt.Fprintf(aResponseWriter, "value=%s\n", val)
		} 
	}

	fmt.Fprintf(aResponseWriter, "Method=%s\n", aRequest.Method)
	fmt.Fprintf(aResponseWriter, "%s\n", responseText)
}

func HandlerMessage(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	body := aRequest.Form
	log.Printf("aRequest.Form=%s", body)
	bytesBody, err := ioutil.ReadAll(aRequest.Body)
	if(err != nil) {
		log.Printf("Error reading body, err=%s", err.Error())
	} else {
		log.Printf("bytesBody=%s", string(bytesBody))
	}
	
	headerAuthentication := aRequest.Header.Get(STR_Authorization)
	log.Printf("headerAuthentication=%s", headerAuthentication)
	
	reportMessage := new(objects.ReportMessage)
	json.Unmarshal(bytesBody, reportMessage)
	log.Printf("report.Message=%s, report.Sequence=%d, report.Time=%d", reportMessage.Message, reportMessage.Sequence, reportMessage.Time)
	
	result := new(objects.Result)
	result.ErrorMessage = STR_EMPTY
	result.ResultCode = http.StatusOK
	ServeResult(aResponseWriter, result, STR_template_result)
}

func HandlerUploadImage(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	requestMethod := aRequest.Method
	if requestMethod == STR_GET {
		result := new(objects.Result)
		result.ErrorMessage = STR_error
		result.ResultCode = http.StatusMethodNotAllowed
		ServeResult(aResponseWriter, result, STR_template_result)
	} else if requestMethod == STR_POST {
		//get message part
		errorParse := aRequest.ParseMultipartForm(8388608)
		if errorParse != nil {
			log.Printf("errorParse=%s", errorParse.Error())
		}
//		myform := aRequest.MultipartForm
//		valuesMap := myform.Value //map[string][]string
//		arrayMessage := valuesMap["message"]
//		log.Printf("arrayMessage=%d", len(arrayMessage))
		
		strMessage := aRequest.FormValue(API_KEY_message)
		log.Printf("strMessage=%s", strMessage)
	
		//get file part
		multipartFile, multipartFileHeader, err := aRequest.FormFile(API_KEY_image)
		if err != nil {
			log.Printf("Error getting file from FormFile, err=%s", err.Error())
			result := new(objects.Result)
			result.ErrorMessage = err.Error()
			result.ResultCode = http.StatusBadRequest
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
		defer multipartFile.Close()
	
		imageFilePath := fmt.Sprintf(STR_img_filepathSave_template, multipartFileHeader.Filename)
		fileName := imageFilePath[0:(len(imageFilePath) - 4)]
		fileExstension := imageFilePath[(len(imageFilePath) - 4):len(imageFilePath)]
		fileNum := 0;
		var errorFileExists error
		_, errorFileExists = os.Stat(imageFilePath)
		for(!os.IsNotExist(errorFileExists)) {
			fileNum++
			imageFilePath = fileName + strconv.Itoa(fileNum) + fileExstension
			_, errorFileExists = os.Stat(imageFilePath)
		}
		log.Printf("imageFilePath=%s", imageFilePath)
		
		fileOut, errOut := os.Create(imageFilePath)
		if errOut != nil {
			log.Printf("Error creating fileOut, errOut=%s", errOut.Error())
			return
		}
		defer fileOut.Close()
	
		written, errWrite := io.Copy(fileOut, multipartFile)
		if errWrite != nil {
			log.Printf("Erro copying file, errWrite=%s", errWrite.Error())
			return
		}
	
		log.Printf("Bytes written=%d", written)
		
		result := new(objects.Result)
		result.ErrorMessage = STR_EMPTY
		result.ResultCode = http.StatusOK
		ServeResult(aResponseWriter, result, STR_template_result)
	}
}

func HandlerUploadFile(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	requestMethod := aRequest.Method
	if requestMethod == STR_GET {
		result := new(objects.Result)
		result.ErrorMessage = STR_error
		result.ResultCode = http.StatusMethodNotAllowed
		ServeResult(aResponseWriter, result, STR_template_result)
	} else if requestMethod == STR_POST {
		multipartFile, multipartFileHeader, err := aRequest.FormFile(API_KEY_file)
		if err != nil {
			log.Printf("Error getting file from FormFile, err=%s", err.Error())
			result := new(objects.Result)
			result.ErrorMessage = err.Error()
			result.ResultCode = http.StatusBadRequest
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
		defer multipartFile.Close()
	
		imageFilePath := fmt.Sprintf(STR_img_filepathSave_template, multipartFileHeader.Filename)
		fileName := imageFilePath[0:(len(imageFilePath) - 4)]
		fileExstension := imageFilePath[(len(imageFilePath) - 4):len(imageFilePath)]
		fileNum := 0;
		var errorFileExists error
		_, errorFileExists = os.Stat(imageFilePath)
		for(!os.IsNotExist(errorFileExists)) {
			fileNum++
			imageFilePath = fileName + strconv.Itoa(fileNum) + fileExstension
			_, errorFileExists = os.Stat(imageFilePath)
		}
		log.Printf("imageFilePath=%s", imageFilePath)
		
		fileOut, errOut := os.Create(imageFilePath)
		if errOut != nil {
			log.Printf("Error creating fileOut, errOut=%s", errOut.Error())
			return
		}
		defer fileOut.Close()
	
		written, errWrite := io.Copy(fileOut, multipartFile)
		if errWrite != nil {
			log.Printf("Erro copying file, errWrite=%s", errWrite.Error())
			return
		}
	
		log.Printf("Bytes written=%d", written)
		result := new(objects.Result)
		result.ErrorMessage = STR_EMPTY
		result.ResultCode = http.StatusOK
		ServeResult(aResponseWriter, result, STR_template_result)
	}
}

func HandlerLogin(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm();
	
	if request.Method == STR_GET {
		ServeLogin(responseWriter, STR_EMPTY);	
	} else {
		var userName string = request.FormValue(API_KEY_username);
		var password string = request.FormValue(API_KEY_password);
		if userName == STR_EMPTY || password == STR_EMPTY {
			ServeLogin(responseWriter, "Please enter username and password");
			return;
		}
		
		var userId = -1;
		var errorUser error = nil
		userId, errorUser = DbGetUser(userName, password, nil)
		if errorUser != nil {
			log.Printf("handlerLogin, errorUser=%s", errorUser.Error())
		}
		if (userId > -1) {
			token := DbAddToken(userId, nil)
			AddCookie(responseWriter, token)
			http.Redirect(responseWriter, request, API_URL_Content, 301)
		} else {
			ServeLogin(responseWriter, "Wrong username or password");
		}
	}
}

func HandlerRegister(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	if request.Method == STR_GET {
		ServeRegister(responseWriter, STR_EMPTY);
	} else {
		email := request.FormValue(API_KEY_email);
		password := request.FormValue(API_KEY_password);
		if (email == STR_EMPTY || password == STR_EMPTY) {
			ServeRegister(responseWriter, STR_MSG_register);
			return;
		}
		
		isUserExists, isUserAdded, errorUser := DbAddUser(email, password, nil);
		if errorUser != nil {
			log.Printf("handleRegister, errorUser=%s", errorUser.Error())
		}
		if isUserExists {
			ServeRegister(responseWriter, "Username is already taken.");
		} else if isUserAdded == false {
			ServeRegister(responseWriter, "Cannot create user.");
		} else {
			ServeLogin(responseWriter, STR_EMPTY);
		}
	}
}

func HandlerApiKeys(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	token := GetHeaderToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return;
	}
	
	var apiKeys []string
	apiKeys = DbGetApiKey(userId, nil)
	
	templateApiKeys, err := template.ParseFiles(STR_template_list_apikeys_html)
	if err != nil {
		log.Printf("Error parsing template %s, error=%s", STR_template_list_apikeys_html, err.Error())
	}
	errorExecute := templateApiKeys.Execute(responseWriter, apiKeys)
	if errorExecute != nil {
		log.Printf("Error executing template, %s, error=%s", STR_template_list_apikeys_html, errorExecute.Error())
	}
}

func HandlerReports(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	apiKey := request.Form["apiKey"]
	strStartNum := request.Form["startNum"]
	strPageSize := request.Form["pageSize"]
	startNum, errStartNum := strconv.Atoi(strStartNum[0])
	if errStartNum != nil {
		log.Printf("Error converting %s to int, error=%s", strStartNum, errStartNum.Error())
	}
	pageSize, errPageSize := strconv.Atoi(strPageSize[0])
	if errPageSize != nil {
		log.Printf("Error converting %s to int, error=%s", strPageSize, errPageSize.Error())
	}
	var sliceReports []*objects.Report
	var endNum int
	sliceReports, endNum = DbGetReportsByApiKey(apiKey[0], "clientId", startNum, pageSize, nil)
	log.Printf("HandlerReports, endNum=%d", endNum)
	
	templateReport, err := template.ParseFiles(STR_template_list_reports_for_apikey)
	if err != nil {
		log.Printf("Error parsing template, %s, error=%s", STR_template_list_reports_for_apikey, err.Error())
	}
	errorExecute := templateReport.Execute(responseWriter, sliceReports)
	if errorExecute != nil {
		log.Printf("Error executing template, %s, error=%s", STR_template_list_reports_for_apikey, errorExecute.Error())
	}
}

/* Utils */
func GetHeaderToken(aRequest *http.Request) string {
	headers := aRequest.Header
	tokens := headers["Token"]
	token := tokens[0]
	return token
}