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

func TestExport(t *testing.T) {
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
	cognitoSvc.EXPECT().DescribePoolAttributes("abc123").Return([]string{"name", "given_name", "family_name"}, nil)

	exportCmd := &ExportCmd{
		UserPoolID: "abc123",
	}

	buf := &bytes.Buffer{}

	err := exportCmd.Run(&Context{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})

	expected := "Username,name,given_name,family_name,Enabled,LastModified\nwolfeidau,,,,false,2016-08-15 00:00:00 +0000 UTC\n"

	assert.NoError(err)
	assert.Equal(expected, buf.String())
}
