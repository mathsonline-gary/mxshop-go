package gen

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	e := errorWrapper{
		Errors: []*errorInfo{
			{
				Name:       "USER_NOT_FOUND",
				Value:      "0",
				CamelValue: "UserNotFound",
			},
			{
				Name:       "USER_EXISTS",
				Value:      "1",
				CamelValue: "UserExists",
			},
		},
	}
	s := e.execute()
	fmt.Println(s)
}
