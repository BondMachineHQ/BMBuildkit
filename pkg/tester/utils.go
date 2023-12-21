package tester

import (
	"net/http"
	"testing"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func FromDir(path string) (string, error) {
	return "", nil
}

func Invoke(t *testing.T, input string, handler HandlerFunc) error {
	t.Helper()
	// err := handler.Handle(req, &resp)
	// if err != nil {
	// 	return &resp, err
	// }

	// expected, err := toObjectMap(b.Scheme, b.ExpectedOutput)
	// if err != nil {
	// 	return &resp, err
	// }

	// collected, err := toObjectMap(b.Scheme, resp.Collected)
	// if err != nil {
	// 	return &resp, err
	// }

	// assert.Equal(t, b.ExpectedDelay, resp.Delay)
	return nil
}

func DefaultTest(t *testing.T, path string, handler HandlerFunc) {
	t.Helper()
	t.Run(path, func(t *testing.T) {
		input, err := FromDir(path)
		if err != nil {
			t.Fatal(err)
		}
		err = Invoke(t, input, handler)
		if err != nil {
			t.Fatal(err)
		}
	})
}
