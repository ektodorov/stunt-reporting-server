package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"objects"
	"time"
)

func DbInit() error {
	var err error = nil
	var db *sql.DB = nil
	var stmt *sql.Stmt = nil

	db, err = sql.Open(DB_TYPE, DB_NAME)
	if err != nil {
		log.Println("init, Error opening db", DB_NAME, ", err=", err)
		return err;
	}
	defer db.Close();
	
	//init table users
	stmt, err = db.Prepare(STMT_CREATE_TABLE_USERS)
	if err != nil {
		log.Println("init, Error preparing create table users stmt, err=", err)
		return err;
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("init, Error exec create table stmt, err=", err)
	}
	stmt.Close();
	
	//init table tokens
	stmt, err = db.Prepare(STMT_CREATE_TABLE_TOKENS)
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", STMT_CREATE_TABLE_TOKENS, err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error creating table, %s, error=%s", STMT_CREATE_TABLE_TOKENS, err.Error())
	}
	stmt.Close()
	
	//init table apikeys	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_APIKEYS)
	if err != nil {
		log.Println("init, Error preparing, %s, err=", STMT_CREATE_TABLE_APIKEYS, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("init, Error executing, %s, err=",STMT_CREATE_TABLE_APIKEYS, err)
	}
	stmt.Close();
	return nil
}


// Adds a user in the users table if it does not exist.
// ApiKey is also generated for the user. If apiKey creation fails for some reason,
// we would not fail the user creation. We would just have to check if a apiKey exists for that user when we are presenting it to the
// user and if it doesn't we should create it then.
func DbAddUser(aEmail string, aPassword string, aDb *sql.DB) (isUserExists bool, isUserAdded bool, errorUser error) {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return false, false, err
		}
		defer db.Close()
	}
	
	isUserExists = false
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email), err.Error())
		return
	}
	rows, err = stmt.Query(aEmail)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("select * from %s where %s=%s", TABLE_users, TABLE_USERS_COLUMN_email, aEmail), err.Error())
		return
	}
	for rows.Next() {
		var id int
		var email string
		var password string
		var salt string
		rows.Scan(&id, &email, &password, &salt)
		if email == aEmail {
			isUserExists = true
			break;
		}
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	if isUserExists == false {
		salt := GenerateRandomString(SALT_LENGTH)
		passwordHash, err := HashSha1(fmt.Sprintf("%s%s", aPassword, salt))
		if err != nil {
			log.Printf("Error hashing string=%s, error=%s", fmt.Sprintf("%s%s", aPassword, salt), err.Error())
			return isUserExists, false, err
		}
	
		stmt, err = db.Prepare(STMT_INSERT_INTO_USERS)
		if err != nil {
			log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_USERS, err.Error())
			return
		}
		_, err = stmt.Exec(aEmail, passwordHash, salt)
		if err != nil {
			log.Printf("Error executing %s, error=%s", STMT_INSERT_INTO_USERS, err.Error())
			return isUserExists, false, err		
		}
		stmt.Close()
		
		//create and insert apiKey
		stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where %s=?", TABLE_USERS_COLUMN_id, TABLE_users, TABLE_USERS_COLUMN_email))
		if err != nil {
			log.Printf("Error preparing, select id from users where email=?, error=%s", err.Error())
		}
		rows, err = stmt.Query(aEmail)
		if err != nil {
			log.Printf("Error query, select id from users where email=%s, error=%s", aEmail, err.Error())
		}
		var userId int = -1
		if rows.Next() {
			err = rows.Scan(&userId)
			if err != nil {
				log.Printf("Error scanning userId, error=%s", err.Error())
			}
		}
		if rows != nil {rows.Close()}
		if stmt != nil {stmt.Close()}
		
		if userId < 0 {
			return isUserExists, false, nil
		}
		//DbAddApiKey(userId, STR_EMPTY, db)
		
		return isUserExists, true, nil
	} else {
		return isUserExists, false, nil
	}
}

