package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"objects"
	"regexp"
	"strconv"
	"time"
	"html/template"
)

func HandlerRoot(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    aRequest.ParseForm()
    if !(IsTokenValid(aResponseWriter, aRequest)) {
    	return
    }
    
    //ServeError(aResponseWriter, STR_MSG_NOTFOUND, STR_template_page_error_html)
    
     //must test req.URL.Path == "/" and process only then, 
     //because this handler will be called for every path, since "/" matches all paths starting with "/"
     
//    aRequest.ParseForm();
//	
//	cookie, errCookie := aRequest.Cookie(API_KEY_token)
//	if errCookie != nil {
//		log.Printf("handleRoot, Error reading cookie, error=%s", errCookie.Error())
//		ServeLogin(aResponseWriter, STR_MSG_login);
//		return
//	}
//	isTokenValid, userId := DbIsTokenValid(cookie.Value, nil)
//	if errCookie == nil && !isTokenValid {
//		ServeLogin(aResponseWriter, STR_MSG_login);
//		return
//	}
//	
//	user, errorUser := DbGetUserLoad(userId, nil);
//	if errorUser != nil {
//		log.Printf("errorUser=%s", errorUser.Error())
//	}
//	log.Printf("cookie.value=%s", cookie.Value)
//	
//	//Check if the file in the url path exists
//	templateFile, err := template.ParseFiles(aRequest.URL.Path[1:]);
//	if err != nil {
//		ServeError(aResponseWriter, STR_MSG_404, STR_template_page_error_html);
//	} else {	
//		AddCookie(aResponseWriter, cookie.Value)
//		if aRequest.URL.Path[1:] == "templates/Content.html" && user.Email != STR_EMPTY {
//			err = templateFile.Execute(aResponseWriter, user);
//		} else {
//			err = templateFile.Execute(aResponseWriter, nil);
//		}
//		if err != nil {
//			log.Printf("handleRoot, Error=", err.Error());
//		}
//	}
}

//func HandlerEcho(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
//	responseText := STR_EMPTY 
//	
//	aRequest.ParseForm()
//	
//	body := aRequest.Form
//	log.Printf("aRequest.Form=%s", body)
//	
//	bytesBody, err := ioutil.ReadAll(aRequest.Body)
//	if(err != nil) {
//		log.Printf("Error reading body, err=%s", err.Error())
//	} else {
//		log.Printf("bytesBody=%s", string(bytesBody))
//		responseText = string(bytesBody)
//	}
//	
//	headers := aRequest.Header
//	for key, value := range headers {
//		log.Printf("header=%s\n", key)
//		fmt.Fprintf(aResponseWriter, "Header=%s\n", key)		
//		for idx, val := range value {
//			log.Printf("idx=%d, value=%s", idx, val)
//			fmt.Fprintf(aResponseWriter, "value=%s\n", val)
//		} 
//	}
//
//	fmt.Fprintf(aResponseWriter, "Method=%s\n", aRequest.Method)
//	fmt.Fprintf(aResponseWriter, "%s\n", responseText)
//}

