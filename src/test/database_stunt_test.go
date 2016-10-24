package test

import (
	"strconv"
    "testing"
    "utils"
)

func TestDbInit(t *testing.T) {
	utils.DbInit()
}

func TestDbAddUser(t *testing.T) {
	var user = "mail0@mail.com"
	var password = "password0"
	isUserExists, isUserAdded, err := utils.DbAddUser(user, password, nil)
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if err != nil {
		t.Logf("error=%s", err.Error())
	}
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
}

func TestDbDeleteUser(t *testing.T) {
	var user = "mail1@mail.com"
	var password = "password1"
	isUserExists, isUserAdded, err := utils.DbAddUser(user, password, nil)
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if err != nil {
		t.Logf("error=%s", err.Error())
	}
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, errorGetUser := utils.DbGetUser(user, password, nil)
	if err != nil {
		t.Error("Error getting user id, errorGetUser=%s", errorGetUser.Error())
	}
	
	isUserDeleted := utils.DbDeleteUser(id, nil)
	t.Logf("isUserDeleted=%t", isUserDeleted)
}

func TestDbGetUser(t *testing.T) {
	var user = "mail2@mail.com"
	var password = "password2"
	isUserExists, isUserAdded, err := utils.DbAddUser(user, password, nil)
	if err != nil {
		t.Logf("error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err := utils.DbGetUser(user, password, nil)
	if err != nil {
		t.Error("Error getting user id")
	}
	t.Logf("user email=%s, id=%d", user, id)
}

func TestDbGetToken(t *testing.T) {
	var username = "mail3@mail.com"
	var password = "password3"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Logf("error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	sliceTokens := utils.DbGetToken(id, nil)
	for idx, token := range sliceTokens {
		t.Logf("sliceTokens, idx=%d, token=%s", idx, token)
	}
}

func TestDbAddToken(t *testing.T) {
	var username = "mail4@mail.com"
	var password = "password4"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Log("TestDbAddToken, error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	isTokenAdded := utils.DbAddToken(id, nil)
	t.Logf("isTokenAdded=%t", isTokenAdded)
}

func TestDbAddReport(t *testing.T) {
	var username = "mail5@mail.com"
	var password = "password5"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Log("TestDbAddToken, error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	sliceTokens := utils.DbGetToken(id, nil)
	for idx, token := range sliceTokens {
		t.Logf("sliceTokens, idx=%d, token=%s", idx, token)
	}
	
	if len(sliceTokens) > 0 {
		utils.DbAddReport(id, sliceTokens[0], "clientid4", 1234, 4, "message4", "filepath4", nil)
	}
}

func TestDbDeleteReport(t *testing.T) {
	var username = "mail5@mail.com"
	var password = "password5"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Log("TestDbAddToken, error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	sliceTokens := utils.DbGetToken(id, nil)
	for idx, token := range sliceTokens {
		t.Logf("sliceTokens, idx=%d, token=%s", idx, token)
	}
	
	if len(sliceTokens) > 0 {
		for x := 0; x < 3; x++ {
			utils.DbAddReport(id, sliceTokens[0], ("clientid_" + strconv.Itoa(x)), (12345 + x), x, ("message5_" + strconv.Itoa(x)), ("filepath5_" + strconv.Itoa(x)), nil)
		}
	}
	
	sliceReports, endNumber := utils.DbGetReportsByToken(sliceTokens[0], "client6", 0, 2, nil)
	t.Logf("endNumber=%d", endNumber)
	for idx, report := range sliceReports {
		t.Logf("idx=%d, report=%+v\n", idx, report)
		utils.DbDeleteReport(sliceTokens[0], report.Id, nil)
	}
}

func TestDbClearReports(t *testing.T) {
	var username = "mail6@mail.com"
	var password = "password6"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Log("TestDbAddToken, error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	sliceTokens := utils.DbGetToken(id, nil)
	for idx, token := range sliceTokens {
		t.Logf("sliceTokens, idx=%d, token=%s", idx, token)
	}
	
	if len(sliceTokens) > 0 {
		for x := 0; x < 2; x++ {
			utils.DbAddReport(id, sliceTokens[0], "clientid6", (123456 + x), x, ("message6_" + strconv.Itoa(x)), ("filepath6_" + strconv.Itoa(x)), nil)
		}
	}

	utils.DbClearReports(sliceTokens[0], nil)
}

func TestDbGetReportsByToken(t *testing.T) {
	var username = "mail6@mail.com"
	var password = "password6"
	var err error
	var id int
	
	isUserExists, isUserAdded, err := utils.DbAddUser(username, password, nil)
	if err != nil {
		t.Log("TestDbAddToken, error=%s", err.Error())
	}
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
	
	id, err = utils.DbGetUser(username, password, nil)
	if err != nil {
		t.Errorf("Error getting user, username=%s, password=%s", username, password)
	}
	t.Logf("user username=%s, id=%d", username, id)
	
	sliceTokens := utils.DbGetToken(id, nil)
	for idx, token := range sliceTokens {
		t.Logf("sliceTokens, idx=%d, token=%s", idx, token)
	}
	
	if len(sliceTokens) > 0 {
		for x := 0; x < 10; x++ {
			utils.DbAddReport(id, sliceTokens[0], "clientid6_" + strconv.Itoa(x), 123456 + x, x, ("message6_" + strconv.Itoa(x)), "filepath6_" + strconv.Itoa(x), nil)
		}
	}
	
	
	sliceReports, endNumber := utils.DbGetReportsByToken(sliceTokens[0], "client6", 0, 2, nil)
	t.Logf("endNumber=%d", endNumber)
	for idx, report := range sliceReports {
		t.Logf("idx=%d, report=%+v", idx, report)
	}
}

