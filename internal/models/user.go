package models

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
)

type User struct {
	ID uint `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email string `json:"email" db:"email"`
	Role string `json:"role" db:"role"`
	Password []byte `json:"password" db:"password"`
}

func (u *User) SetID(id uint) {
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

func (ufc UserFactory) ParsePassword(password string) []byte {
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
	ListOfRoles []Role
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

func (ufc UserFactory) validate(username string, email string, role Role ,password string) error {
	if err := ufc.validateEmail(email); err != nil {
		return err
	}

	if len(username) < ufc.fc.MinUsernameLen {
		return UsernameTooShortError{
			MinUsernameLen:   ufc.fc.MinUsernameLen,
			ProvidedUsernameLen: len(username),
		}
	}

	if len(password) < ufc.fc.MinPasswordLen {
		return PasswordTooShortError{
			MinPasswordLenLen:   ufc.fc.MinPasswordLen,
			ProvidedPasswordLen: len(password),
		}
	}

	if _, ok := CheckPermission(role.String()); !ok {
		return RoleDoesNotExistError{
			ProvidedRole: role.String(),
			ListOfRoles: GetRolesList(),
		}
	}
	return nil
}

func (ufc UserFactory) NewUser(username string, email string, role Role ,password string) (User, error) {
	if err := ufc.validate(username, email, role, password); err != nil {
		return User{}, err
	}

	return User{
		Username: username,
		Email:    email,
		Role: 	  role.String(),
		Password: ufc.ParsePassword(password),
	}, nil
}

func (ufc UserFactory) NewUserWithEmail(email string, role Role ,password string) (User, error) {
	if err := ufc.validateEmail(email); err != nil {
		return User{}, err
	}

	if len(password) < ufc.fc.MinPasswordLen {
		return User{}, PasswordTooShortError{
			MinPasswordLenLen:   ufc.fc.MinPasswordLen,
			ProvidedPasswordLen: len(password),
		}
	}

	if _, ok := CheckPermission(role.String()); !ok {
		return User{}, RoleDoesNotExistError{
			ProvidedRole: role.String(),
			ListOfRoles: GetRolesList(),
		}
	}

	return User{
		Username: "",
		Email:    email,
		Role: 	  role.String(),
		Password: ufc.ParsePassword(password),
	}, nil
}

func (ufc UserFactory) NewUserWithUsername(username string, role Role ,password string) (User, error) {
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

	if _, ok := CheckPermission(role.String()); !ok {
		return User{}, RoleDoesNotExistError{
			ProvidedRole: role.String(),
			ListOfRoles: GetRolesList(),
		}
	}

	return User{
		Username: username,
		Email:    "",
		Role: 	  role.String(),
		Password: ufc.ParsePassword(password),
	}, nil
}

type UserWithToken struct {
	User *User `json:"user"`
	Token string `json:"token"`
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
	PageSettings *struct{
		PageSize uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
	} `json:"page_settings"`
}