func HandlerMessage(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	body := aRequest.Form
	log.Printf("aRequest.Form=%s", body)
	bytesBody, err := ioutil.ReadAll(aRequest.Body)
	if(err != nil) {
		log.Printf("Error reading body, err=%s", err.Error())
	}
//	log.Printf("bytesBody=%s", string(bytesBody))
	
	//check Header Token
//	headerAuthentication := aRequest.Header.Get(STR_Authorization)
//	isValid, userId := DbIsTokenValid(headerAuthentication, nil)
//	log.Printf("HandlerMessage, headerAuthentication=%s, isValid=%t, userId=%d", headerAuthentication, isValid, userId)
//	if !isValid {
//		result := new(objects.Result)
//		result.ErrorMessage = STR_MSG_login
//		result.ResultCode = http.StatusOK
//		ServeResult(aResponseWriter, result, STR_template_result)
//		return
//	}
	
	report := new(objects.Report)
	json.Unmarshal(bytesBody, report)
	log.Printf("HandlerMessage, report.ApiKey=%s, report.ClientId=%s, report.Message=%s, report.Sequence=%d, report.Time=%d", 
			report.ApiKey, report.ClientId, report.Message, report.Sequence, report.Time)
	var isApiKeyValid = false
	if report.ApiKey != STR_EMPTY {
		isApiKeyValid, _ = IsApiKeyValid(report.ApiKey)
	}
	if !isApiKeyValid {
		result := new(objects.Result)
		result.ErrorMessage = STR_MSG_invalidapikey
		result.ResultCode = http.StatusOK
		ServeResult(aResponseWriter, result, STR_template_result)
		return
	}
	
	DbAddReport(report.ApiKey, report.ClientId, report.Time, report.Sequence, report.Message, report.FilePath, nil) 
	
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

		//Check ApiKey
		strMessage := aRequest.FormValue(API_KEY_message)
		log.Printf("strMessage=%s", strMessage)
		report := new(objects.Report)
		json.Unmarshal([]byte(strMessage), report)
		log.Printf("HandlerUploadImage, report.ApiKey=%s, report.ClientId=%s, report.Message=%s, report.Sequence=%d, report.Time=%d", 
				report.ApiKey, report.ClientId, report.Message, report.Sequence, report.Time)
		var isApiKeyValid = false
		if report.ApiKey != STR_EMPTY {
			isApiKeyValid, _ = IsApiKeyValid(report.ApiKey)
		}
		if !isApiKeyValid {
			result := new(objects.Result)
			result.ErrorMessage = STR_MSG_invalidapikey
			result.ResultCode = http.StatusOK
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
	
		//get file part
		multipartFile, multipartFileHeader, err := aRequest.FormFile(API_KEY_image)
		if err != nil {
			log.Printf("HandlerUploadImage, Error getting file from FormFile, err=%s", err.Error())
			result := new(objects.Result)
			result.ErrorMessage = err.Error()
			result.ResultCode = http.StatusBadRequest
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
		defer multipartFile.Close()
	
		imageFilePath := fmt.Sprintf(STR_filepath_upload_template, multipartFileHeader.Filename)
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
		log.Printf("HandlerUploadImage, imageFilePath=%s", imageFilePath)
		
		fileOut, errOut := os.Create(imageFilePath)
		if errOut != nil {
			log.Printf("Error creating fileOut, errOut=%s", errOut.Error())
			return
		}
		defer fileOut.Close()
	
		written, errWrite := io.Copy(fileOut, multipartFile)
		if errWrite != nil {
			log.Printf("HandlerUploadImage, Erro copying file, errWrite=%s", errWrite.Error())
			return
		}
		log.Printf("HandlerUploadImage, Bytes written=%d", written)
		
		//add report to Database
		report.FilePath = imageFilePath
		DbAddReport(report.ApiKey, report.ClientId, report.Time, report.Sequence, report.Message, report.FilePath, nil) 
		
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
		//get message part
		errorParse := aRequest.ParseMultipartForm(8388608)
		if errorParse != nil {
			log.Printf("HandlerUploadFile, errorParse=%s", errorParse.Error())
		}		
		//Check ApiKey
		strMessage := aRequest.FormValue(API_KEY_message)
		log.Printf("strMessage=%s", strMessage)
		report := new(objects.Report)
		json.Unmarshal([]byte(strMessage), report)
		log.Printf("HandlerUploadFile, report.ApiKey=%s, report.ClientId=%s, report.Message=%s, report.Sequence=%d, report.Time=%d", 
				report.ApiKey, report.ClientId, report.Message, report.Sequence, report.Time)
		var isApiKeyValid = false
		if report.ApiKey != STR_EMPTY {
			isApiKeyValid, _ = IsApiKeyValid(report.ApiKey)
		}
		if !isApiKeyValid {
			result := new(objects.Result)
			result.ErrorMessage = STR_MSG_invalidapikey
			result.ResultCode = http.StatusOK
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
		
		multipartFile, multipartFileHeader, err := aRequest.FormFile(API_KEY_file)
		if err != nil {
			log.Printf("HandlerUploadFile, Error getting file from FormFile, err=%s", err.Error())
			result := new(objects.Result)
			result.ErrorMessage = err.Error()
			result.ResultCode = http.StatusBadRequest
			ServeResult(aResponseWriter, result, STR_template_result)
			return
		}
		defer multipartFile.Close()
	
		filePath := fmt.Sprintf(STR_filepath_upload_template, multipartFileHeader.Filename)
		fileName := filePath[0:(len(filePath) - 4)]
		fileExstension := filePath[(len(filePath) - 4):len(filePath)]
		fileNum := 0;
		var errorFileExists error
		_, errorFileExists = os.Stat(filePath)
		for(!os.IsNotExist(errorFileExists)) {
			fileNum++
			filePath = fileName + strconv.Itoa(fileNum) + fileExstension
			_, errorFileExists = os.Stat(filePath)
		}
		log.Printf("HandlerUploadFile, filePath=%s", filePath)
		
		fileOut, errOut := os.Create(filePath)
		if errOut != nil {
			log.Printf("HandlerUploadFile, Error creating fileOut, errOut=%s", errOut.Error())
			return
		}
		defer fileOut.Close()
	
		written, errWrite := io.Copy(fileOut, multipartFile)
		if errWrite != nil {
			log.Printf("HandlerUploadFile, Error copying file, errWrite=%s", errWrite.Error())
			return
		}
		log.Printf("Bytes written=%d", written)
		
		//add report to Database
		report.FilePath = filePath
		DbAddReport(report.ApiKey, report.ClientId, report.Time, report.Sequence, report.Message, report.FilePath, nil)
		
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
			log.Printf("HandlerLogin, errorUser=%s", errorUser.Error())
		}
		if (userId > -1) {
			token := DbAddToken(userId, nil)
			AddCookie(responseWriter, token)
			http.Redirect(responseWriter, request, GetApiUrlListApiKeys(), 301)
		} else {
			ServeLogin(responseWriter, "Wrong username or password");
		}
	}
}

func HandlerLogout(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	if !(IsTokenValid(responseWriter, request)) {return}
	
	AddCookie(responseWriter, "no token")
	ServeLogin(responseWriter, STR_EMPTY);
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
	
	token := GetCookieToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerApiKeys, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return
	}
	
	var strApiKey string = STR_EMPTY
	sliceApiKeys := request.Form[API_KEY_apikey]
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		strApiKey = sliceApiKeys[0]
	}
	if strApiKey != STR_EMPTY {
		sliceInviteIds := request.Form[API_KEY_inviteid]
		if sliceInviteIds != nil && len(sliceInviteIds) > 0 {
			DbInviteAddApiKey(userId, sliceInviteIds[0], strApiKey, strApiKey, nil)
		}
	}
	
	var apiKeys []*objects.ApiKey
	apiKeys = DbGetApiKey(userId, nil)
	log.Printf("HandlerApiKeys, apiKeys=%s", apiKeys)
	
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
	
	log.Printf("HandlerReports, url=%s", request.URL.RawPath)
	
	token := GetCookieToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerReports, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return;
	}
	
	var strApiKey string = STR_EMPTY
	var strAppName string = STR_EMPTY
	var strClientId string = STR_EMPTY
	var startNum int
	var pageNum int
	var pageSize int
	var strPageNum []string
	var strPageSize []string
	var err error
	apiKey := request.Form[API_KEY_apikey]
	appName := request.Form[API_KEY_appname]
	clientId := request.Form[API_KEY_clientid]
	strPageNum = request.Form[API_KEY_pagenum]
	strPageSize = request.Form[API_KEY_pagesize]
	if apiKey != nil && len(apiKey) > 0 {
		strApiKey = apiKey[0]
	}
	if appName != nil && len(appName) > 0 {
		strAppName = appName[0]
	}
	if clientId != nil && len(clientId) > 0 {
		strClientId = clientId[0]
	}
	if strPageNum != nil && len(strPageNum) > 0 {
		regex, errRegEx := regexp.Compile("[^-][^0-9]")
		if errRegEx != nil {
			log.Printf("HandlerReports, errRegEx=%s", errRegEx.Error())
		}
		strPage := regex.ReplaceAllString(strPageNum[0], STR_EMPTY)
		log.Printf("HandlerReports, strPageNum[0]=%s, strPage=%s", strPageNum[0], strPage)		
		pageNum, err = strconv.Atoi(strPage)
		if err != nil {
			log.Printf("Error converting %s to int, error=%s", strPageNum, err.Error())
			pageNum = 0
		}
	} else {
		pageNum = 0
	}
	if strPageSize != nil && len(strPageSize) > 0 {
		pageSize, err = strconv.Atoi(strPageSize[0])
		if err != nil {
			log.Printf("Error converting %s to int, error=%s", strPageSize, err.Error())
			pageSize = REPORTS_PAGE_SIZE
		}	
	} else {
		pageSize = REPORTS_PAGE_SIZE
	}
	var sliceReports []*objects.Report
	
	log.Printf("HandlerReports, clientId=%s", clientId)
	log.Printf("HandlerReports, strClientId=%s", strClientId)
	if(pageNum < 0) {
		var rowCount int64
		sliceReports, rowCount = DbGetReportsLastPage(strApiKey, strClientId, pageSize, nil)
		pageNum = int(rowCount / int64(pageSize))
		log.Printf("HandlerReports, startNum=%d, pageNum=%d, pageSize=%d", startNum, pageNum, pageSize)
	} else {
		startNum = pageNum * pageSize
		var endNum int
		sliceReports, endNum = DbGetReportsByApiKey(strApiKey, strClientId, startNum, pageSize, nil)
		log.Printf("HandlerReports, startNum=%d, endNum=%d, pageNum=%d, pageSize=%d", startNum, endNum, pageNum, pageSize)
	}
	
	var clientInfo = new(objects.ClientInfo)
	clientInfo.Name = STR_All_clients
	if strClientId != STR_EMPTY {
		clientInfo = DbGetClientInfo(strApiKey, strClientId, nil)
	}

	var pagePrevious int
	var pageNext int
	var reportLast *objects.Report
	length := len(sliceReports)
	if length == 0 || length < pageSize {
		pagePrevious = pageNum - 1
		pageNext = -1
	} else {
		pagePrevious = pageNum - 1
		pageNext = pageNum + 1
	}
	if length == 0 {
		reportLast = nil
	} else {
		reportLast = sliceReports[len(sliceReports) - 1]
	}
	
	templateReport, err := template.ParseFiles(STR_template_list_reports_for_apikey_html)
	if err != nil {
		log.Printf("Error parsing template, %s, error=%s", STR_template_list_reports_for_apikey_html, err.Error())
	}
	//errorExecute := templateReport.Execute(responseWriter, sliceReports)
	var templateData = struct {
							Reports []*objects.Report 
							ApiKey string
							AppName string
							ClientId string
							PageNumStart int
							PageNumPrevious int
							PageNumNext int
							PageNumLast int
							ReportLast *objects.Report
							ClientInfo *objects.ClientInfo
							}{
								sliceReports, 
								strApiKey,
								strAppName,
								strClientId,
								pageNum,
								pagePrevious,
								pageNext,
								-1,
								reportLast,
								clientInfo,
							}
	errorExecute := templateReport.Execute(responseWriter, templateData)
	if errorExecute != nil {
		log.Printf("Error executing template, %s, error=%s", STR_template_list_reports_for_apikey_html, errorExecute.Error())
	}
}

