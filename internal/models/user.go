package models

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
)

const (
	ADMIN = "admin"
	USER = "user"
	VIEWER = "viewer"
)
var roles = map[string]string{
	ADMIN: ADMIN,
	USER: USER,
	VIEWER: VIEWER,
}

type User struct {
	ID int `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email string `json:"email" db:"email"`
	Role string `json:"role" db:"role"`
	Password []byte `json:"password" db:"password"`
}

func (u *User) SetID(id int) {
	u.ID = id
}

//emailValidate - should check when ever email is valid for user creation
type emailValidate func(email string) error
//passwordHash - function should return hashed password
type passwordHash func(password string) []byte

type FactoryConfig struct {
	//Just basic constraints, can be added more
	MinPasswordLen int
	MinUsernameLen int
	ValidateEmail emailValidate
	ParsePassword passwordHash
}

func (f FactoryConfig) Validate() error {
	var err error
	if f.MinUsernameLen < 4 {
		err = multierr.Append(
			err,
			errors.New(fmt.Sprintf("MinUsernameLen should be grater then 4, got: %d", f.MinUsernameLen)),
		)
	}

	if f.MinPasswordLen < 4 {
		err = multierr.Append(
			err,
			errors.New(fmt.Sprintf("MinPasswordLen should be grater then 4, got: %d", f.MinPasswordLen)),
		)
	}

	if f.ParsePassword == nil {
		err = multierr.Append(
			err,
			errors.New("ParsePassword function not specified, passwords must be hashed"),
		)
	}
	return err
}

type UserFactory struct {
	fc FactoryConfig
}

func NewUserFactory(fc FactoryConfig) (UserFactory, error) {
	if err := fc.Validate(); err != nil {
		return UserFactory{}, err
	}
	return UserFactory{fc: fc}, nil
}

func (ufc UserFactory) validateEmail(email string) error {
	if ufc.fc.ValidateEmail == nil {
		return nil
	}

	if err := ufc.fc.ValidateEmail(email); err != nil {
		return err
	}

	return nil
}

func (ufc UserFactory) parsePassword(password string) []byte {
	if ufc.fc.ParsePassword == nil {
		return []byte(password)
	}

	return ufc.fc.ParsePassword(password)
}

type PasswordTooShortError struct {
	MinPasswordLenLen int
	ProvidedPasswordLen int
}

func (p PasswordTooShortError) Error() string {
	return fmt.Sprintf(
		"Provided password is too short, min length: %d, provided password lentgh: %d",
		p.MinPasswordLenLen,
		p.ProvidedPasswordLen,
	)
}

type UsernameTooShortError struct {
	MinUsernameLen int
	ProvidedUsernameLen int
}

func (p UsernameTooShortError) Error() string {
	return fmt.Sprintf(
		"Username password is too short, min length: %d, provided password lentgh: %d",
		p.MinUsernameLen,
		p.ProvidedUsernameLen,
	)
}

type RoleDoesNotExistError struct {
	ListOfRoles map[string]string
	ProvidedRole string
}

func (p RoleDoesNotExistError) Error() string {
	return fmt.Sprintf(
		"Provided user role does not exist, user roles: %v, provided user role: %q",
		p.ListOfRoles,
		p.ProvidedRole,
	)
}

//func (ufc UserFactory) NewUnAuthUser(username, password string) (User, error) {
//	if len(username) < ufc.fc.MinUsernameLen {
//		return User{}, UsernameTooShortError{
//			MinUsernameLen:   ufc.fc.MinUsernameLen,
//			ProvidedUsernameLen: len(username),
//		}
//	}
//
//	if len(password) < ufc.fc.MinPasswordLen {
//		return User{}, PasswordTooShortError{
//			MinPasswordLenLen:   ufc.fc.MinPasswordLen,
//			ProvidedPasswordLen: len(password),
//		}
//	}
//
//	return User{
//		Username: username,
//		Password: ufc.parsePassword(password),
//	}, nil
//}

func (ufc UserFactory) NewUser(username, email, role, password string) (User, error) {
	if err := ufc.validateEmail(email); err != nil {
		return User{}, err
	}

	if len(username) < ufc.fc.MinUsernameLen {
		return User{}, UsernameTooShortError{
			MinUsernameLen:   ufc.fc.MinUsernameLen,
			ProvidedUsernameLen: len(username),
		}
	}

	if len(password) < ufc.fc.MinPasswordLen {
		return User{}, PasswordTooShortError{
			MinPasswordLenLen:   ufc.fc.MinPasswordLen,
			ProvidedPasswordLen: len(password),
		}
	}

	if _, ok := roles[role]; !ok {
		return User{},RoleDoesNotExistError{
			ProvidedRole: role,
			ListOfRoles: roles,
		}
	}

	return User{
		Username: username,
		Email:    email,
		Role: 	  role,
		Password: ufc.parsePassword(password),
	}, nil
}

type FindUserRequest struct {
	Username *struct {
		Like string `json:"LIKE"`
		Eq   string `json:"EQ"`
	} `json:"username"`
	Email *struct {
		Like string `json:"LIKE"`
		Eq   string `json:"EQ"`
	} `json:"email"`
	Role *struct{
		Eq   string `json:"EQ"`
	} `json:"role"`
	ID *struct {
		Eq int `json:"EQ"`
	} `json:"id"`
}