package test

import (
    "testing"
    "utils"
)

func TestDbInit(t *testing.T) {
	utils.DbInit()
}

func TestDbAddUser(t *testing.T) {
	var user = "mail6@mail.com"
	var password = "password6"
	isUserExists, isUserAdded := utils.DbAddUser(user, password)
	t.Logf("isUserExists=%t, isUserAdded=%t", isUserExists, isUserAdded)
	if isUserExists == false && isUserAdded == false {
		t.Error("Error adding user")
	}
}