func HandlerReportsClearConfirm(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	token := GetCookieToken(aRequest)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerReports, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(aResponseWriter, STR_MSG_login)
		return;
	}
	
	var strApiKey string = STR_EMPTY
	var strAppName string = STR_EMPTY
	sliceApiKeys := aRequest.Form[API_KEY_apikey]
	sliceAppNames := aRequest.Form[API_KEY_appname]
	log.Printf("HandlerReportsClearConfirm, sliceApiKeys=%s, sliceApiNames=%s", sliceApiKeys, sliceAppNames)
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		strApiKey = sliceApiKeys[0]
	} else {
		strApiKey = STR_EMPTY
	}
	if sliceAppNames != nil && len(sliceAppNames) > 0 {
		strAppName = sliceAppNames[0]
	}
	if strApiKey == STR_EMPTY {return}
	
	var apiKey *objects.ApiKey = new(objects.ApiKey)
	apiKey.ApiKey = strApiKey
	apiKey.AppName = strAppName
	
	templateDeleteConfirm, err := template.ParseFiles(STR_template_reports_deleteconfirm_html)
	if err != nil {
		log.Printf("HandlerReportsClearConfirm, Error parsing template %s, error=%s", STR_template_reports_deleteconfirm_html, err.Error())
	}
	errorExecute := templateDeleteConfirm.Execute(aResponseWriter, apiKey)
	if errorExecute != nil {
		log.Printf("HandlerReportsClearConfirm, Error executing template, %s, error=%s", STR_template_reports_deleteconfirm_html, errorExecute.Error())
	}
}

