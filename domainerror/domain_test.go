//go:build unit

package domainerror_test

import (
	"fmt"
	"github.com/julioisaac/users/domainerror"
	"github.com/julioisaac/users/test"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		domainerror map[string]string
		expectedErr *domainerror.Error
	}{
		{
			name: "valid domain error",
			domainerror: map[string]string{
				"Code":    fmt.Sprintf("%v", domainerror.Default),
				"Message": "default",
				"Path":    "",
			},
			expectedErr: domainerror.New(domainerror.Default, "default", nil),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := domainerror.New(tc.expectedErr.Code, tc.expectedErr.Message, nil)
			test.AssertDomainError(t, err, tc.expectedErr)
		})
	}
}
