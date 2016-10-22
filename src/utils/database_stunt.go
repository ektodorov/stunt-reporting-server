package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
//	"time"
)

func DbInit() {
	var err error = nil;
	var db *sql.DB = nil;
	var stmt *sql.Stmt = nil;

	db, err = sql.Open(DB_TYPE, DB_NAME);
	if err != nil {
		log.Println("init, Error opening db", DB_NAME, ", err=", err);
		return;
	}
	defer db.Close();
	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_USERS);
	if err != nil {
		log.Println("init, Error preparing create table users stmt, err=", err);
		return;
	}
	_, err = stmt.Exec();
	if err != nil {
		log.Println("init, Error exec create table stmt, err=", err);
	}
	stmt.Close();
	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_TOKENS);
	if err != nil {
		log.Println("init, Error preparing create table tokens, err=", err);
	}
	_, err = stmt.Exec();
	if err != nil {
		log.Println("init, Error executing create table tokens, err=", err);
	}
	stmt.Close();
	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_REPORTS_TABLES);
	if err != nil {
		log.Println("init, Error preparing create table reports, err=", err);
	}
	_, err = stmt.Exec();
	if err != nil {
		log.Println("init, Error executing create table reports, err=", err);
	}
	stmt.Close();
}


// Adds a user in the users table if it does not exist.
// Token is also generated for the user. If token creation fails for some reason,
// we would not fail the user creation. We would just have to check if a token exists for that user when we are presenting it to the
// user and if it doesn't we should create it then.
func DbAddUser(aEmail string, aPassword string) (isUserExists bool, isUserAdded bool) {
	var err error = nil;
	var db *sql.DB = nil;
	var stmt *sql.Stmt = nil;
	var rows *sql.Rows = nil;
	
	db, err = sql.Open(DB_TYPE, DB_NAME)
	if err != nil {
		log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
		return
	}
	defer db.Close()
	
	isUserExists = false
	stmt, err = db.Prepare("select * from users where email=?")
	if err != nil {
		log.Printf("Error preparing 'select * from users where email=?', error=%s", err.Error())
		return
	}
	rows, err = stmt.Query(aEmail)
	if err != nil {
		log.Printf("Error executing 'select * from users where email=%s', error=%s", aEmail, err.Error())
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
			return isUserExists, false
		}
	
		stmt, err = db.Prepare(STMT_INSERT_INTO_USERS)
		if err != nil {
			log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_USERS, err.Error())
			return
		}
		_, err = stmt.Exec(aEmail, passwordHash, salt)
		if err != nil {
			log.Printf("Error executing %s, error=%s", STMT_INSERT_INTO_USERS, err.Error())
			return isUserExists, false		
		}
		stmt.Close()
		
		//create and insert token
		stmt, err = db.Prepare("select id from users where email=?")
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
			return isUserExists, true
		}
		DbAddToken(userId, db)
		
		return isUserExists, true
	} else {
		return isUserExists, false
	}
}

