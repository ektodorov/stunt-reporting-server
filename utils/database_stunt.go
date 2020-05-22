package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"objects"
	"time"
)

//DbInit initalizes the database - creates tables used by the application.
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
	if stmt != nil {stmt.Close();}
	
	//init table tokens
	stmt, err = db.Prepare(STMT_CREATE_TABLE_TOKENS)
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", STMT_CREATE_TABLE_TOKENS, err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error creating table, %s, error=%s", STMT_CREATE_TABLE_TOKENS, err.Error())
	}
	if stmt != nil {stmt.Close()}
	
	//init table apikeys	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_APIKEYS)
	if err != nil {
		log.Printf("init, Error preparing, %s, err=%s", STMT_CREATE_TABLE_APIKEYS, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("init, Error executing, %s, err=%s",STMT_CREATE_TABLE_APIKEYS, err)
	}
	if stmt != nil {stmt.Close();}
	
	return nil
}


//DbAddUser adds a user in the users table if it does not exist.
//ApiKey is also generated for the user.
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

//Deletes a user and his apiKeys and the reports<apiKey> tables of that user.
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

//DbGetUser gets a user record from the users table.
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

//DbGetUserLoad gets a user record from the users table.
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

//DbAddToken adds a token for a user to the tokens table.
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

//DbDeleteApiKey deletes API key from api keys table.
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

//DbIsTokenValid checks if a token is valid. The check is against the TABLE_tokens table.
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
	
	rows, err = db.Query(fmt.Sprintf("select * from %s", TABLE_tokens))
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

//DbIsApiKeyValid checks if an API Key is valid. The check is against the TABLE_apikeys table.
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
		var id int
		var userId int
		var apiKey string
		var appName string
		err = rows.Scan(&id, &userId, &apiKey, &appName)
		if err != nil {
			log.Printf("IsApiKeyValid, Error scan %s, err=%s", TABLE_apikeys, err.Error())
		}
		if apiKey == aApiKey {
			return true, userId;
		}
	}
	
	return false, -1
}

//DbCleanTokens deletes expired token records from the TABLE_tokens table.
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
	db.Exec(fmt.Sprintf("delete from %s where %s=? AND %s + %s < ?", 
		TABLE_tokens, TABLE_TOKENS_COLUMN_userid, TABLE_TOKENS_COLUMN_issued, TABLE_TOKENS_COLUMN_expires), 
		aUserId, now)
}

//DbGetApiKey gets all apiKeys of a user.
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
		var id int
		var userId int
		var apiKey string
		var appName string
		var objApiKey  *objects.ApiKey
		objApiKey = new(objects.ApiKey)
		rows.Scan(&id, &userId, &apiKey, &appName)
		objApiKey.UserId = userId
		objApiKey.ApiKey = apiKey
		objApiKey.AppName = appName
		sliceApiKeys = append(sliceApiKeys, objApiKey)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceApiKeys
}

//DbAddApiKey adds additional API Keys for a user. Initial API Key is added in the DbAddUser.
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
		log.Printf("Error executing %s, values userId=%d, apiKey=%s, appName=%s, error=%s", STMT_INSERT_INTO_APIKEYS, aUserId, apiKey, aAppName, err.Error())
		if stmt != nil {stmt.Close()}
		return false
	}
	if stmt != nil {stmt.Close()}
	return true
}

//DbAddClientInfo adds client info to the TABLE_clientinfo for API key table.
func DbAddClientInfo(aApiKey string, aClientId string, aName string, aManufacturer string, aModel string, aDeviceId string, aDb *sql.DB) error {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return err
		}
		defer db.Close()
	}	
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_CREATE_TABLE_CLIENTINFO, aApiKey))
	if err != nil {
		log.Printf("Error creating, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_CLIENTINFO, aApiKey), err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_CLIENTINFO, aApiKey), err.Error())
		return err
	}
	if stmt != nil {stmt.Close()}
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_INSERT_INTO_CLIENTINFO, aApiKey))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_CLIENTINFO, aApiKey), err.Error())
		return err
	}
	_, err = stmt.Exec(aClientId, aName, aManufacturer, aModel, aDeviceId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_CLIENTINFO, aApiKey), err.Error())
		return err
	}
	if stmt != nil {stmt.Close()}
	return nil
}

