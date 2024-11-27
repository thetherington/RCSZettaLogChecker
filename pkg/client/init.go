package client

import (
	"net/http"
	"time"
)

type Options struct {
	Username string
	Password string
	Secret   string
}

type ZettaIface interface {
	GetUnmarshalJson(requestURL string, v any) error
	GetRawPayload(requestURL string) ([]byte, error)
}

type api struct {
	Options Options
	Client  http.Client
}

func New(options Options) ZettaIface {
	return api{
		Options: options,
		Client: http.Client{
			Timeout: 10 * time.Second,
			Transport: &ZettaTransport{
				transport: http.DefaultTransport,
				username:  options.Username,
				password:  options.Password,
				secret:    options.Secret,
			},
		},
	}
}
