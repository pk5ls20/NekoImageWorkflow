// Code generated by Kitex v0.9.1. DO NOT EDIT.

package fileuploadservice

import (
	uploadClient "NekoImageWorkflowKitex/uploadClient/kitex_gen/uploadClient"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	HandleFilePreUpload(ctx context.Context, Req *uploadClient.FilePreRequest, callOptions ...callopt.Option) (r *uploadClient.FilePreResponse, err error)
	HandleFilePostUpload(ctx context.Context, Req *uploadClient.FilePostRequest, callOptions ...callopt.Option) (r *uploadClient.FilePostResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kFileUploadServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kFileUploadServiceClient struct {
	*kClient
}

func (p *kFileUploadServiceClient) HandleFilePreUpload(ctx context.Context, Req *uploadClient.FilePreRequest, callOptions ...callopt.Option) (r *uploadClient.FilePreResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleFilePreUpload(ctx, Req)
}

func (p *kFileUploadServiceClient) HandleFilePostUpload(ctx context.Context, Req *uploadClient.FilePostRequest, callOptions ...callopt.Option) (r *uploadClient.FilePostResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleFilePostUpload(ctx, Req)
}
