package utils

import (
	"crypto/sha1"
	"github.com/nu7hatch/gouuid"
	"fmt"
	"math/rand"
	"time"
)

const STR_EMPTY = ""
const STR_BLANK = " "
const STR_MSG_NOTFOUND = "404 Not found"
const STR_GET = "GET"
const STR_POST = "POST"
const STR_error = "error"
const STR_Authorization = "Authorization"

const STR_template_page_error_html = "templates/page_error.html"
const STR_template_result = "templates/result.json"

const STR_img_filepathSave_template = "templates/images/%s"
const STR_img_filepathSrc_template = "images/%s"
const STR_img_filepathTemplates_template = "templates/%s"

var Letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const PATH_ROOT = "/"
const PATH_ECHO = "/echo"
const PATH_MESSAGE = "/message"
const PATH_UPLOADIMAGE = "/uploadimage"
const PATH_UPLOADFILE = "/uploadfile"
const PATH_STATIC_TEMPLATES = "./templates"
const PORT_8080 = ":8080"

const API_KEY_image = "image"
const API_KEY_file = "file"
const API_KEY_filename = "filename"
const API_KEY_sequence = "sequence"
const API_KEY_time = "time"
const API_KEY_message = "message"

const DB_TYPE = "sqlite3"
const DB_NAME = "stunt.sqlite"
const TABLE_USERS_COLUMN_id = "id"
const TABLE_USERS_COLUMN_email = "email"
const TABLE_USERS_COLUMN_password = "password"
const TABLE_USERS_COLUMN_name = "name"
const TABLE_USERS_COLUMN_salt = "salt"
const TABLE_TOKENS_COLUMN_userid = "userid"
const TABLE_TOKENS_COLUMN_token = "token"
const TABLE_REPORTS_COLUMN_clientid = "clientid"
const TABLE_REPORTS_COLUMN_time = "time"
const TABLE_REPORTS_COLUMN_sequence = "sequence"
const TABLE_REPORTS_COLUMN_message = "message"
const TABLE_REPORTS_COLUMN_filepath = "filepath"
const STMT_CREATE_TABLE_USERS = "create table if not exists users('id' integer primary key, 'email' text unique, 'password' text, 'name' text, 'salt' text)"
const STMT_CREATE_TABLE_TOKENS = "create table if not exists tokens('userid' integer, 'token' text unique)"
const STMT_CREATE_TABLE_REPORTS = "create table if not exists reports%s('clientid' string, 'time' integer, 'sequence' integer, 'message' text, 'filepath' text)"
const STMT_CREATE_TABLE_REPORTS_TABLES = "create table if not exists reports('id' integer primary key, 'tablename' string)"
const STMT_INSERT_INTO_USERS = "insert or ignore into users(email, password, name, salt) values(?, ?, ?, ?)"
const STMT_INSERT_INTO_TOKENS = "insert or ignore into tokens(userid, token) values(?, ?)"
const STMT_INSERT_INTO_REPORTS = "insert or ignore into reports%s(clientid, time, sequence, message, filepath) values(?, ?, ?, ?, ?)"
const STMT_INSERT_INTO_REPORTS_TABLES = "insert or ignore into reports(id, tablename) values(?, ?)"

func HashSha1(aValue string) (string, error) {
	hashSha1 := sha1.New();
	_, err := hashSha1.Write([]byte(aValue))
	hashed := hashSha1.Sum(nil);
	return fmt.Sprintf("%x", hashed), err
}

func GenerateUUID() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}

func GenerateRandomString(aLength int) string {
	rand.Seed(time.Now().Unix());
	b := make([]rune, aLength)
	for i := range b {
		b[i] = Letters[rand.Intn(len(Letters))]
	}
	return string(b)
}