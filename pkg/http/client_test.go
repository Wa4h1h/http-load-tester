package http

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestParseHeaders(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		input  []string
		output []*Header
	}{
		{"empty input then empty output", make([]string, 0), make([]*Header, 0)},
		{
			"list of header strings should return list of Header objects",
			[]string{"k1:v1", "k2:v2"},
			[]*Header{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}},
		},
		{
			"list of header strings and malformed strings should only return list of correct Header objects",
			[]string{"k1:v1", "k2:v2", "malformed"},
			[]*Header{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}},
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			hs := parseHeaders(row.input)

			assert.Equal(t, row.output, hs)
		})
	}
}

func TestDo(t *testing.T) {
	t.Parallel()

	go func() {
		http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		http.HandleFunc("GET /timeout", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(5) * time.Second)

			w.WriteHeader(http.StatusOK)
		})

		http.ListenAndServe(":8000", nil)
	}()

	time.Sleep(2 * time.Second)

	type doInput struct {
		url     string
		method  string
		headers []string
		body    string
		timeout float64
	}

	testsMap := map[string][]struct {
		name        string
		input       *doInput
		output      *DoResponse
		outputError error
	}{
		"Sent": {
			{
				name: "server should response",
				input: &doInput{
					url: "http://localhost:8000/", method: "GET", timeout: 4,
				},
				output: &DoResponse{
					Code:   http.StatusOK,
					Status: "200 OK",
				},
			},
		},
		"Error": {
			{
				name: "hostname can not be resolved",
				input: &doInput{
					url: "not known", method: "GET", timeout: 4,
				},
				outputError: errors.New("error: calling http url: Get \"not%20known\": unsupported protocol scheme \"\""),
			},
			{
				name: "timeout exceeded",
				input: &doInput{
					url: "http://localhost:8000/timeout", method: "GET", timeout: 4,
				},
				outputError: errors.New("error: calling http url: Get \"http://localhost:8000/timeout\": context deadline exceeded"),
			},
		},
	}

	for _, row := range testsMap["Sent"] {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res, err := Do(row.input.url, row.input.method, row.input.headers, row.input.body, row.input.timeout)
			require.NoError(t, err)
			require.NotNil(t, res)

			assert.Equal(t, row.output.Status, res.Status)
			assert.Equal(t, row.output.Code, res.Code)
		})
	}

	for _, row := range testsMap["Error"] {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res, err := Do(row.input.url, row.input.method, row.input.headers, row.input.body, row.input.timeout)
			require.Error(t, err)
			require.Nil(t, res)

			assert.Equal(t, err.Error(), row.outputError.Error())
		})
	}
}
