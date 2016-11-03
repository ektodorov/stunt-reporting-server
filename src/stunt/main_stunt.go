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
    http.HandleFunc(utils.PATH_ECHO, utils.HandlerEcho)
    http.HandleFunc(utils.PATH_MESSAGE, utils.HandlerMessage)
    http.HandleFunc(utils.PATH_UPLOADIMAGE, utils.HandlerUploadImage)
    http.HandleFunc(utils.PATH_UPLOADFILE, utils.HandlerUploadFile)
    http.HandleFunc(utils.PATH_Login, utils.HandlerLogin)
    http.HandleFunc(utils.PATH_Register, utils.HandlerRegister)
    http.HandleFunc(utils.PATH_ApiKeys, utils.HandlerApiKeys)
    http.HandleFunc(utils.PATH_Reports, utils.HandlerReports)
    
    //http.Handle("/templates/", http.FileServer(http.Dir("./templates")))
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
    
    http.ListenAndServe(":8080", nil)
}