//DbDeleteClientInfo deletes client info from the TABLE_clientinfo for API Key table.
func DbDeleteClientInfo(aApiKey string, aClientId string, aDb *sql.DB) (isDeleted bool) {
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
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid))
	if err != nil {
		log.Printf("Error preparing, %s", fmt.Sprintf("delete from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid), err.Error())
		return false
	}
	_, err = stmt.Exec(aClientId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("delete from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid), err.Error())
		return false
	}
	if stmt != nil {stmt.Close()}
	return true
}

//DbGetClientInfo gets client info from the TABLE_clientinfo for an API Key.
func DbGetClientInfo(aApiKey string, aClientId string, aDb *sql.DB) *objects.ClientInfo {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var row *sql.Row = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("DbDeleteApiKey, Error opening db login.sqlite, err=%s", err.Error())
			return nil
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid))
	if err != nil {
		log.Printf("Error preparing, %s", fmt.Sprintf("select * from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid), err.Error())
		return nil
	}
	row = stmt.QueryRow(aClientId)
	var clientId string
	var name string
	var manufacturer string
	var model string
	var deviceId string
	err = row.Scan(&clientId, &name, &manufacturer, &model, &deviceId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("select * from %s%s where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_clientid), err.Error())
		return nil
	}
	clientInfo := new(objects.ClientInfo)
	clientInfo.ClientId = clientId
	clientInfo.Name = name
	clientInfo.Manufacturer = manufacturer
	clientInfo.Model = model
	clientInfo.DeviceId = deviceId
	
	if stmt != nil {stmt.Close()}
	
	return clientInfo
}

//DbGetClientInfos gets all client info records from TABLE_clientinfo for API key table.
func DbGetClientInfos(aApiKey string, aDb *sql.DB) []*objects.ClientInfo {
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
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s", TABLE_clientinfo, aApiKey))
	if err != nil {
		log.Printf("DbGetClientInfos, Error preparing %s, error=%s", fmt.Sprintf("select * from %s%s", TABLE_clientinfo, aApiKey), err.Error())
		return nil
	}
	
	rows, err = stmt.Query()
	if err != nil {
		log.Printf("DbGetClientInfos, Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s", TABLE_clientinfo, aApiKey), err.Error())
	}
	var sliceClientInfo []*objects.ClientInfo = make([]*objects.ClientInfo, 0, 16)
	for rows.Next() {
		var clientId string
		var name string
		var manufacturer string
		var model string
		var deviceId string
		var objClientInfo  *objects.ClientInfo
		objClientInfo = new(objects.ClientInfo)
		rows.Scan(&clientId, &name, &manufacturer, &model, &deviceId)
		objClientInfo.ApiKey = aApiKey
		objClientInfo.ClientId = clientId
		objClientInfo.Name = name
		objClientInfo.Manufacturer = manufacturer
		objClientInfo.Model = model
		objClientInfo.DeviceId = deviceId
		sliceClientInfo = append(sliceClientInfo, objClientInfo)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceClientInfo
}

//DbUpdateClientInfo updates client info in TABLE_clientinfo for API key table.
func DbUpdateClientInfo(aApiKey string, aClientId string, aName string, aDb *sql.DB) error {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("DbUpdateClientInfo, Error opening db login.sqlite, err=%s", err.Error())
			return err
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("update %s%s set %s=? where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_name, TABLE_CLIENTINFO_clientid))
	if err != nil {
		log.Printf("DbUpdateClientInfo, Error preparing, %s", 
			fmt.Sprintf("update %s%s set %s=? where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_name, TABLE_CLIENTINFO_clientid), err.Error())
		return err
	}
	
	_, err = stmt.Exec(aName, aClientId)
	if err != nil {
		log.Printf("DbUpdateClientInfo, Error executing, %s", 
			fmt.Sprintf("update %s%s set %s=? where %s=?", TABLE_clientinfo, aApiKey, TABLE_CLIENTINFO_name, TABLE_CLIENTINFO_clientid), err.Error())
	}
	return err
}

//DbClearClientInfo deletes all client info records from TABLE_clientinfo for API key table.
func DbClearClientInfo(aApiKey string, aDb *sql.DB) error {
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var row *sql.Row = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return nil
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("select name from sqlite_master where type='table' and name='%s%s'", TABLE_clientinfo, aApiKey))
	if err != nil {
		log.Printf("DbClearClientInfo, Error preparing %s, error=%s", 
			fmt.Sprintf("select name from sqlite_master where type='table' and name='%s'", TABLE_clientinfo, aApiKey), err.Error())
		return nil
	}
	row = stmt.QueryRow()
	if row != nil {
		var name string
		err = row.Scan(&name)
		if err != nil {
			stmt.Close()
			return err
		} else {
			stmt.Close()
		}
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s", TABLE_clientinfo, aApiKey))
	if err != nil {
		log.Printf("DbClearClientInfo, Error preparing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_clientinfo, aApiKey), err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("DbClearClientInfo, Error executing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_clientinfo, aApiKey), err.Error())
	}
	if stmt != nil {stmt.Close()}
	return err
}

//DbAddReport adds report to TABLE_reports for API key table.
func DbAddReport(aApiKey string, aClientId string, aTime int64, aSequence int, aMessage string, aFilePath string, aDb *sql.DB) {
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

//DbDeleteReport deletes report from TABLE_reports for API key table.
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

//DbClearReports deletes all records in the TABLE_reports for API key table.
func DbClearReports(aApiKey string, aDb *sql.DB) error {
	// Instead of deleting all from reports<aApiKey> we can just - drop table if exists reports<aApiKey>
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return err
		}
		defer db.Close()
	}
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s", TABLE_reports, aApiKey))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aApiKey), err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aApiKey), err.Error())
	}
	if stmt != nil {stmt.Close()}
	return err
}

//DbGetReportsByApiKey deletes reports from TABLE_reports for API key. All records will be deleted if no clientId is supplied, otherwise only 
//the records for the supplied clientId will be deleted.
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
	
	if aClientId != STR_EMPTY {
		stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s=? order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
		log.Printf("DbGetReportsByApiKey, %s", fmt.Sprintf("select * from %s%s where %s=? order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
	} else {
		stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
		log.Printf("DbGetReportsByApiKey, %s", fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
	}
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, endNum
	}
	
	if aClientId != STR_EMPTY {
		log.Printf("DbGetReportsByApiKey, Query, aClientId=%s", aClientId)
		rows, err = stmt.Query(aClientId, aStartNum, aPageSize)
	} else {
		log.Printf("DbGetReportsByApiKey, Query, without clientId, aClientId=%s", aClientId)
		rows, err = stmt.Query(aStartNum, aPageSize)
	}
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, endNum
	}
	
	for rows.Next() {
		var id int
		var clientId string
		var reportTime int64
		var sequence int
		var message string
		var filePath string
		err = rows.Scan(&id, &clientId, &reportTime, &sequence, &message, &filePath)
		if err != nil {
			log.Printf("Error scanning, error=%s", err.Error())
		}
		var report = new(objects.Report)
		report.Id = id
		report.ClientId = clientId
		report.Time = reportTime
		report.Sequence = sequence
		report.Message = message
		report.FilePath = filePath
		report.TimeString = fmt.Sprintf("%s", time.Unix(reportTime, 0))
		sliceReports = append(sliceReports, report)
		endNum++
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceReports, endNum
}

//DbGetReports gets all reports from TABLE_reports for API key table.
func DbGetReports(aApiKey string, aId int, aPageSize int, aDb *sql.DB) (sliceReports []*objects.Report) {
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
			fmt.Sprintf("select * from %s%s where %s > ? order by %s, %s limit ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports
	}
	
	rows, err = stmt.Query(aId, aPageSize)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s > ? order by %s, %s limit ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports
	}
	
	for rows.Next() {
		var id int
		var clientId string
		var time int64
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
		report.Sequence = sequence
		report.Message = message
		report.FilePath = filePath
		sliceReports = append(sliceReports, report)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceReports
}

//DbGetReportsLastPage gets the last records from the TABLE_reports for API key, that are in the last page according to the pagination.
func DbGetReportsLastPage(aApiKey string, aClientId string, aPageSize int, aDb *sql.DB) (sliceReports []*objects.Report, count int64) {
	sliceReports = make([]*objects.Report, 0, 64)
	var err error = nil
	var db *sql.DB = aDb
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	var rowCount int64 = 0
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
			return
		}
		defer db.Close()
	}
	
	if aClientId != STR_EMPTY {
		stmt, err = db.Prepare(fmt.Sprintf("select Count(*) from %s%s where %s=?", TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid))
	} else {
		stmt, err = db.Prepare(fmt.Sprintf("select Count(*) from %s%s",TABLE_reports, aApiKey))
	}
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select Count(*) from %s%s",TABLE_reports, aApiKey), err.Error())
		return sliceReports, rowCount
	}
	
	if aClientId != STR_EMPTY {
		rows, err = stmt.Query(aClientId)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		log.Printf("Error quering, %s, error=%s", fmt.Sprintf("select Count(*) from %s%s", TABLE_reports, aApiKey), err.Error())
		return sliceReports, rowCount
	}
	for rows.Next() {
		err = rows.Scan(&rowCount)
		if err != nil {
			log.Printf("Error scanning, error=%s", err.Error())
		}
	}
	log.Printf("DbGetReportsLastPage, rowCount=%d", rowCount)
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	if aClientId != STR_EMPTY {
		stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s=? order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
	} else {
		stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id))
	}
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, rowCount
	}
	
	if aClientId != STR_EMPTY {
		rows, err = stmt.Query(aClientId, (rowCount - int64(aPageSize)), aPageSize)
	} else {
		rows, err = stmt.Query((rowCount - int64(aPageSize)), aPageSize)
	}
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s order by %s, %s limit ?, ?",
					TABLE_reports, aApiKey, TABLE_REPORTS_COLUMN_clientid, TABLE_REPORTS_COLUMN_id),
			err.Error())
		return sliceReports, rowCount
	}
	
	for rows.Next() {
		var id int
		var clientId string
		var reportTime int64
		var sequence int
		var message string
		var filePath string
		err = rows.Scan(&id, &clientId, &reportTime, &sequence, &message, &filePath)
		if err != nil {
			log.Printf("Error scanning, error=%s", err.Error())
		}
		var report = new(objects.Report)
		report.Id = id
		report.ClientId = clientId
		report.Time = reportTime
		report.Sequence = sequence
		report.Message = message
		report.FilePath = filePath
		report.TimeString = fmt.Sprintf("%s", time.Unix(reportTime, 0))
		sliceReports = append(sliceReports, report)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceReports, rowCount
}

