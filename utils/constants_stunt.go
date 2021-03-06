package utils

import (
	"crypto/sha1"
	"github.com/nu7hatch/gouuid"
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"os"
	"strings"
	"time"
)

var Port string
var FileLog *os.File
var Letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
const SALT_LENGTH = 32
const TOKEN_VALIDITY_SECONDS = 60 * 60 * 24
const TOKEN_VALIDITY_MS = TOKEN_VALIDITY_SECONDS * 1000
const REPORTS_PAGE_SIZE = 100
const INVITE_VALIDITY_MS = 60 * 15 * 1000

const STR_EMPTY = ""
const STR_BLANK = " "
const STR_MSG_NOTFOUND = "404 Not found"
const STR_GET = "GET"
const STR_POST = "POST"
const STR_error = "error"
const STR_Authorization = "Authorization"
const STR_symbol_dash = "-"
const STR_symbol_dot = "."
const STR_symbol_colon = ":"
const STR_id = "id"
const STR_All_clients = "All clients"
const STR_export = "export"
const STR_txt = "txt"
const STR_MSG_404 = "404, Page not found"
const STR_MSG_login = "Please login"
const STR_MSG_register = "Please enter email and password"
const STR_MSG_invalidapikey = "Invalid api key"
const STR_MSG_invalidclientid = "Invalid client id"
const STR_MSG_server_error = "500 Server error"
const STR_PORT = "PORT"

const STR_templates_login_html = "templates/login.html"
const STR_templates_register_html = "templates/register.html"
const STR_templates_Content_html = "templates/Page1.html"
const STR_template_page_error_html = "templates/page_error.html"
const STR_template_result = "templates/result.json"
const STR_template_list_apikeys_html = "templates/list_apikeys.html"
const STR_template_list_reports_for_apikey_html = "templates/list_reports_for_apikey.html"
const STR_template_add_apikey_html = "templates/add_apikey.html"
const STR_template_apikey_deleteconfirm_html = "templates/apikey_deleteconfirm.html"
const STR_template_reports_deleteconfirm_html = "templates/reports_deleteconfirm.html"
const STR_template_list_clientids_for_apikey_html = "templates/list_clientids_for_apikey.html"
const STR_template_export_html = "templates/export.html"
const STR_template_invite_html = "templates/invite.html"

const STR_filepath_upload_template = "resources/reportfiles/%s"
const STR_filepath_download_template = "resources/exports/%s"
const STR_filepath_filelog = "resources/logs/logfile.txt"

const PATH_ROOT = "/"
const PATH_MESSAGE = "/message"
const PATH_CLIENTINFO = "/sendclientinfo"
const PATH_ClientInfoUpdate = "/updateclientinfo"
const PATH_UPLOADIMAGE = "/uploadimage"
const PATH_UPLOADFILE = "/uploadfile"
const PATH_STATIC_TEMPLATES = "./templates"
const PATH_Login = "/login"
const PATH_logout = "/logout"
const PATH_Register = "/register"
const PATH_ApiKeys = "/apikeys"
const PATH_Reports = "/reports"
const PATH_AddApiKey = "/apikeyadd"
const PATH_ApiKeyDeleteConfirm = "/apikeydeleteconfirm"
const PATH_ClientIds = "/clientids"
const PATH_ApiKeyDelete = "/apikeydelete"
const PATH_download = "/exports"
const PATH_filelog_delete = "/filelogdelete"
const PATH_reports_delete_confirm = "/reportsdeleteconfirm"
const PATH_reports_delete = "/reportsdelete"
const PATH_invite = "/invite"
const PORT_8080 = "8080"

const API_KEY_image = "image"
const API_KEY_file = "file"
const API_KEY_filename = "filename"
const API_KEY_sequence = "sequence"
const API_KEY_time = "time"
const API_KEY_message = "message"
const API_KEY_username = "username"
const API_KEY_password = "password"
const API_KEY_email = "email"
const API_KEY_token = "token"
const API_KEY_apikey = "apikey"
const API_KEY_startnum = "startnum"
const API_KEY_pagenum = "pagenum"
const API_KEY_pagesize = "pagesize"
const API_KEY_appname = "appname"
const API_KEY_clientid = "clientid"
const API_KEY_name = "name"
const API_KEY_inviteid = "inviteid"