func HandlerReportsClear(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	token := GetCookieToken(aRequest)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerReportsClear, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(aResponseWriter, STR_MSG_login)
		return;
	}
	
	var strApiKey string
	apiKeys := aRequest.Form[API_KEY_apikey]
	if apiKeys != nil && len(apiKeys) > 0 {
		strApiKey = apiKeys[0]
	} else {
		strApiKey = STR_EMPTY
	}
	if strApiKey == STR_EMPTY {return}
	
	DbClearReports(strApiKey, nil)
	DbClearClientInfo(strApiKey, nil)
	http.Redirect(aResponseWriter, aRequest, GetApiUrlListApiKeys(), http.StatusMovedPermanently)
}

func HandlerAddApiKey(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	log.Printf("HandlerAddApiKey, url=%s", request.URL.Path)
	token := GetCookieToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerAddApiKey, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return;
	}
	
	var appName string = STR_EMPTY
	sliceAppNames := request.Form[API_KEY_appname]
	if sliceAppNames != nil && len(sliceAppNames) > 0 {
		appName = sliceAppNames[0]
	} else {
		ServeAddApiKey(responseWriter)
		return
	}
	
	log.Printf("HandlerAddApiKey, appName=%s", appName)
	isAdded := DbAddApiKey(userId, appName, nil)
	log.Printf("HadlerAddApiKey, isAdded=%t", isAdded)
	if isAdded {
		http.Redirect(responseWriter, request, GetApiUrlListApiKeys(), http.StatusMovedPermanently)
	}
}

