package cognito

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/cognito-cli/awsmocks"
)

func TestLogout(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	csvc := awsmocks.NewMockCognitoIdentityProviderAPI(ctrl)

	csvc.EXPECT().AdminUserGlobalSignOut(&cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String("abc123"),
		Username:   aws.String("mark@example.com"),
	}).Return(&cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil)

	cs := cognitoService{csvc: csvc}

	err := cs.Logout("abc123", "mark@example.com")
	assert.NoError(err)
}

func TestDescribePoolAttributes(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	csvc := awsmocks.NewMockCognitoIdentityProviderAPI(ctrl)

	csvc.EXPECT().DescribeUserPool(&cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String("abc123"),
	}).Return(&cognitoidentityprovider.DescribeUserPoolOutput{
		UserPool: &cognitoidentityprovider.UserPoolType{
			SchemaAttributes: []*cognitoidentityprovider.SchemaAttributeType{
				{Name: aws.String("email")},
			},
		},
	}, nil)

	cs := cognitoService{csvc: csvc}

	attributes, err := cs.DescribePoolAttributes("abc123")
	assert.NoError(err)
	assert.Equal([]string{"email"}, attributes)
}
