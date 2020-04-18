package commands

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/cognito-cli/mocks"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

func TestLogout(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cognitoSvc := mocks.NewMockService(ctrl)

	callbackFunc := func(userPoolID string, callback func(p *cognito.UsersPage) bool) {

		p := &cognito.UsersPage{
			Users: []*cognitoidentityprovider.UserType{
				{
					Username:             aws.String("wolfeidau"),
					Attributes:           []*cognitoidentityprovider.AttributeType{},
					UserLastModifiedDate: aws.Time(t1),
				},
			},
		}

		callback(p)
	}

	cognitoSvc.EXPECT().ListUsers("abc123", gomock.Any()).Do(callbackFunc).Return(nil)
	cognitoSvc.EXPECT().Logout("abc123", "wolfeidau").Return(nil)

	logoutcmd := &LogoutCmd{
		UserPoolID: "abc123",
	}

	buf := &bytes.Buffer{}

	err := logoutcmd.Run(&CLIContext{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})
	assert.NoError(err)
	assert.Equal("Found users commencing logout for count=1\n", buf.String())
}
