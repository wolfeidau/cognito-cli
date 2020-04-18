package commands

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jedib0t/go-pretty/table"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
	"github.com/wolfeidau/cognito-cli/pkg/wildcard"
)

// CLIContext cli context used for common options
type CLIContext struct {
	Debug            bool
	DisableLocalTime bool
	Cognito          cognito.Service
	Writer           io.Writer
}

func convertMap(m map[string]string) map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range m {
		result[k] = v
	}

	return result
}

func attrToMap(attributes []*cognitoidentityprovider.AttributeType) map[string]string {
	m := map[string]string{}

	for _, attr := range attributes {
		m[aws.StringValue(attr.Name)] = aws.StringValue(attr.Value)
	}

	return m
}

func buildTableHeader(names []string) table.Row {
	h := table.Row{}

	for _, name := range names {
		h = append(h, name)
	}

	h = append(h, "Enabled")
	h = append(h, "LastModified")

	return h
}

func matchFilters(filters map[string]string, attributes map[string]string) bool {

	// we only need to match one attribute to conclude that a set of attributes match the filter
	for fk, fv := range filters {
		for ak, av := range attributes {
			if wildcard.Match(fk, ak) {
				if wildcard.Match(fv, av) {
					return true
				}
			}
		}
	}

	return false
}

func awsTimeLocal(t *time.Time, local bool) time.Time {
	if local {
		return aws.TimeValue(t).Local()
	}

	return aws.TimeValue(t)
}
