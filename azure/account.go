package azure

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type Account struct {
	ID   string `json:"id"`
	User User
}

type User struct {
	Name string `json:"name"`
}

type Login struct {
	Page string
	Code string
}

func AccountShow() (bool, Account, error) {
	var account Account

	err := az(outTypeJSON, &account, "account", "show")

	if err != nil {
		expectedErr := "Please run 'az login' to setup account."
		if strings.Contains(err.Error(), expectedErr) {
			return false, account, nil
		}
		return false, account, err
	}

	return true, account, nil
}

func AccountLogin(l chan Login, a chan Account) error {
	defer close(l)
	defer close(a)

	var sb strings.Builder
	var out []Account

	c := make(chan string)

	go func() {
		azStart(c, "login", "--use-device-code")
	}()

	l <- parseLogin(<-c)

	for out := range c {
		sb.WriteString(out)
	}

	err := json.Unmarshal([]byte(sb.String()), &out)
	if err != nil {
		return err
	}

	if len(out) == 0 {
		return errors.New("error logging in")
	}

	a <- out[0]

	return err
}

func parseLogin(msg string) Login {
	var login Login
	rePage := regexp.MustCompile("page (.*?) and")
	reCode := regexp.MustCompile("code (.*?) to")
	pageMatch := rePage.FindStringSubmatch(msg)
	codeMatch := reCode.FindStringSubmatch(msg)
	login.Page = pageMatch[1]
	login.Code = codeMatch[1]
	return login
}

func AccountLogout() (bool, error) {
	err := az(outTypeNil, nil, "logout")

	if err != nil {
		expectedErr := "There are no active accounts."
		if strings.Contains(err.Error(), expectedErr) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func verifyLoggedIn() error {
	isLoggedIn, _, err := AccountShow()
	if err != nil {
		return err
	}
	if !isLoggedIn {
		return errors.New("no one's logged in")
	}
	return nil
}
