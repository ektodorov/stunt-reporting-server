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
    http.ListenAndServe(":8080", nil)
}

