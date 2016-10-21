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
	
//	stmt, err = db.Prepare(STMT_INSERT_INTO_USERS)
//	if err != nil {
//		log.Println("init, Error prepare insert into users, err=", err);
//	}
//	passwordHash, err := HashSha1("password1salt1")
//	if err != nil {
//		log.Println("init, Erro hashing, err=", err)
//	}
//	_, err = stmt.Exec("user1@mail.com", passwordHash, "salt1");
//	if err != nil {
//		log.Println("init, Error insert into users, err=", err);
//	}
//	stmt.Close();
	
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
		if stmt != nil {stmt.Close()}
		if rows != nil {rows.Close()}
		log.Printf("DbAddUser, aEmail=%s, userId=%d", aEmail, userId)
		
		if userId < 0 {
			return isUserExists, true
		}
		token, err := GenerateToken()
		stmt, err = db.Prepare(STMT_INSERT_INTO_TOKENS)
		if err != nil {
			log.Printf("Error preparing %s, error=%s", STMT_INSERT_INTO_TOKENS, err.Error())
		}
		_, err = stmt.Exec(userId, token)
		if err != nil {
			log.Printf("Error executing %s, error=%s", STMT_INSERT_INTO_TOKENS, err.Error())
		}
		
		return isUserExists, true
	} else {
		return isUserExists, false
	}
}

//Deletes a user and his tokens
func DbDeleteUser(aUserId int) {
	
}

func DbGetUser(aEmail string, aPassword string) {

}

func DbGetUserById(aUserId int) {

}

func DbGetToken(aUserId int) {

}

//Token is added in DbAddUser.
//This method we can use when we want to add additional tokencs for a user, or if the user does not have a token when we present it to him.
func DbAddToken(aUserId int) {
	
}

//check to see if we already have reports<aToken> table and use it, if we do not, we have to first create it
func DbAddReport(aToken string, aClientId string, aTime int, aSequence int, aMessage string, aFilePath string) {
	
}

func DbDeleteReport() {

}

func DbGetReportsByToken(aToken string, aClientId string) {

}

//we first have to get the token for that user and then we can call the other function
//this will fetch all reports from all reports<token> tables if the user has more than one token
func DbGetReportsByUserId(aUserId int, aClientId string) {
	
}