//Deletes a user and his tokens and the reports<token> tables of that user
func DbDeleteUser(aUserId int) bool {
	var err error = nil
	var db *sql.DB = nil
	var stmt *sql.Stmt = nil
	var rows *sql.Rows = nil
	
	db, err = sql.Open(DB_TYPE, DB_NAME)
	if err != nil {
		log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
		return false
	}
	defer db.Close()
	
	stmt, err = db.Prepare("delete from users where id=?")
	if err != nil {
		log.Printf("Error preparing, delete from users where id=?, error=%s", err.Error())
		return false
	}
	_, err = stmt.Exec(aUserId)
	if err != nil {
		log.Printf("Error executing, delete from users where id=%d, error=%s", aUserId, err.Error())
		return false
	}
	stmt.Close()
	
	stmt, err = db.Prepare("delete from tokens where userid=?")
	if err != nil {
		log.Printf("Error preparing, delete from tokens where userid=?, error=%s", err.Error())
		return false
	}
	_, err = stmt.Exec(aUserId)
	if err != nil {
		log.Printf("Error executing, delete from tokens where userid=%d, error=%s", aUserId, err.Error())
		return false
	}
	stmt.Close()
	
	stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where %s=?", TABLE_REPORTS_COLUMN_tablename, TABLE_reports, TABLE_REPORTS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing, select tablename from reports where userid=?, error=%s", err.Error())
		return false
	}
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, select tablename from reports where userid=%d, error=%s", aUserId, err.Error())
		return false
	}	
	var tableNames []string
	tableNames = make([]string, 0, 16)
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tableNames = append(tableNames, tableName)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	if len(tableNames) == 0 {
		return true
	}
	
	stmt, err = db.Prepare("drop table if exists ?")
	if err != nil {
		log.Printf("Error preparing, drop table if exists ?, error=%s", err.Error())
	}
	var stmtDelete *sql.Stmt
	stmtDelete, err = db.Prepare(fmt.Sprintf("delete from %s where %s=?", TABLE_reports, TABLE_REPORTS_COLUMN_tablename))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf("delete from %s where %s=?", TABLE_reports, TABLE_REPORTS_COLUMN_tablename), err.Error())
	}
	for _, name := range tableNames {
		_, err = stmt.Exec(name)
		if err != nil {
			log.Printf("Error executing, drop table if exists %s, error=%s", name, err.Error())
		}
		
		_, err = stmtDelete.Exec(name)
		if 	err != nil {
			log.Printf("Error preparing, %s, error=%s", fmt.Sprintf("delete from %s where %s=?", TABLE_reports, TABLE_REPORTS_COLUMN_tablename), err.Error())
		}
	}
	if stmt != nil {stmt.Close()}
	if stmtDelete != nil {stmtDelete.Close()}
	
	return true
}

func DbGetUser(aEmail string, aPassword string) (id int, err error){
	var db *sql.DB = nil
	var stmt *sql.Stmt = nil
	var row *sql.Row = nil
	
	db, err = sql.Open(DB_TYPE, DB_NAME)
	if err != nil {
		log.Printf("Error opening database=%s, error=%s", DB_NAME, err.Error())
		return
	}
	defer db.Close()
	
	db.Prepare(fmt.Sprintf("select * from %s where %s=?", TABLE_users, TABLE_USERS_COLUMN_email))
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

func DbGetUserById(aUserId int) {

}

func DbGetToken(aUserId int, aDb *sql.DB) []string {
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
	
	stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where userid=?", TABLE_TOKENS_COLUMN_token, TABLE_users))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("select %s from %s where userid=?", TABLE_TOKENS_COLUMN_token, TABLE_users), err.Error())
	}
	
	var sliceTokens []string = make([]string, 0, 16)
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", fmt.Sprintf("select %s from %s where userid=?", TABLE_TOKENS_COLUMN_token, TABLE_users), err.Error())
	}
	
	var token string
	for rows.Next() {
		rows.Scan(&token)
		sliceTokens = append(sliceTokens, token)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceTokens
}

//Token is added in DbAddUser.
//This method we can use when we want to add additional tokencs for a user, or if the user does not have a token when we present it to him.
func DbAddToken(aUserId int, aDb *sql.DB) bool {
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
	
	var token string = STR_EMPTY
	token, err = GenerateToken()
	stmt, err = db.Prepare(STMT_INSERT_INTO_TOKENS)
	if err != nil {
		log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_TOKENS, err.Error())
	}
	_, err = stmt.Exec(aUserId, token)
	if err != nil {
		log.Printf("Error executing %s, error=%s", STMT_INSERT_INTO_TOKENS, err.Error())
		stmt.Close()
		return false
	}
	stmt.Close()
	return true
}

// Check to see if we already have reports<aToken> table and use it, if we do not, we have to first create it
func DbAddReport(aToken string, aClientId string, aTime int, aSequence int, aMessage string, aFilePath string) {
	//select tablename from reports where id=token
	
}


func DbDeleteReport(aToken string, aId int) {
	
}

// Delete all records in the reports<token> table
func DbClearReports(aToken string) {
	
}

func DbGetReportsByToken(aToken string, aClientId string) {

}

// We first have to get the token for that user and then we can call the other function
// this will fetch all reports from all reports<token> tables if the user has more than one token
func DbGetReportsByUserId(aUserId int, aClientId string) {
	
}