const API_URL_domain = "localhost"
const API_URL = "http://localhost"

const DB_TYPE = "sqlite3"
const DB_NAME = "stunt.sqlite"
const TABLE_users = "users"
const TABLE_tokens = "tokens"
const TABLE_apikeys = "apikeys"
const TABLE_reports = "reports"
const TABLE_clientinfo = "clientinfo"
const TABLE_invites = "invites"
const TABLE_USERS_COLUMN_id = STR_id
const TABLE_USERS_COLUMN_email = "email"
const TABLE_USERS_COLUMN_password = "password"
const TABLE_USERS_COLUMN_salt = "salt"
const TABLE_TOKENS_COLUMN_id = STR_id
const TABLE_TOKENS_COLUMN_userid = "userid"
const TABLE_TOKENS_COLUMN_token = "token"
const TABLE_TOKENS_COLUMN_issued = "issued"
const TABLE_TOKENS_COLUMN_expires = "expires"
const TABLE_APIKEYS_COLUMN_userid = "userid"
const TABLE_APIKEYS_COLUMN_apikey = "apikey"
const TABLE_APIKEYS_COLUMN_appname = "appname"
const TABLE_REPORTS_COLUMN_clientid = "clientid"
const TABLE_REPORTS_COLUMN_time = "time"
const TABLE_REPORTS_COLUMN_sequence = "sequence"
const TABLE_REPORTS_COLUMN_message = "message"
const TABLE_REPORTS_COLUMN_filepath = "filepath"
const TABLE_REPORTS_COLUMN_id = STR_id
const TABLE_CLIENTINFO_clientid = "clientid"
const TABLE_CLIENTINFO_name = "name"
const TABLE_CLIENTINFO_manufacturer = "manufacturer"
const TABLE_CLIENTINFO_model = "model"
const TABLE_CLIENTINFO_deviceid = "deviceid"
const TABLE_INVITES_COLUMN_inviteid = "inviteid"
const TABLE_INVITES_COLUMN_apikey = TABLE_APIKEYS_COLUMN_apikey
const TABLE_INVITES_COLUMN_issued = TABLE_TOKENS_COLUMN_issued
const TABLE_INVITES_COLUMN_expires = TABLE_TOKENS_COLUMN_expires
const STMT_CREATE_TABLE_USERS = "create table if not exists users('id' integer primary key, 'email' text unique, 'password' text, 'salt' text)"
const STMT_CREATE_TABLE_TOKENS = "create table if not exists tokens('id' integer primary key, 'userid' integer, 'token' text, 'issued' integer, 'expires' integer)"
const STMT_CREATE_TABLE_APIKEYS = "create table if not exists apikeys('id' integer primary key, 'userid' integer, 'apikey' text, 'appname' text)"
const STMT_CREATE_TABLE_REPORTS = "create table if not exists reports%s('id' integer primary key, 'clientid' text, 'time' integer, 'sequence' integer, 'message' text, 'filepath' text)"
const STMT_CREATE_TABLE_CLIENTINFO = "create table if not exists clientinfo%s('clientid' text unique, 'name' text, 'manufacturer' text, 'model' text, 'deviceid' text)"
const STMT_CREATE_TABLE_INVITES = "create table if not exists invites%s('id' integer primary key, 'inviteid' text, 'apikey' text, 'issued' integer, 'expires' integer)"
const STMT_INSERT_INTO_USERS = "insert or ignore into users(email, password, salt) values(?, ?, ?)"
const STMT_INSERT_INTO_TOKENS = "insert or ignore into tokens(userid, token, issued, expires) values(?, ?, ?, ?)"
const STMT_INSERT_INTO_APIKEYS = "insert or ignore into apikeys(userid, apikey, appname) values(?, ?, ?)"
const STMT_INSERT_INTO_REPORTS = "insert or ignore into reports%s(clientid, time, sequence, message, filepath) values(?, ?, ?, ?, ?)"
const STMT_INSERT_INTO_CLIENTINFO = "insert or ignore into clientinfo%s(clientid, name, manufacturer, model, deviceid) values(?, ?, ?, ?, ?)"
const STMT_INSERT_INTO_INVITES = "insert or ignore into invites%s(inviteid, apikey, issued, expires) values(?, ?, ?, ?)"

