package commands

import (
	"bytes"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/cognito-cli/mocks"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

var (
	t1 = time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)
)

func TestFind(t *testing.T) {
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

	fcmd := &FindCmd{
		UserPoolID: "abc123",
		Attributes: []string{"Username"},
		CSV:        true,
	}

	buf := &bytes.Buffer{}

	err := fcmd.Run(&Context{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})

	expected := "Username,Enabled,LastModified\nwolfeidau,false,2016-08-15 00:00:00 +0000 UTC\n"

	assert.NoError(err)
	assert.Equal(expected, buf.String())
}