//DbInviteAddApiKey adds an API key to a user that has been invited. 
//The TABLE_invites for API key is queried to check if the invitation has expired and if the it exists.
func DbInviteAddApiKey(aUserId int, aInviteId string, aApiKey string, aAppName string, aDb *sql.DB) {
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
	
	DbInviteCreateTable(aApiKey, db)
	
	var isValidInvite = false
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s=?", TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_inviteid))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s=?", TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_inviteid), err.Error())
		return
	}
	rows, err = stmt.Query(aInviteId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", fmt.Sprintf("select * from %s%s where %s=?", TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_inviteid), err.Error())
		return
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	for rows.Next() {
		var id string
		var inviteId string
		var apiKey string
		var issued int64
		var expires int64
		err = rows.Scan(&id, &inviteId, &apiKey, &issued, &expires)
		if err != nil {
			log.Printf("Error scanning, error=%s", err.Error())
			return
		}
		
		if apiKey == aApiKey && now < (issued + expires) {
			isValidInvite = true
			break
		}
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	if !isValidInvite {return}
	
	stmt, err = db.Prepare(STMT_INSERT_INTO_APIKEYS)
	if err != nil {
		log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_APIKEYS, err.Error())
		return
	}
	_, err = stmt.Exec(aUserId, aApiKey, aAppName)
	if err != nil {
		log.Printf("Error executing %s, values userId=%d, apiKey=%s, appName=%s, error=%s", STMT_INSERT_INTO_APIKEYS, aUserId, aApiKey, aAppName, err.Error())
		return
	}
	if stmt != nil {stmt.Close()}
}

