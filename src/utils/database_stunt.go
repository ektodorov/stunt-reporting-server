package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
	
	stmt, err = db.Prepare(STMT_INSERT_INTO_USERS)
	if err != nil {
		log.Println("init, Error prepare insert into users, err=", err);
	}
	passwordHash, err := HashSha1("password1salt1")
	if err != nil {
		log.Println("init, Erro hashing, err=", err)
	}
	_, err = stmt.Exec("user1@mail.com", passwordHash, "address1", "salt1");
	if err != nil {
		log.Println("init, Error insert into users, err=", err);
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

func DbAddUser(aEmail string, aPassword string) {
	//we have to create a token as well
}

func DbDeleteUser(aUserId int) {
	//we have to delete the tokens for that user as well
}

func DbGetUser(aEmail string, aPassword string) {

}

func DbGetUserById(aUserId int) {

}

func DbGetToken(aUserId int) {

}

func DbAddToken(aUserId int) {
	//we generate the token in this method
}

func DbAddReport(aToken string, aClientId string, aTime int, aSequence int, aMessage string, aFilePath string) {
	//check to see if we already have reports<aToken> table and use it, if not we have to first create it
}

func DbDeleteReport() {

}

func DbGetReportsByToken(aToken string, aClientId string) {

}

func DbGetReportsByUserId(aUserId int, aClientId string) {
	//we first have to get the token for that user and then we can call the other function
}