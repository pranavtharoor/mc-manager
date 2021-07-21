package azure

import (
	"errors"
	"strings"
)

func isValidVMName(s string) bool {
	if len(s) == 0 {
		return false
	}
	tokens := strings.Split("~ ! @ # $ % ^ & * ( ) = + _ [ ] { } \\ | ; : . ' \" , < > / ?", " ")
	tokens = append(tokens, " ")
	for _, token := range tokens {
		if strings.Contains(s, token) {
			return false
		}
	}
	return true
}

func sanitizeVMInputs(resourceGroup string, name string) error {
	if !isValidVMName(resourceGroup) {
		return errors.New("invalid VM resource group provided")
	}
	if !isValidVMName(name) {
		return errors.New("invalid VM name provided")
	}
	return nil
}

func VMStart(resourceGroup string, name string) error {
	if err := verifyLoggedIn(); err != nil {
		return err
	}
	if err := sanitizeVMInputs(resourceGroup, name); err != nil {
		return err
	}
	err := az(outTypeNil, nil, "vm", "start", "-g", resourceGroup, "-n", name, "--no-wait")
	return err
}

func VMDeallocate(resourceGroup string, name string) error {
	if err := verifyLoggedIn(); err != nil {
		return err
	}
	if err := sanitizeVMInputs(resourceGroup, name); err != nil {
		return err
	}
	return az(outTypeNil, nil, "vm", "deallocate", "-g", resourceGroup, "-n", name, "--no-wait")
}
