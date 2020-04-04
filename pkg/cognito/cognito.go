package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

const (
	listUsersPageLimit  = 10
	listPoolsMaxResults = 60
)

// UsersPage return a page of AWS Cognito users
type UsersPage struct {
	Users []*cognitoidentityprovider.UserType
}

// UserPoolsPage return a page of AWS Cognito pools
type UserPoolsPage struct {
	UserPools []*cognitoidentityprovider.UserPoolDescriptionType
}

// UsersPageFunc callback which is invoked per page of AWS Cognito users returned
type UsersPageFunc func(p *UsersPage) bool

// UserPoolsPageFunc callback which is invoked per page of AWS Cognito pools returned
type UserPoolsPageFunc func(p *UserPoolsPage) bool

// Service provides AWS Cognito related operations
type Service interface {
	ListUsers(userPoolID string, f UsersPageFunc) error
	ListPools(f UserPoolsPageFunc) error
	Logout(userPoolID string, username string) error
	DescribePoolAttributes(userPoolID string) ([]string, error)
}

type cognitoService struct {
	csvc cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// NewService create a new AWS Cognito service adapter
func NewService(config ...*aws.Config) Service {

	sess := session.Must(session.NewSession(config...))
	csvc := cognitoidentityprovider.New(sess)

	return &cognitoService{csvc: csvc}
}

func (ul *cognitoService) ListUsers(userPoolID string, f UsersPageFunc) error {
	return ul.csvc.ListUsersPages(&cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(userPoolID),
		Limit:      aws.Int64(listUsersPageLimit),
	},
		func(p *cognitoidentityprovider.ListUsersOutput, lastPage bool) bool {
			return f(&UsersPage{Users: p.Users})
		})
}

func (ul *cognitoService) ListPools(f UserPoolsPageFunc) error {
	return ul.csvc.ListUserPoolsPages(&cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(listPoolsMaxResults),
	},
		func(p *cognitoidentityprovider.ListUserPoolsOutput, lastPage bool) bool {
			return f(&UserPoolsPage{UserPools: p.UserPools})
		})
}

func (ul *cognitoService) Logout(userPoolID string, username string) error {
	_, err := ul.csvc.AdminUserGlobalSignOut(&cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return err
	}
	return nil
}

func (ul *cognitoService) DescribePoolAttributes(userPoolID string) ([]string, error) {

	res, err := ul.csvc.DescribeUserPool(&cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil {
		return nil, err
	}

	result := []string{}

	for _, sattr := range res.UserPool.SchemaAttributes {
		result = append(result, aws.StringValue(sattr.Name))
	}

	return result, nil
}