func HandlerApiKeyDeleteConfirm(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	log.Printf("HandlerApiKeyDelete, url=%s", request.URL.Path)
	if !(IsTokenValid(responseWriter, request)) {return}
	
	var strApiKey string = STR_EMPTY
	var strAppName string = STR_EMPTY
	sliceApiKeys := request.Form[API_KEY_apikey]
	sliceAppNames := request.Form[API_KEY_appname]
	log.Printf("HandlerApiKeyDelete, sliceApiKeys=%s, sliceApiNames=%s", sliceApiKeys, sliceAppNames)
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		strApiKey = sliceApiKeys[0]
	} else {
		ServeAddApiKey(responseWriter)
		return
	}
	
	if sliceAppNames != nil && len(sliceAppNames) > 0 {
		strAppName = sliceAppNames[0]
	}
	
	var apiKey *objects.ApiKey = new(objects.ApiKey)
	apiKey.ApiKey = strApiKey
	apiKey.AppName = strAppName
	
	templateDeleteConfirm, err := template.ParseFiles(STR_template_apikey_deleteconfirm_html)
	if err != nil {
		log.Printf("HandlerApiKeyDeleteConfirm, Error parsing template %s, error=%s", STR_template_apikey_deleteconfirm_html, err.Error())
	}
	errorExecute := templateDeleteConfirm.Execute(responseWriter, apiKey)
	if errorExecute != nil {
		log.Printf("HandlerApiKeyDeleteConfirm, Error executing template, %s, error=%s", STR_template_apikey_deleteconfirm_html, errorExecute.Error())
	}
}

