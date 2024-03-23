package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func parseHeaders(headers []string) []*Header {
	h := make([]*Header, 0, len(headers))

	for _, header := range headers {
		sHeader := strings.Split(header, ":")

		h = append(h, &Header{Key: sHeader[0], Value: sHeader[1]})
	}

	return h
}

func Do(url string, method string, headers []string, body string, timeout float64) (*DoResponse, error) {
	var ctx context.Context

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		defer cancel()
	} else {
		ctx = context.Background()
	}

	var (
		bodyReader io.Reader
		res        *http.Response
		err        error
		req        *http.Request
		doResp     *DoResponse
	)

	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error: creating http requset: %w", err)
	}

	for _, header := range parseHeaders(headers) {
		req.Header.Add(header.Key, header.Value)
	}

	start := time.Now()
	res, err = http.DefaultClient.Do(req)
	end := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("error: calling http url: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println("error: closing connection", err.Error())
		}
	}()

	var bytes []byte

	bytes, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error: reading response body: %w", err)
	}

	doResp = &DoResponse{
		Url:    req.URL.RawPath,
		Status: res.Status,
		Code:   res.StatusCode,
		Time:   end.Milliseconds(),
		Body:   string(bytes),
	}

	return doResp, nil
}
