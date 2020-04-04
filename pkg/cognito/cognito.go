package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

// UsersPage return a page of congito users
type UsersPage struct {
	Users []*cognitoidentityprovider.UserType
}

// UsersPageFunc callback which is invoked per page of cognito users returned
type UsersPageFunc func(p *UsersPage) bool

// Service provides cognito related operations
type Service interface {
	ListUsers(userPoolID string, f UsersPageFunc) error
}

type cognitoService struct {
	csvc cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// NewService create a new cognito service adapter
func NewService(config ...*aws.Config) Service {

	sess := session.Must(session.NewSession(config...))
	csvc := cognitoidentityprovider.New(sess)

	return &cognitoService{csvc: csvc}
}

func (ul *cognitoService) ListUsers(userPoolID string, f UsersPageFunc) error {
	return ul.csvc.ListUsersPages(&cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(userPoolID),
		Limit:      aws.Int64(10),
	},
		func(p *cognitoidentityprovider.ListUsersOutput, lastPage bool) bool {
			return f(&UsersPage{Users: p.Users})
		})
}
