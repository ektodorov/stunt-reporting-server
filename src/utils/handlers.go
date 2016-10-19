package utils

import (
	"io"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"objects"
	"strconv"
)

func HandlerRoot(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    ServeError(aResponseWriter, STR_MSG_NOTFOUND, STR_template_page_error_html)
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
		log.Printf("header=%s", key)
		fmt.Fprintf(aResponseWriter, "Header=%s\n", key)		
		for idx, val := range value {
			log.Printf("idx=%d, value=%s", idx, val)
			fmt.Fprintf(aResponseWriter, "value=%s\n", val)
		} 
	}

	fmt.Fprintf(aResponseWriter, "Method=%s \n", aRequest.Method)
	fmt.Fprintf(aResponseWriter, "%s", responseText)
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

/* Utils */
func getHeaderToken(aRequest *http.Request) string {
	headers := aRequest.Header
	tokens := headers["Token"]
	token := tokens[0]
	return token
}