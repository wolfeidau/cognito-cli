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

func TestLs(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cognitoSvc := mocks.NewMockService(ctrl)

	callbackFunc := func(callback func(p *cognito.UserPoolsPage) bool) {

		p := &cognito.UserPoolsPage{
			UserPools: []*cognitoidentityprovider.UserPoolDescriptionType{
				{
					Id:           aws.String("abc123"),
					Name:         aws.String("test pool"),
					CreationDate: aws.Time(t1),
				},
			},
		}

		callback(p)
	}

	cognitoSvc.EXPECT().ListPools(gomock.Any()).Do(callbackFunc).Return(nil)

	fcmd := &LsCmd{
		CSV: true,
	}

	buf := &bytes.Buffer{}

	err := fcmd.Run(&Context{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})
	assert.NoError(err)
	assert.Equal("ID,Name,Created\nabc123,test pool,2016-08-15 00:00:00 +0000 UTC\n", buf.String())
}

func TestLsNoPools(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cognitoSvc := mocks.NewMockService(ctrl)

	cognitoSvc.EXPECT().ListPools(gomock.Any()).Return(nil)

	fcmd := &LsCmd{}

	buf := &bytes.Buffer{}

	err := fcmd.Run(&Context{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})
	assert.NoError(err)
	assert.Equal("No pools found.\n", buf.String())
}