//Deletes a user and his apiKeys and the reports<apiKey> tables of that user
func DbDeleteUser(aUserId int, aDb *sql.DB) bool {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return false
		}
		defer db.Close()
	}
	
	//Delete user from table users
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_id))
	if err != nil {
		log.Printf("Error preparing, delete from users where id=?, error=%s", err.Error())
		return false
	}
	_, err = stmt.Exec(aUserId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("delete from %s where %s=%d", TABLE_users, TABLE_USERS_COLUMN_id, aUserId), err.Error())
		return false
	}
	stmt.Close()
	
	stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where %s=?", TABLE_APIKEYS_COLUMN_apikey, TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_APIKEYS_COLUMN_apikey, TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid), err.Error())
		return false
	}
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=%d", TABLE_APIKEYS_COLUMN_apikey, TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid, aUserId), err.Error())
		return false
	}	
	//get all apiKeys of user
	var sliceApiKeys []string = make([]string, 0, 16)
	for rows.Next() {
		var apiKey string
		rows.Scan(&apiKey)
		sliceApiKeys = append(sliceApiKeys, apiKey)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	if len(sliceApiKeys) == 0 {
		return true
	}
	
	//Delete apiKeys of user (delete from table apiKeys)
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s where %s=?", TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing, delete from %s where userid=?, error=%s", TABLE_apikeys, err.Error())
		return false
	}
	_, err = stmt.Exec(aUserId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("delete from %s where %s=%d", TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid, aUserId), err.Error())
		return false
	}
	stmt.Close()
	
	//Drop reports<apiKey> tables for the user
	for _, name := range sliceApiKeys {
		stmt, err = db.Prepare(fmt.Sprintf("drop table if exists %s%s", TABLE_reports, name))
		if err != nil {
			log.Printf("Error preparing, drop table if exists ?, error=%s", err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Printf("Error executing, drop table if exists %s, error=%s", (TABLE_reports + name), err.Error())
		}
	}
	if stmt != nil {stmt.Close()}
	
	return true
}

func DbGetUser(aEmail string, aPassword string, aDb *sql.DB) (id int, err error){
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var row *sql.Row = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return -1, err
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email), err.Error())
	}
	row = stmt.QueryRow(aEmail)
	
	var email string
	var password string
	var salt string
	err = row.Scan(&id, &email, &password, &salt)
	if err != nil && err == sql.ErrNoRows {
		return -1, err
	}
	passwordHash, err := HashSha1(fmt.Sprintf("%s%s", aPassword, salt))
	if err != nil {
		log.Printf("Error hashing string=%s, error=%s", fmt.Sprintf("%s%s", aPassword, salt), err.Error())
	}
	if passwordHash != password {
		return -1, err
	}
	
	return id, err
}

func DbGetUserLoad(aUserId int, aDb *sql.DB) (user *objects.User, err error) {
	user = new(objects.User)
	user.Id = aUserId
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var row *sql.Row = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return user, err
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_id))
	if stmt != nil {defer stmt.Close()}
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email), err.Error())
	}
	row = stmt.QueryRow(aUserId)
	
	var id int
	var email string
	var password string
	var salt string
	err = row.Scan(&id, &email, &password, &salt)
	if err != nil && err == sql.ErrNoRows {
		return user, err
	} else {
		user.Id = id
		user.Email = email
	}
	
	return user, err
}

func DbAddToken(aUserId int, aDb *sql.DB) (token string) {
	DbCleanTokens(aUserId, aDb)

	var err error = nil
	var db *sql.DB = aDb
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("AddToken, Error opening db, err=", err)
			return
		}
		defer db.Close()
	}
	
	token, errUUID := GenerateUUID()
	if errUUID != nil {
		fmt.Println("Erro generating uuid, err=", errUUID)
	}
	issued := time.Now().UnixNano() / int64(time.Millisecond)
	expires := TOKEN_VALIDITY_MS
	
	_, err = db.Exec("insert or ignore into tokens(userid, token, issued, expires) values(?, ?, ?, ?)", aUserId, token, issued, expires)
	if err != nil {
		fmt.Println("AddToken, Error inserting into tokens, err=", err)
	}
	return token
}