//HashSha1 hashes the passed in string value using sha1 hash.
func HashSha1(aValue string) (string, error) {
	hashSha1 := sha1.New();
	_, err := hashSha1.Write([]byte(aValue))
	hashed := hashSha1.Sum(nil);
	return fmt.Sprintf("%x", hashed), err
}

//GenerateUUID creates a UUID.
func GenerateUUID() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}

//GenerateRandomString creates a random string with the passed in length.
func GenerateRandomString(aLength int) string {
	rand.Seed(time.Now().Unix());
	b := make([]rune, aLength)
	for i := range b {
		b[i] = Letters[rand.Intn(len(Letters))]
	}
	return string(b)
}

//GenerateToken creates an Authentication token to be used by the web application client.
func GenerateToken() (string, error) {
	tokenuuid, error := GenerateUUID()
	if error != nil {
		return STR_EMPTY, error
	}
	
	token := strings.Replace(tokenuuid, STR_symbol_dash, STR_EMPTY, -1)  
	return token, error
}

//AddCookie adds the passed in token to the Cookie header.
func AddCookie(responseWriter http.ResponseWriter, token string) {
	cookie := new(http.Cookie)
	cookie.Name = API_KEY_token
	cookie.Value = token
	cookie.Domain = API_URL_domain
	cookie.MaxAge = TOKEN_VALIDITY_SECONDS
	cookie.Path = "/"
	http.SetCookie(responseWriter, cookie)
}

//GetCookieToken gets the authentication token from the Cookie header.
func GetCookieToken(aRequest *http.Request) string {
//We don't want to use header, because we would have to write our AJAX application, that is why we use cookies
//	headers := aRequest.Header
//	tokens := headers["Token"]
//	var token = STR_EMPTY
//	if(tokens != nil && len(tokens) > 0) {
//		token = tokens[0]
//	}
//  return token
	
	cookie, errCookie := aRequest.Cookie(API_KEY_token)
	if errCookie != nil {
		log.Printf("GetHeaderToken, Error reading cookie, error=%s", errCookie.Error())
		return STR_EMPTY
	}
	token := cookie.Value
	return token
}

//IsTokenValid validates the token obtained from the Cookie header, against the database.
func IsTokenValid(responseWriter http.ResponseWriter, request *http.Request) bool {
	token := GetCookieToken(request)
	isValid, userId := DbIsTokenValid(token, nil)
	log.Printf("IsTokenValid, token=%s, isValid=%t, userId=%d", token, isValid, userId)
	if !isValid {
		ServeLogin(responseWriter, STR_MSG_login)
		return false
	}
	return true
}

//IsApiKeyValid validates the passed in API key against the database.
func IsApiKeyValid(aApiKey string) (isValid bool, userId int) {
	isValid, userId = DbIsApiKeyValid(aApiKey, nil)
	if !isValid {
		return false, -1
	}
	return isValid, userId
}

//FileLogCreate creates a log file to which logs can be written instead to standart output.
func FileLogCreate() {
	var err error
	FileLog, err = os.OpenFile("resources/logs/logfile.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Printf("main_stunt, init, error Open logfile, error=%s", err.Error())
		return
	} else {
		log.SetOutput(FileLog)
		
//		import "os/signal"
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, os.Interrupt)
//		go func(){
//    		for sig := range c {
//        		// sig is a ^C, handle it
//        		log.Printf("main_stunt, handle os.Interrupt, sig=%d", sig)
//        		utils.FileLog.Close()
//    		}
//		}()
	}
}

//GetApiUrlListApiKeys gets URL to PATH_ApiKeys.
func GetApiUrlListApiKeys() string {
	if Port == PORT_8080 {
		return API_URL + STR_symbol_colon + Port + PATH_ApiKeys
	} else {
		return API_URL + PATH_ApiKeys
	}
}

//GetApiUrlListClientIds gets URL to PATH_ClientIds.
func GetApiUrlListClientIds() string {
	if Port == PORT_8080 {
		return API_URL + STR_symbol_colon + Port + PATH_ClientIds
	} else {
		return API_URL + PATH_ClientIds
	}
}

//GetApiUrlInvite gets URL to PATH_ApiKeys.
func GetApiUrlInvite() string {
	if Port == PORT_8080 {
		return API_URL + STR_symbol_colon + Port + PATH_ApiKeys
	} else {
		return API_URL + PATH_ApiKeys
	}
}