func HandlerApiKeyDelete(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	log.Printf("HandlerApiKeyDelete, url=%s", request.URL.Path)
	if !(IsTokenValid(responseWriter, request)) {return}
	
	var apiKey string = STR_EMPTY
	sliceApiKeys := request.Form[API_KEY_apikey]
	log.Printf("HandlerApiKeyDelete, sliceApiKeys=%s", sliceApiKeys)
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		apiKey = sliceApiKeys[0]
	} else {
		ServeAddApiKey(responseWriter)
		return
	}
	isDeleted := DbDeleteApiKey(apiKey, nil)
	log.Printf("HandlerApiKeyDelete, isDeleted=%t", isDeleted)
	http.Redirect(responseWriter, request, GetApiUrlListApiKeys(), 301)
}

func HandlerClientInfoSend(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	body := request.Form
	log.Printf("aRequest.Form=%s", body)
	bytesBody, err := ioutil.ReadAll(request.Body)
	if(err != nil) {
		log.Printf("Error reading body, err=%s", err.Error())
	}
//	log.Printf("bytesBody=%s", string(bytesBody))
	
	clientInfo := new(objects.ClientInfo)
	json.Unmarshal(bytesBody, clientInfo)
	log.Printf("HandlerClientInfo, clientInfo.ApiKey=%s, clientInfo.ClientId=%s, clientInfo.Name=%s, clientInfo.Manufacturer=%s, clientInfo.Model=%s, clientInfo.DeviceId=%s", 
			clientInfo.ApiKey, clientInfo.ClientId, clientInfo.Name, clientInfo.Manufacturer, clientInfo.Model, clientInfo.DeviceId)
	var isApiKeyValid = false
	if clientInfo.ApiKey != STR_EMPTY {
		isApiKeyValid, _ = IsApiKeyValid(clientInfo.ApiKey)
	}
	if !isApiKeyValid {
		result := new(objects.Result)
		result.ErrorMessage = STR_MSG_invalidapikey
		result.ResultCode = http.StatusOK
		ServeResult(responseWriter, result, STR_template_result)
		return
	}
	
	errorAdd := DbAddClientInfo(clientInfo.ApiKey, clientInfo.ClientId, clientInfo.Name, clientInfo.Manufacturer, 
							clientInfo.Model, clientInfo.DeviceId, nil)
	if errorAdd != nil {
		log.Printf("HandlerClientInfo, errorAdd=%s", errorAdd.Error())
		ServeError(responseWriter, errorAdd.Error(), STR_template_page_error_html)
	}
}

