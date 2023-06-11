package model

import "net/http"

const (
	HeaderRequestId      = "X-Request-ID"
	HeaderVersionId      = "X-Version-ID"
	HeaderEtagServer     = "ETag"
	HeaderEtagClient     = "If-None-Match"
	HeaderForwarded      = "X-Forwarded-For"
	HeaderForwardedProto = "X-Forwarded-Proto"
	HeaderRealIP         = "X-Real-IP"
	HeaderToken          = "token"
	HeaderCsrfToken      = "X-CSRF-Token"
	HeaderBearer         = "BEARER"
	HeaderAuthorization  = "Authorization"

	Status          = "status"
	StatusOK        = "OK"
	StatusFail      = "FAIL"
	StatusUnhealthy = "UNHEALTHY"

	ClientDir = "client"

	ApiURLSuffix = "/api"
)

type Response struct {
	StatusCode    int
	RequestId     string
	Etag          string
	ServerVersion string
	Header        http.Header
}

type Client struct {
	URL        string
	ApiURL     string
	HTTPClient *http.Client
	AuthToken  string
	AuthType   string
	HTTPHeader map[string]string
}
