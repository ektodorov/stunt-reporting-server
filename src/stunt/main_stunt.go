package main 

import (
	"log"
    "net/http"
    "utils"
)

func init() {
	log.SetFlags(log.Lshortfile)
	utils.DbInit()
}

func main() {
    http.HandleFunc(utils.PATH_ROOT, utils.HandlerRoot)
//    http.HandleFunc(utils.PATH_ECHO, utils.HandlerEcho)
    http.HandleFunc(utils.PATH_MESSAGE, utils.HandlerMessage)
    http.HandleFunc(utils.PATH_UPLOADIMAGE, utils.HandlerUploadImage)
    http.HandleFunc(utils.PATH_UPLOADFILE, utils.HandlerUploadFile)
    http.HandleFunc(utils.PATH_Login, utils.HandlerLogin)
    http.HandleFunc(utils.PATH_logout, utils.HandlerLogout)
    http.HandleFunc(utils.PATH_Register, utils.HandlerRegister)
    http.HandleFunc(utils.PATH_ApiKeys, utils.HandlerApiKeys)
    http.HandleFunc(utils.PATH_Reports, utils.HandlerReports)
    http.HandleFunc(utils.PATH_AddApiKey, utils.HandlerAddApiKey)
    http.HandleFunc(utils.PATH_ApiKeyDeleteConfirm, utils.HandlerApiKeyDeleteConfirm)
    http.HandleFunc(utils.PATH_ApiKeyDelete, utils.HandlerApiKeyDelete)
    http.HandleFunc(utils.PATH_CLIENTINFO, utils.HandlerClientInfoSend)
    http.HandleFunc(utils.PATH_ClientInfoUpdate, utils.HandlerClientInfoUpdate)
    http.HandleFunc(utils.PATH_ClientIds, utils.HandlerClientIds)
    http.HandleFunc(utils.PATH_download, utils.HandlerDownload)
    
    //http.Handle("/templates/", http.FileServer(http.Dir("./templates")))
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
    
    http.ListenAndServe(":8080", nil)
}