func HandlerClientInfoUpdate(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	body := request.Form
	log.Printf("aRequest.Form=%s", body)
	var strApiKey string
	var strClientId string
	var strName string
	apiKeys := body[API_KEY_apikey]
	clientIds := body[API_KEY_clientid]
	names := body[API_KEY_name]
	if apiKeys != nil && len(apiKeys) > 0 {
		strApiKey = apiKeys[0]
	} else {
		ServeError(responseWriter, STR_MSG_invalidapikey, STR_template_page_error_html)
		return
	}
	if clientIds != nil && len(clientIds) > 0 {
		strClientId = clientIds[0]
	} else {
		ServeError(responseWriter, STR_MSG_invalidclientid, STR_template_page_error_html)
		return
	}
	if names != nil && len(names) > 0 {
		strName = names[0]
	}

	var isApiKeyValid = false
	if strApiKey != STR_EMPTY {
		isApiKeyValid, _ = IsApiKeyValid(strApiKey)
	}
	if !isApiKeyValid {
		result := new(objects.Result)
		result.ErrorMessage = STR_MSG_invalidapikey
		result.ResultCode = http.StatusOK
		ServeResult(responseWriter, result, STR_template_result)
		return
	}
	
	errorUpdate := DbUpdateClientInfo(strApiKey, strClientId, strName, nil)
	if errorUpdate != nil {
		log.Printf("HandlerClientInfo, errorUpdate=%s", errorUpdate.Error())
		ServeError(responseWriter, errorUpdate.Error(), STR_template_page_error_html)
		return
	}
	
	http.Redirect(responseWriter, request, (GetApiUrlListClientIds() + "?apikey=" + strApiKey + "&clientid=" + strClientId), 301)
}

func HandlerClientIds(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	var apiKey string = STR_EMPTY
	sliceApiKeys := request.Form[API_KEY_apikey]
	log.Printf("HandlerClientIds, sliceApiKeys=%s", sliceApiKeys)
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		apiKey = sliceApiKeys[0]
	} else {
		ServeLogin(responseWriter, STR_MSG_login)
		return
	}
	sliceClientInfos := DbGetClientInfos(apiKey, nil)
	if sliceClientInfos == nil {
		sliceClientInfos = make([]*objects.ClientInfo, 0, 16)
//		ServeError(responseWriter, STR_MSG_server_error, STR_template_page_error_html)
//		return
	}
	
	templateClientIds, err := template.ParseFiles(STR_template_list_clientids_for_apikey_html)
	if err != nil {
		log.Printf("HandlerClientIds, Error parsing template, %s, error=%s", STR_template_list_clientids_for_apikey_html, err.Error())
	}
	var templateData = struct {
							ClientIds []*objects.ClientInfo 
							ApiKey string
							}{
								sliceClientInfos, 
								apiKey,
							}
	errorExecute := templateClientIds.Execute(responseWriter, templateData)
	if errorExecute != nil {
		log.Printf("Error executing template, %s, error=%s", STR_template_list_clientids_for_apikey_html, errorExecute.Error())
		ServeError(responseWriter, STR_MSG_server_error, STR_template_page_error_html)
	}
}

