package commands

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/cognito-cli/mocks"
)

func TestFind(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cognitoSvc := mocks.NewMockService(ctrl)

	cognitoSvc.EXPECT().ListUsers("abc123", gomock.Any()).Return(nil)

	fcmd := &FindCmd{
		UserPoolID: "abc123",
		Attributes: []string{"Username"},
	}

	err := fcmd.Run(&Context{Debug: true, Cognito: cognitoSvc})

	assert.NoError(err)
}
