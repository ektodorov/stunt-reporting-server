package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"objects"
//	"time"
)

func DbInit() error {
	var err error = nil;
	var db *sql.DB = nil;
	var stmt *sql.Stmt = nil;

	db, err = sql.Open(DB_TYPE, DB_NAME);
	if err != nil {
		log.Println("init, Error opening db", DB_NAME, ", err=", err);
		return err;
	}
	defer db.Close();
	
	stmt, err = db.Prepare(STMT_CREATE_TABLE_USERS);
	if err != nil {
		log.Println("init, Error preparing create table users stmt, err=", err);
		return err;
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
	return nil
}


// Adds a user in the users table if it does not exist.
// Token is also generated for the user. If token creation fails for some reason,
// we would not fail the user creation. We would just have to check if a token exists for that user when we are presenting it to the
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
		
		//create and insert token
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
		DbAddToken(userId, db)
		
		return isUserExists, true, nil
	} else {
		return isUserExists, false, nil
	}
}

//Deletes a user and his tokens and the reports<token> tables of that user
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
	
	stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where %s=?", TABLE_TOKENS_COLUMN_token, TABLE_tokens, TABLE_TOKENS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_TOKENS_COLUMN_token, TABLE_tokens, TABLE_TOKENS_COLUMN_userid), err.Error())
		return false
	}
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=%d", TABLE_TOKENS_COLUMN_token, TABLE_tokens, TABLE_TOKENS_COLUMN_userid, aUserId), err.Error())
		return false
	}	
	//get all tokens of user
	var sliceTokens []string = make([]string, 0, 16)
	for rows.Next() {
		var token string
		rows.Scan(&token)
		sliceTokens = append(sliceTokens, token)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	if len(sliceTokens) == 0 {
		return true
	}
	
	//Delete tokens of user (delete from table tokens)
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s where %s=?", TABLE_tokens, TABLE_TOKENS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing, delete from %s where userid=?, error=%s", TABLE_tokens, err.Error())
		return false
	}
	_, err = stmt.Exec(aUserId)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf("delete from %s where %s=%d", TABLE_tokens, TABLE_TOKENS_COLUMN_userid, aUserId), err.Error())
		return false
	}
	stmt.Close()
	
	//Drop reports<token> tables for the user
	stmt, err = db.Prepare("drop table if exists ?")
	if err != nil {
		log.Printf("Error preparing, drop table if exists ?, error=%s", err.Error())
	}
	for _, name := range sliceTokens {
		_, err = stmt.Exec((TABLE_reports + name))
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

// Get all tokens of a user.
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
	
	stmt, err = db.Prepare(fmt.Sprintf("select %s from %s where %s=?", TABLE_TOKENS_COLUMN_token, TABLE_users, TABLE_TOKENS_COLUMN_userid))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_TOKENS_COLUMN_token, TABLE_users, TABLE_TOKENS_COLUMN_userid), err.Error())
	}
	
	rows, err = stmt.Query(aUserId)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select %s from %s where %s=?", TABLE_TOKENS_COLUMN_token, TABLE_users, TABLE_TOKENS_COLUMN_userid), err.Error())
	}
	var sliceTokens []string = make([]string, 0, 16)
	for rows.Next() {
		var token string
		rows.Scan(&token)
		sliceTokens = append(sliceTokens, token)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceTokens
}

// Token is added in DbAddUser.
// This method we can use when we want to add additional tokencs for a user, or if the user does not have a token when we present it to him.
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
		log.Printf("Error executing %s, values userId=%d, token=%s, error=%s", STMT_INSERT_INTO_TOKENS, aUserId, token, err.Error())
		stmt.Close()
		return false
	}
	if stmt != nil {stmt.Close()}
	return true
}

// Check to see if we already have reports<aToken> table and use it, if we do not, we have to first create it
func DbAddReport(aUserId int, aToken string, aClientId string, aTime int, aSequence int, aMessage string, aFilePath string, aDb *sql.DB) {
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
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aToken))
	if err != nil {
		log.Printf("Error creating, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aToken), err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_CREATE_TABLE_REPORTS, aToken), err.Error())
	}
	if stmt != nil {stmt.Close()}
	
	stmt, err = db.Prepare(fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aToken))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aToken), err.Error())
	}
	_, err = stmt.Exec(aClientId, aTime, aSequence, aMessage, aFilePath)
	if err != nil {
		log.Printf("Error executing, %s, error=%s", fmt.Sprintf(STMT_INSERT_INTO_REPORTS, aToken), err.Error())
	}
	if stmt != nil {stmt.Close()}
}


func DbDeleteReport(aToken string, aId int, aDb *sql.DB) {
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
	
	stmt, err = db.Prepare(fmt.Sprintf("delete from %s%s where %s=?", TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", fmt.Sprintf("delete from %s%s where %s=?", TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id), err.Error())
	}
	_, err = stmt.Exec(aId)
	if err != nil {
		log.Printf("Error deleting, %s, error=%s", 
				fmt.Sprintf("delete from %s%s where %s=aId", TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id, aId), err.Error())
	}
	if stmt != nil {stmt.Close()}
}

// Delete all records in the reports<token> table
func DbClearReports(aToken string, aDb *sql.DB) {
	// Instead of deleting all from reports<aToken> we can just - drop table if exists reports<aToken>
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
	
	stmt, err = db.Prepare(fmt.Sprintf("delete * from %s%s", TABLE_reports, aToken))
	if err != nil {
		log.Printf("Error preparing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aToken), err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Error executing %s, error=%s", fmt.Sprintf("delete * from %s%s", TABLE_reports, aToken), err.Error())
	}
	stmt.Close()
}

func DbGetReportsByToken(aToken string, aClientId string, aStartNum int, aPageSize int, aDb *sql.DB) (sliceReports []*objects.Report, endNum int) {
	endNum = -1
	sliceReports = nil
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
	
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s%s where %s > ? order by %s limit ?",
					TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_id))
	if err != nil {
		log.Printf("Error preparing, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s > ? order by %s limit ?", TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_id),
			err.Error())
	}
	
	rows, err = stmt.Query(aToken)
	if err != nil {
		log.Printf("Error quering, %s, error=%s", 
			fmt.Sprintf("select * from %s%s where %s > ? order by %s limit ?", TABLE_reports, aToken, TABLE_REPORTS_COLUMN_id, TABLE_REPORTS_COLUMN_id),
			err.Error())
	}
	
	sliceReports = make([]*objects.Report, 0, 64)
	for rows.Next() {
		var id int
		var clientId string
		var time int
		var sequence int
		var message string
		var filePath string
		rows.Scan(&id, &clientId, &sequence, &message, &filePath)
		var report = new(objects.Report)
		report.Id = id
		report.ClientId = clientId
		report.Time = time
		report.Message = message
		report.FilePath = filePath
		sliceReports = append(sliceReports, report)
	}
	if rows != nil {rows.Close()}
	if stmt != nil {stmt.Close()}
	
	return sliceReports, endNum
}
