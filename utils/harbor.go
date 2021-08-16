package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type User struct {
	Email    string
	Username string
}

func GetToken(loginId string, passwd string) string {
	byteString := []byte(loginId + ":" + passwd)
	return base64.StdEncoding.EncodeToString(byteString)
}

func GetUsers(harborHost string, token string) (*[]User, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", harborHost+"/api/v2.0/users?sort=ascending&page=1&page_size=100", nil)
	if err != nil {
		log.Fatal(err)
		return &[]User{}, err
	}
	req.Header.Add("Authorization", "Basic "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed %v", err)
		return &[]User{}, err
	}

	if resp.StatusCode != 200 {
		return &[]User{}, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return &[]User{}, err
	}

	users := make([]User, 30, 100)
	err = json.Unmarshal(body, &users)
	return &users, err

}

func GetUserEmail(username string, harborHost string, token string) (string, error) {
	users, err := GetUsers(harborHost, token)
	if err != nil {
		log.Fatal(err)
	}

	email := ""
	for _, v := range *users {
		if v.Username == username {
			email = v.Email
		}
	}
	if email != "" {
		return email, err
	} else {
		return "", fmt.Errorf("No such user")
	}

}

func GetCurrentUser(harborHost string, token string) (*User, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", harborHost+"/api/v2.0/users/current", nil)
	if err != nil {
		log.Fatal(err)
		return &User{}, err
	}

	req.Header.Add("Authorization", "Basic "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed %v", err)
		return &User{}, err
	}

	if resp.StatusCode != 200 {
		return &User{}, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return &User{}, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	return &user, err
}
