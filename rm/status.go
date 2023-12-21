package nexusrm

import (
	"context"
	"net/http"
)

const (
	restStatusReadable = "service/rest/v1/status"
	restStatusWritable = "service/rest/v1/status/writable"
)

func StatusReadableContext(ctx context.Context, rm RM) (_ bool) {
	_, resp, err := rm.Get(ctx, restStatusReadable)
	return err == nil && resp.StatusCode == http.StatusOK
}

// StatusReadable returns true if the RM instance can serve read requests
func StatusReadable(rm RM) (_ bool) {
	return StatusReadableContext(context.Background(), rm)
}

func StatusWritableContext(ctx context.Context, rm RM) (_ bool) {
	_, resp, err := rm.Get(ctx, restStatusWritable)
	return err == nil && resp.StatusCode == http.StatusOK
}

// StatusWritable returns true if the RM instance can serve read requests
func StatusWritable(rm RM) (_ bool) {
	return StatusWritableContext(context.Background(), rm)
}
