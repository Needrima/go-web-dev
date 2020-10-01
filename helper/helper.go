package helper

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CheckInSlice(xs []string, s string) (int, bool) {
	for i, v := range xs {
		if v == s {
			return i, true
		}
	}
	return -1, false
}

func Found(xexp []string, exp string) bool {
	for _, v := range xexp {
		if v == exp {
			return true
		}
	}
	return false
}

func CheckName(name string) bool {
	//check length of name
	name_length := len(name) < 3
	//chech if name has only english alphabets
	english, _ := regexp.MatchString("[A-Za-z]", name)

	if name_length || !english {
		return false
	}
	return true
}

func ConvertName(name string) string {
	return strings.ToLower(name)
}

func CheckPassword(password string) bool {
	//check if password is longer than 7 characters
	pass_length := len(password) < 7
	//check if password contains only english numbers and alphabets
	english, _ := regexp.MatchString("[A-Za-z0-9]", password)
	//check if password contains at least one number
	contains_any_num := strings.ContainsAny(password, "0123456789")

	if pass_length || !english || !contains_any_num {
		return false
	}
	return true
}

func CheckEmail(email string) bool {
	//check email length
	email_length := len(email) < 15
	//check if email ends with .com or .co.uk
	//chech if only english alphabets
	containsat, _ := regexp.MatchString("@{1}", email)

	if !containsat || email_length {
		return false
	}
	return true
}

func CheckDOB(year, month, day string) bool {
	//get current year
	currentyear := time.Now().Year()
	//convert values to int
	yr, _ := strconv.Atoi(year)
	mth, _ := strconv.Atoi(month)
	dy, _ := strconv.Atoi(day)
	//check if user is underage
	if currentyear-yr < 15 {
		return false
	} else if mth == 2 && yr%4 == 0 && dy > 29 {
		return false
	} else if mth == 2 && yr%4 != 0 && dy > 28 {
		return false
	} else if mth == 4 || mth == 6 || mth == 9 || mth == 11 && dy > 30 {
		return false
	} else if mth == 1 || mth == 3 || mth == 5 || mth == 7 || mth == 8 || mth == 10 && dy > 31 {
		return false
	} else {
		return true
	}
}

type user struct {
	Sname, Fname, Oname, Uname, Pword, Email string
	Dob                                      []string
}

var sessiondb = map[string]string{}
var userdb = map[string]user{}

func GetUser(w http.ResponseWriter, r *http.Request) user {
	cookie, _ := r.Cookie("session")

	un := sessiondb[cookie.Value]
	u, ok := userdb[un]
	if !ok {
		http.Error(w, "User does not exist", 500)
	}
	return u
}
