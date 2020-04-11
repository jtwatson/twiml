package twiml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ttacon/libphonenumber"
)

// validPhoneNumber checks that a valid phone number is provided
// param of "allowempty" will allow a nil value
func validPhoneNumber(v interface{}, param string) error {
	switch num := v.(type) {
	case string:
		if num == "" {
			if param == "allowempty" {
				return nil
			}
			return errors.New("Required")
		}
		return validatePhoneNumber(num)
	case *string:
		if num == nil {
			if param == "allowempty" {
				return nil
			}
			return errors.New("Required")
		}
		return validatePhoneNumber(*num)
	default:
		return fmt.Errorf("validDatastoreKey: Unexpected type %T", num)
	}
}

func validatePhoneNumber(num string) error {
	n, err := libphonenumber.Parse(num, "US")
	if err != nil {
		return errors.New("Invalid phone number")
	}
	if !libphonenumber.IsValidNumber(n) {
		return errors.New("Invalid phone number")
	}
	return nil
}

func validateKeyPadEntry(v interface{}, param string) error {
	switch num := v.(type) {
	case string:
		if num == "" {
			if param == "allowempty" {
				return nil
			}
			return errors.New("Required")
		}
		return validateNumericPoundStar(num)
	case *string:
		if num == nil {
			if param == "allowempty" {
				return nil
			}
			return errors.New("Required")
		}
		return validateNumericPoundStar(*num)
	default:
		return fmt.Errorf("validDatastoreKey: Unexpected type %T", num)
	}
}

func validateNumericPoundStar(v string) error {
	return characterList(v, "0123456789#*")
}

// characterList checks a string against a list of acceptable characters.
// returns an erro if a character is found which is not in charList
func characterList(s string, charList string) error {
	for _, c := range s {
		if strings.Index(charList, string(c)) == -1 {
			return fmt.Errorf("Invalid: character '%s' is not allowed", string(c))
		}
	}
	return nil
}
