package bot

import (
	"fmt"

	"github.com/pranavtharoor/mc-manager/azure"
)

func azureLogin(send func(string)) {
	accChan := make(chan azure.Account)
	loginChan := make(chan azure.Login)

	go func() {
		err := azure.AccountLogin(loginChan, accChan)
		if err != nil {
			send(err.Error())
		}
	}()

	login := <-loginChan
	loginMsg := fmt.Sprintf("Log in: %s\nCode: `%s`", login.Page, login.Code)
	send(loginMsg)

	account := <-accChan
	accMsg := fmt.Sprintf("Logged in as: _%s_", account.User.Name)
	send(accMsg)
}

func azureAccount() string {
	isLoggedIn, account, err := azure.AccountShow()
	if err != nil {
		return err.Error()
	}

	if !isLoggedIn {
		return "No one's logged in"
	}

	return fmt.Sprintf("Logged in as: _%s_", account.User.Name)
}

func azureLogout() string {
	wasLoggedIn, err := azure.AccountLogout()
	if err != nil {
		return err.Error()
	}

	if !wasLoggedIn {
		return "No one's logged in"
	}

	return "Logged out"
}
