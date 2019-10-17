package kncloudevents

import (
	nethttp "net/http"

	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
)

type ConnectionArgs struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

func NewDefaultClient(target ...string) (cloudevents.Client, error) {
	tOpts := []http.Option{cloudevents.WithBinaryEncoding()}
	if len(target) > 0 && target[0] != "" {
		tOpts = append(tOpts, cloudevents.WithTarget(target[0]))
	}

	// Make an http transport for the CloudEvents client.
	t, err := cloudevents.NewHTTPTransport(tOpts...)
	if err != nil {
		return nil, err
	}
	return NewDefaultClientGivenHttpTransport(t)
}

// NewDefaultClientGivenHttpTransport creates a new CloudEvents client using the provided cloudevents HTTP
// transport. Note that it does modify the provided cloudevents HTTP Transport by different connnection options
// to its Client, in case they are specified.
func NewDefaultClientGivenHttpTransport(t *cloudevents.HTTPTransport, connectionArgs ...ConnectionArgs) (cloudevents.Client, error) {
	// Add connection options to the underlying transport.
	var transport = nethttp.DefaultTransport
	if len(connectionArgs) > 0 {
		httpTransport := transport.(*nethttp.Transport)
		httpTransport.MaxIdleConns = connectionArgs[0].MaxIdleConns
		httpTransport.MaxIdleConnsPerHost = connectionArgs[0].MaxIdleConnsPerHost
	}
	t.Client = &nethttp.Client{
		Transport: transport,
	}

	// Use the transport to make a new CloudEvents client.
	c, err := cloudevents.NewClient(t,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)

	if err != nil {
		return nil, err
	}
	return c, nil
}