func HandlerDownload(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	token := GetCookieToken(aRequest)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerDownload, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(aResponseWriter, STR_MSG_login)
		return;
	}
	
	var strApiKey string
	var strClientId string = STR_EMPTY
	var pageSize int
	var strPageSize []string
	var err error
	apiKey := aRequest.Form[API_KEY_apikey]
	clientId := aRequest.Form[API_KEY_clientid]
	strPageSize = aRequest.Form[API_KEY_pagesize]
	if apiKey != nil && len(apiKey) > 0 {
		strApiKey = apiKey[0]
	} else {
		strApiKey = STR_EMPTY
	}
	if clientId != nil && len(clientId) > 0 {
		strClientId = clientId[0]
	}
	if strPageSize != nil && len(strPageSize) > 0 {
		pageSize, err = strconv.Atoi(strPageSize[0])
		if err != nil {
			log.Printf("HandlerDownload, Error converting %s to int, error=%s", strPageSize, err.Error())
			pageSize = REPORTS_PAGE_SIZE
		}	
	} else {
		pageSize = REPORTS_PAGE_SIZE
	}
	log.Printf("HandlerDownload, clientId=%s", clientId)
	log.Printf("HandlerDownload, strClientId=%s", strClientId)
	
	var sliceReports []*objects.Report

	var buffer = new(bytes.Buffer)
	var hasMore = true
	var startNum = 0
	for ; hasMore; {
		sliceReports, startNum = DbGetReportsByApiKey(strApiKey, strClientId, startNum, pageSize, nil)
		log.Printf("HandlerDownload, startNum=%d, pageSize=%d", startNum, pageSize)	
		
		var written int
		count := len(sliceReports)
		for x := 0; x < count; x++ {
			report := sliceReports[x]
			written, err = buffer.WriteString("Id=")
			log.Printf("HandlerDownload, written=%s", written)
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(strconv.Itoa(report.Id))
			if err != nil {
				log.Printf("HandlerDownload, error writing report.Id=%s, error=%s", report.Id, err.Error())
			}
			written, err = buffer.WriteString(", ClientId=")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(report.ClientId)
			if err != nil {
				log.Printf("HandlerDownload, error writing report.ClientId=%s, error=%s", report.ClientId, err.Error())
			}
			written, err = buffer.WriteString(", Time=")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(strconv.FormatInt(report.Time, 10))
			if err != nil {
				log.Printf("HandlerDownload, error writing report.Time=%d, error=%s", report.Time, err.Error())
			}
			written, err = buffer.WriteString(", Sequence=")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(strconv.Itoa(report.Sequence))
			if err != nil {
				log.Printf("HandlerDownload, error writing report.Sequence=%d, error=%s", report.Sequence, err.Error())
			}
			written, err = buffer.WriteString(", Message=")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(report.Message)
			if err != nil {
				log.Printf("HandlerDownload, error writing report.Message=%s, error=%s", report.Message, err.Error())
			}
			written, err = buffer.WriteString(", FilePath=")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
			written, err = buffer.WriteString(report.FilePath)
			if err != nil {
				log.Printf("HandlerDownload, error writing report.FilePath=%s, error=%s", report.FilePath, err.Error())
			}
			written, err = buffer.WriteString("\n")
			if err != nil {
				log.Printf("HandlerDownload, error writing, error=%s", err.Error())
			}
		}
		
		if count == 0 || count < pageSize {
			hasMore = false
		}
	}
	data := buffer.Bytes()
	
    aResponseWriter.Header().Set("Content-Type", "application/octet-stream")
    aResponseWriter.Header().Set("Content-Disposition", "attachment; filename=" + STR_export + ".txt")
    aResponseWriter.Header().Set("Content-Transfer-Encoding", "binary")
    aResponseWriter.Header().Set("Expires", "0")
    http.ServeContent(aResponseWriter, aRequest, STR_EMPTY, time.Now(), bytes.NewReader(data))
}

func HandlerFileLogDelete(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
	aRequest.ParseForm()
	
	token := GetCookieToken(aRequest)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerDownload, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(aResponseWriter, STR_MSG_login)
		return;
	}
	
	var err error
	err = os.Remove("resources/logs/logfile.txt")
	if err != nil {
		log.Printf("main_stunt, init, error Removing file, error=%s", err.Error())
	} else {
		FileLogCreate()
	}
}

func HandlerInviteCreate(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	
	token := GetCookieToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("HandlerInvite, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return
	}
	
	var strApiKey string = STR_EMPTY
	var strAppName string = STR_EMPTY
	sliceApiKeys := request.Form[API_KEY_apikey]
	sliceAppNames := request.Form[API_KEY_appname]
	if sliceApiKeys != nil && len(sliceApiKeys) > 0 {
		strApiKey = sliceApiKeys[0]
	}
	if sliceAppNames != nil && len(sliceAppNames) > 0 {
		strAppName = sliceAppNames[0]
	}
	
	var inviteId = DbInviteAdd(strApiKey, nil)
	
	templateApiKeys, err := template.ParseFiles(STR_template_invite_html)
	if err != nil {
		log.Printf("Error parsing template %s, error=%s", STR_template_invite_html, err.Error())
	}
	var templateData = struct {
							ApiUrlInvite string
							InviteId string
							ApiKey string
							AppName string
						}{
							GetApiUrlInvite(),
							inviteId,
							strApiKey,
							strAppName,
						}
	errorExecute := templateApiKeys.Execute(responseWriter, templateData)
	if errorExecute != nil {
		log.Printf("Error executing template, %s, error=%s", STR_template_invite_html, errorExecute.Error())
	}
}