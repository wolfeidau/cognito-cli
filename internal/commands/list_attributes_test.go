package commands

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/cognito-cli/mocks"
)

func TestListAttributes(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cognitoSvc := mocks.NewMockService(ctrl)

	cognitoSvc.EXPECT().DescribePoolAttributes("abc123").Return([]string{"name", "given_name", "family_name"}, nil)

	listAttributesCmd := &ListAttributesCmd{
		UserPoolID: "abc123",
	}

	buf := &bytes.Buffer{}

	err := listAttributesCmd.Run(&CLIContext{Debug: true, DisableLocalTime: true, Cognito: cognitoSvc, Writer: buf})

	expected := "+-------------+\n| NAME        |\n+-------------+\n| name        |\n| given_name  |\n| family_name |\n+-------------+\n"

	assert.NoError(err)
	assert.Equal(expected, buf.String())
}
