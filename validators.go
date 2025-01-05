package twiml

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ttacon/libphonenumber"
)

const allowempty = "allowempty"

// validFromOrTo checks that a valid phone number or sip uri is provided
// param of "allowempty" will allow a nil value
func validFromOrTo(v interface{}, param string) error {
	err := validPhoneNumber(v, param)
	if err == nil {
		return nil
	}

	if err := validSIPURI(v, param); err == nil {
		return nil
	}

	return err
}

// validPhoneNumber checks that a valid phone number is provided
// param of "allowempty" will allow a nil value
func validPhoneNumber(v interface{}, param string) error {
	switch num := v.(type) {
	case string:
		if num == "" {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}

		return validatePhoneNumber(num)
	case *string:
		if num == nil {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}

		return validatePhoneNumber(*num)
	default:
		return fmt.Errorf("validPhoneNumber: Unexpected type %T", num)
	}
}

func validatePhoneNumber(num string) error {
	n, err := libphonenumber.Parse(num, "US")
	if err != nil {
		return errors.New("invalid phone number")
	}
	if !libphonenumber.IsValidNumber(n) {
		return errors.New("invalid phone number")
	}

	return nil
}

func validateKeyPadEntry(v interface{}, param string) error {
	switch num := v.(type) {
	case string:
		if num == "" {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}

		return validateNumericPoundStar(num)
	case *string:
		if num == nil {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}

		return validateNumericPoundStar(*num)
	default:
		return fmt.Errorf("validateKeyPadEntry: Unexpected type %T", num)
	}
}

func validateNumericPoundStar(v string) error {
	return characterList(v, "0123456789#*")
}

// characterList checks a string against a list of acceptable characters.
// returns an erro if a character is found which is not in charList
func characterList(s, charList string) error {
	for _, c := range s {
		if !strings.Contains(charList, string(c)) {
			return fmt.Errorf("invalid: character '%s' is not allowed", string(c))
		}
	}

	return nil
}

// validSIPURI checks that a valid sip uri is provided
// param of "allowempty" will allow a nil value
func validSIPURI(v interface{}, param string) error {
	switch num := v.(type) {
	case string:
		if num == "" {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}
		_, err := parseSIPURI(num)

		return err
	case *string:
		if num == nil {
			if param == allowempty {
				return nil
			}

			return errors.New("Required")
		}
		_, err := parseSIPURI(*num)

		return err
	default:
		return fmt.Errorf("validPhoneNumber: Unexpected type %T", num)
	}
}

func parseSIPURI(uri string) (*url.URL, error) {
	uri = strings.ToLower(uri)

	if !strings.HasPrefix(uri, "sip") {
		return nil, errors.New("invalid SIP URI")
	}

	// Insert the // after the Schema to enable full parsing and avoid Opaque
	if strings.HasPrefix(uri, "sips:") {
		uri = strings.Replace(uri, "sips:", "sips://", 1)
	} else if strings.HasPrefix(uri, "sip:") {
		uri = strings.Replace(uri, "sip:", "sip://", 1)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("invalid SIP URI")
	}

	// Schema should be valid
	if u.Scheme != "sip" && u.Scheme != "sips" {
		return u, errors.New("invalid SIP URI")
	}

	// Path and Opaque should be empty
	if u.Path != "" || u.Opaque != "" {
		return u, errors.New("invalid SIP URI")
	}

	// Host and User should be provided
	if u.Host == "" || u.User.String() == "" {
		return u, errors.New("invalid SIP URI")
	}

	return u, nil
}