//DbInviteAdd adds an invitation to TABLE_invites for API key.
func DbInviteAdd(aApiKey string, aDb *sql.DB) (inviteId string) {
	var err error = nil
	var stmt *sql.Stmt = nil
	var db *sql.DB = aDb
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("DbInviteAdd, Error opening db, err=", err)
			return
		}
		defer db.Close()
	}
	
	DbInviteCreateTable(aApiKey, db)
	DbInviteClean(aApiKey, db)
	
	inviteId, errUUID := GenerateUUID()
	if errUUID != nil {
		fmt.Println("Error generating uuid, err=", errUUID)
	}
	log.Printf("DbInviteAdd, inviteId=%s", inviteId)
	issued := time.Now().UnixNano() / int64(time.Millisecond)
	expires := INVITE_VALIDITY_MS
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_INSERT_INTO_INVITES, aApiKey))
	_, err = stmt.Exec(inviteId, aApiKey, issued, expires)
	if err != nil {
		fmt.Println("DbInviteAdd, Error inserting into tokens, err=", err)
	}
	if stmt != nil {stmt.Close()}
	
	return inviteId
}

//DbInviteClean deletes expired invitation records from TABLE_invites for API key table.
func DbInviteClean(aApiKey string, aDb *sql.DB) {
	var err error = nil
	var stmt *sql.Stmt = nil
	var db *sql.DB = aDb
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("DbInviteClean, Error opening db login.sqlite, err=", err)
			return
		}
		defer db.Close()
	}
	
	now := time.Now().UnixNano() / int64(time.Millisecond)
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s where %s=? AND %s + %s < ?", 
				TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_apikey, TABLE_INVITES_COLUMN_issued, TABLE_INVITES_COLUMN_expires))
	if err != nil {
		log.Printf("DbInviteClean, error preparing %s, error=%s", fmt.Sprintf("delete from %s%s where %s=? AND %s + %s < ?", 
				TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_apikey, TABLE_INVITES_COLUMN_issued, TABLE_INVITES_COLUMN_expires), err.Error())
		return
	}
	_, err = stmt.Exec(aApiKey, now)
	if err != nil {
		log.Printf("DbInviteClean, error executing %s, error=%s", fmt.Sprintf("delete from %s%s where %s=%s AND isssued + expires < %d", 
				TABLE_invites, aApiKey, TABLE_INVITES_COLUMN_apikey, aApiKey, now), err.Error())
	}
}

//DbInviteCreateTable creates TABLE_invites for API key table.
func DbInviteCreateTable(aApiKey string, aDb *sql.DB) {
	var err error
	var stmt *sql.Stmt = nil
	var db *sql.DB = aDb
	
	if db == nil {
		db, err = sql.Open(DB_TYPE, DB_NAME)
		if err != nil {
			fmt.Println("DbInviteClean, Error opening db login.sqlite, err=", err)
			return
		}
		defer db.Close()
	}
	
	//init table invites
	stmt, err = db.Prepare(fmt.Sprintf(STMT_CREATE_TABLE_INVITES, aApiKey))
	if err != nil {
		log.Println("init, Error preparing, %s, err=%s", STMT_CREATE_TABLE_INVITES, err)
	}
	if stmt != nil {defer stmt.Close()}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("init, Error executing, %s, err=%s",STMT_CREATE_TABLE_INVITES, err)
	}
}