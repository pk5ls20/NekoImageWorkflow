// Code generated by Kitex v0.9.1. DO NOT EDIT.

package fileuploadservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	clientTransform "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform"
	proto "google.golang.org/protobuf/proto"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"HandleFilePreUpload": kitex.NewMethodInfo(
		handleFilePreUploadHandler,
		newHandleFilePreUploadArgs,
		newHandleFilePreUploadResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"HandleFilePostUpload": kitex.NewMethodInfo(
		handleFilePostUploadHandler,
		newHandleFilePostUploadArgs,
		newHandleFilePostUploadResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	fileUploadServiceServiceInfo                = NewServiceInfo()
	fileUploadServiceServiceInfoForClient       = NewServiceInfoForClient()
	fileUploadServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return fileUploadServiceServiceInfo
}

// for client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return fileUploadServiceServiceInfoForStreamClient
}

// for stream client
func serviceInfoForClient() *kitex.ServiceInfo {
	return fileUploadServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "FileUploadService"
	handlerType := (*clientTransform.FileUploadService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "protoFile",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.9.1",
		Extra:           extra,
	}
	return svcInfo
}

func handleFilePreUploadHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(clientTransform.FilePreRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(clientTransform.FileUploadService).HandleFilePreUpload(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *HandleFilePreUploadArgs:
		success, err := handler.(clientTransform.FileUploadService).HandleFilePreUpload(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*HandleFilePreUploadResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newHandleFilePreUploadArgs() interface{} {
	return &HandleFilePreUploadArgs{}
}

func newHandleFilePreUploadResult() interface{} {
	return &HandleFilePreUploadResult{}
}

type HandleFilePreUploadArgs struct {
	Req *clientTransform.FilePreRequest
}

func (p *HandleFilePreUploadArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(clientTransform.FilePreRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *HandleFilePreUploadArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *HandleFilePreUploadArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *HandleFilePreUploadArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *HandleFilePreUploadArgs) Unmarshal(in []byte) error {
	msg := new(clientTransform.FilePreRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var HandleFilePreUploadArgs_Req_DEFAULT *clientTransform.FilePreRequest

func (p *HandleFilePreUploadArgs) GetReq() *clientTransform.FilePreRequest {
	if !p.IsSetReq() {
		return HandleFilePreUploadArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *HandleFilePreUploadArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *HandleFilePreUploadArgs) GetFirstArgument() interface{} {
	return p.Req
}

type HandleFilePreUploadResult struct {
	Success *clientTransform.FilePreResponse
}

var HandleFilePreUploadResult_Success_DEFAULT *clientTransform.FilePreResponse

func (p *HandleFilePreUploadResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(clientTransform.FilePreResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *HandleFilePreUploadResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *HandleFilePreUploadResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *HandleFilePreUploadResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *HandleFilePreUploadResult) Unmarshal(in []byte) error {
	msg := new(clientTransform.FilePreResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *HandleFilePreUploadResult) GetSuccess() *clientTransform.FilePreResponse {
	if !p.IsSetSuccess() {
		return HandleFilePreUploadResult_Success_DEFAULT
	}
	return p.Success
}

func (p *HandleFilePreUploadResult) SetSuccess(x interface{}) {
	p.Success = x.(*clientTransform.FilePreResponse)
}

func (p *HandleFilePreUploadResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *HandleFilePreUploadResult) GetResult() interface{} {
	return p.Success
}

func handleFilePostUploadHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(clientTransform.FilePostRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(clientTransform.FileUploadService).HandleFilePostUpload(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *HandleFilePostUploadArgs:
		success, err := handler.(clientTransform.FileUploadService).HandleFilePostUpload(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*HandleFilePostUploadResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newHandleFilePostUploadArgs() interface{} {
	return &HandleFilePostUploadArgs{}
}

func newHandleFilePostUploadResult() interface{} {
	return &HandleFilePostUploadResult{}
}

type HandleFilePostUploadArgs struct {
	Req *clientTransform.FilePostRequest
}

func (p *HandleFilePostUploadArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(clientTransform.FilePostRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *HandleFilePostUploadArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *HandleFilePostUploadArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *HandleFilePostUploadArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *HandleFilePostUploadArgs) Unmarshal(in []byte) error {
	msg := new(clientTransform.FilePostRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var HandleFilePostUploadArgs_Req_DEFAULT *clientTransform.FilePostRequest

func (p *HandleFilePostUploadArgs) GetReq() *clientTransform.FilePostRequest {
	if !p.IsSetReq() {
		return HandleFilePostUploadArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *HandleFilePostUploadArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *HandleFilePostUploadArgs) GetFirstArgument() interface{} {
	return p.Req
}

type HandleFilePostUploadResult struct {
	Success *clientTransform.FilePostResponse
}

var HandleFilePostUploadResult_Success_DEFAULT *clientTransform.FilePostResponse

func (p *HandleFilePostUploadResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(clientTransform.FilePostResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *HandleFilePostUploadResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *HandleFilePostUploadResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *HandleFilePostUploadResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *HandleFilePostUploadResult) Unmarshal(in []byte) error {
	msg := new(clientTransform.FilePostResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *HandleFilePostUploadResult) GetSuccess() *clientTransform.FilePostResponse {
	if !p.IsSetSuccess() {
		return HandleFilePostUploadResult_Success_DEFAULT
	}
	return p.Success
}

func (p *HandleFilePostUploadResult) SetSuccess(x interface{}) {
	p.Success = x.(*clientTransform.FilePostResponse)
}

func (p *HandleFilePostUploadResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *HandleFilePostUploadResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) HandleFilePreUpload(ctx context.Context, Req *clientTransform.FilePreRequest) (r *clientTransform.FilePreResponse, err error) {
	var _args HandleFilePreUploadArgs
	_args.Req = Req
	var _result HandleFilePreUploadResult
	if err = p.c.Call(ctx, "HandleFilePreUpload", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) HandleFilePostUpload(ctx context.Context, Req *clientTransform.FilePostRequest) (r *clientTransform.FilePostResponse, err error) {
	var _args HandleFilePostUploadArgs
	_args.Req = Req
	var _result HandleFilePostUploadResult
	if err = p.c.Call(ctx, "HandleFilePostUpload", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}