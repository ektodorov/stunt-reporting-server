package main 

import (
	"log"
    "net/http"
    "utils"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
    http.HandleFunc(utils.PATH_ROOT, utils.HandlerRoot)
    http.HandleFunc(utils.PATH_ECHO, utils.HandlerEcho)
    http.HandleFunc(utils.PATH_UPLOAD, utils.HandlerUpload)
    http.ListenAndServe(":8080", nil)
}

