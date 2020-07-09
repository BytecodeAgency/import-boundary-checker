package golistimports_test

import (
	"testing"

	"git.bytecode.nl/foss/import-boundry-checker/langs/golistimports"
	"github.com/stretchr/testify/assert"
)

func TestExtractForSourceFile(t *testing.T) {
	tests := []struct {
		input           string
		expectedImports []string
	}{
		{`package user

import (
	"github.com/bytecodeagency/proj/typings/entities"
	"github.com/bytecodeagency/typings/interactors"
	"github.com/go-playground/validator/v10"
)

type DomainInstance interface {
	// CRUD
	Importer(tester string) error
	CreateUser(user entities.NewUser) (jwt string, err error)

	// Authentication
	AuthenticateUser(email string, password string) (jwt *string, err error)
	CheckUserJwt(jwt string) (user *entities.User, err error)

	// Password reset
	StartPasswordReset(email string) error
	SavePasswordReset(resetToken string, password string) error
}

type userDomain struct {
	interactors.DomainInteractor
}

func NewUserDomainInstance(interactor interactors.DomainInteractor) (DomainInstance, error) {
	domain := userDomain{
		interactor,
	}
	validate := validator.New()
	err := validate.Struct(domain)
	return domain, err
}`, []string{"github.com/bytecodeagency/proj/typings/entities", "github.com/bytecodeagency/typings/interactors", "github.com/go-playground/validator/v10"}},
		{`package sometest

import "github.com/bytecodeagency/typings/entities"

func importer () int {
  return 0
}
`, []string{"github.com/bytecodeagency/typings/entities"}},
	}
	for _, test := range tests {
		imports, err := golistimports.ExtractForSourceFile(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedImports, imports)
	}
}