func DbDeleteApiKey(aApiKey string, aDb *sql.DB) (isDeleted bool) {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("DbDeleteApiKey, Error opening db login.sqlite, err=%s", err.Error())
			return
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s where %s=?", TABLE_apikeys, TABLE_APIKEYS_COLUMN_apikey))
	if err != nil {
		log.Printf("Error preparing, delete from %s where %s=?, error=%s", TABLE_apikeys, TABLE_APIKEYS_COLUMN_apikey, err.Error())
		return false
	}
	_, err = stmt.Exec(aApiKey)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("delete from %s where %s=%d", TABLE_apikeys, TABLE_APIKEYS_COLUMN_apikey, aApiKey), err.Error())
		return false
	}
	stmt.Close()
	return true
}

func DbIsTokenValid(aToken string, aDb *sql.DB) (isValid bool, userId int) {
	var err error = nil
	var db *sql.DB = aDb
	var rows *sql.Rows = nil

	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("IsTokenValid, Error opening db login.sqlite, err=", err)
			return false, -1
		}
		defer db.Close()
	}
	
	rows, err = db.Query("select * from tokens")
	if err != nil {
		fmt.Println("IsTokenValid, Error select from tokens, err=", err)
		return false, -1;
	}
	defer rows.Close()
	
	now := time.Now().UnixNano() / int64(time.Millisecond)
	for (rows.Next()) {
		var id int
		var userId int
		var token string
		var issued int64
		var expires int64
		err = rows.Scan(&id, &userId, &token, &issued, &expires)
		if err != nil {
			fmt.Println("IsTokenValid, Error scan tokens, err=", err)
		}
		if token == aToken && now < (issued + expires) {
			return true, userId;
		}
	}
	return false, -1;
}

func DbIsApiKeyValid(aApiKey string, aDb *sql.DB) (isValid bool, userId int) {
	var err error = nil
	var db *sql.DB = aDb
	var rows *sql.Rows = nil

	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("IsTokenValid, Error opening db login.sqlite, err=%s", err.Error())
			return false, -1
		}
		defer db.Close()
	}
	
	rows, err = db.Query(fmt.Sprintf("select * from %s", TABLE_apikeys))
	if err != nil {
		log.Printf("IsTokenValid, Error select from %s, err=%s", TABLE_apikeys, err.Error())
		return false, -1;
	}
	defer rows.Close()
	for (rows.Next()) {
		var userId int
		var apiKey string
		var appName string
		err = rows.Scan(&userId, &apiKey, &appName)
		if err != nil {
			log.Printf("IsApiKeyValid, Error scan %s, err=%s", TABLE_apikeys, err.Error())
		}
		if apiKey == aApiKey {
			return true, userId;
		}
	}
	
	return false, -1
}

func DbCleanTokens(aUserId int, aDb *sql.DB) {
	var err error = nil
	var db *sql.DB = aDb
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("IsTokenValid, Error opening db login.sqlite, err=", err)
			return
		}
		defer db.Close()
	}
	
	now := time.Now().UnixNano() / int64(time.Millisecond)
	db.Exec("delete from tokens where userid=? AND issued + expires < ?", aUserId, now)
}

// Get all apiKeys of a user.
func DbGetApiKey(aUserId int, aDb *sql.DB) []*objects.ApiKey {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return nil
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s where %s=?", TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_APIKEYS_COLUMN_apikey, TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid), err.Error())
	}
	
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_APIKEYS_COLUMN_apikey, TABLE_apikeys, TABLE_APIKEYS_COLUMN_userid), err.Error())
	}
	var sliceApiKeys []*objects.ApiKey = make([]*objects.ApiKey, 0, 16)
	for rows.Next() {
		var userId int
		var apiKey string
		var appName string
		var objApiKey  *objects.ApiKey
		objApiKey = new(objects.ApiKey)
		rows.Scan(&userId, &apiKey, &appName)
		objApiKey.UserId = userId
		objApiKey.ApiKey = apiKey
		objApiKey.AppName = appName
		sliceApiKeys = append(sliceApiKeys, objApiKey)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceApiKeys
}

