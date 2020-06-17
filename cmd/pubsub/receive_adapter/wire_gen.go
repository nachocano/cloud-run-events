// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"github.com/google/knative-gcp/pkg/pubsub/adapter"
	"github.com/google/knative-gcp/pkg/utils/clients"
)

// Injectors from wire.go:

func InitializeAdapter(ctx context.Context, projectID clients.ProjectID, subscription adapter.SubscriptionID, maxConnsPerHost clients.MaxConnsPerHost, name adapter.Name, namespace adapter.Namespace, resourceGroup adapter.ResourceGroup, adapterType adapter.AdapterType, sinkURI adapter.SinkURI, transformerURI adapter.TransformerURI, extensions map[string]string) (*adapter.Adapter, error) {
	client, err := clients.NewPubsubClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	pubsubSubscription := adapter.NewPubSubSubscription(ctx, client, subscription)
	httpClient := clients.NewHTTPClient(ctx, maxConnsPerHost)
	statsReporter, err := adapter.NewStatsReporter(name, namespace, resourceGroup)
	if err != nil {
		return nil, err
	}
	adapterAdapter := adapter.NewAdapter(ctx, pubsubSubscription, httpClient, statsReporter, sinkURI, transformerURI, adapterType, extensions)
	return adapterAdapter, nil
}
