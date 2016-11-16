package main 

import (
	"log"
    "net/http"
	"os"
    "utils"
)

func init() {
	log.SetFlags(log.Lshortfile)
	utils.DbInit()
	//utils.FileLogCreate()
}

func main() {

	utils.Port = os.Getenv(utils.STR_PORT)
	if utils.Port == utils.STR_EMPTY {
        utils.Port = utils.PORT_8080
    }

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
    http.HandleFunc(utils.PATH_filelog_delete, utils.HandlerFileLogDelete)
    
    //http.Handle("/templates/", http.FileServer(http.Dir("./templates")))
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
    
    http.ListenAndServe(":" + utils.Port, nil)
    
    defer func(){
    	utils.FileLog.Close()
	}()
}

