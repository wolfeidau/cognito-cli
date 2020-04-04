package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_matchFilters(t *testing.T) {

	assert := require.New(t)

	type args struct {
		filters    map[string]string
		attributes map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should match with key and value",
			args: args{
				filters: map[string]string{
					"*:customerId": "*-eb74116dbd0b",
				},
				attributes: map[string]string{
					"Username":       "wolfeidau",
					"email":          "mark@example.com",
					"dev:customerId": "fa0bdfda-5b24-4c70-a7c6-eb74116dbd0b",
				},
			},
			want: true,
		},
		{
			name: "should match when key and value match at in one attribute",
			args: args{
				filters: map[string]string{
					"*:customerId": "*-f6faaabe76db",
				},
				attributes: map[string]string{
					"Username":       "wolfeidau",
					"email":          "mark@example.com",
					"dev:customerId": "fa0bdfda-5b24-4c70-a7c6-eb74116dbd0b",
					"dev:inviteId":   "f58d7c7d-6120-40b3-bfe5-f6faaabe76db",
				},
			},
			want: false,
		},
		{
			name: "should not match with key and value",
			args: args{
				filters: map[string]string{
					"customerId": "eb74116dbd0b",
				},
				attributes: map[string]string{
					"Username":           "wolfeidau",
					"email":              "mark@example.com",
					"dev:customerId":     "fa0bdfda-5b24-4c70-a7c6-eb74116dbd0b",
					"dev:dev:customerId": "fa0bdfda-5b24-4c70-a7c6-eb74116dbd0b",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchFilters(tt.args.filters, tt.args.attributes)
			assert.Equal(tt.want, got)
		})
	}
}