// ApiKey is added in DbAddUser.
// This method we can use when we want to add additional apiKeys for a user, or if the user does not have a apiKey when we present it to him.
func DbAddApiKey(aUserId int, aAppName string, aDb *sql.DB) bool {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil	
	
	if aUserId < 0 {
		return false
	}
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return false
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(STMT_INSERT_INTO_APIKEYS)
	if err != nil {
		log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_APIKEYS, err.Error())
	}
	
	var apiKey string = STR_EMPTY
	apiKey, err = GenerateToken()
	if err != nil {
		log.Printf("Error generateToken, error=%s", err.Error())
	}
	if aAppName == STR_EMPTY {
		aAppName = apiKey
	}
	_, err = stmt.Exec(aUserId, apiKey, aAppName)
	if err != nil {
		log.Printf("Error executing %s, values userId=%d, apiKey=%s, error=%s", STMT_INSERT_INTO_APIKEYS, aUserId, apiKey, err.Error())
		if stmt != nil {stmt.Close()}
		return false
	}
	if stmt != nil {stmt.Close()}
	return true
}

func DbAddReport(aApiKey string, aClientId string, aTime int, aSequence int, aMessage string, aFilePath string, aDb *sql.DB) {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return
		}
		defer db.Close()
	}	
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aApiKey))
	if err != nil {
		log.Printf("Error creating, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aApiKey), err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aApiKey), err.Error())
	}
	if stmt != nil {stmt.Close()}
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aApiKey))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aApiKey), err.Error())
	}
	_, err = stmt.Exec(aClientId, aTime, aSequence, aMessage, aFilePath)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aApiKey), err.Error())
	}
	if stmt != nil {stmt.Close()}
}


func DbDeleteReport(aApiKey string, aId int, aDb *sql.DB) {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s where %s=?", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf("delete from %s%s where %s=?", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id), err.Error())
	}
	_, err = stmt.Exec(aId)
	if err != nil {
		log.Printf("Error deleting, %s, error=%s", 
				fmt.Sprintf("delete from %s%s where %s=aId", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, aId), err.Error())
	}
	if stmt != nil {stmt.Close()}
}

// Delete all records in the reports<apiKey> table
func DbClearReports(aApiKey string, aDb *sql.DB) {
	// Instead of deleting all from reports<aApiKey> we can just - drop table if exists reports<aApiKey>
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s", TABLE_reports, aApiKey))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aApiKey), err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aApiKey), err.Error())
	}
	stmt.Close()
}

func DbGetReportsByApiKey(aApiKey string, aClientId string, aStartNum int, aPageSize int, aDb *sql.DB) (sliceReports []*objects.Report, endNum int) {
	endNum = aStartNum
	sliceReports = make([]*objects.Report, 0, 64)
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s > ? order by %s, %s limit ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s > ? order by %s limit ?", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, endNum
	}
	
	rows, err = stmt.Query(aStartNum, (aStartNum + aPageSize))
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s > ? order by %s limit ?", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, endNum
	}
	
	for rows.Next() {
		var id int
		var clientId string
		var time int
		var sequence int
		var message string
		var filePath string
		err = rows.Scan(&id, &clientId, &time, &sequence, &message, &filePath)
		if err != nil {
			log.Printf("Error scanning, error=%s", err.Error())
		}
		var report = new(objects.Report)
		report.Id = id
		report.ClientId = clientId
		report.Time = time
		report.Message = message
		report.FilePath = filePath
		sliceReports = append(sliceReports, report)
		endNum++
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceReports, endNum
}
