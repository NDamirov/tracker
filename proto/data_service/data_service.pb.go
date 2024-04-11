// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/data_service/data_service.proto

package data_service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Task struct {
	TaskId               int64                `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	AuthorId             int64                `protobuf:"varint,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Description          string               `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Status               string               `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetTaskId() int64 {
	if m != nil {
		return m.TaskId
	}
	return 0
}

func (m *Task) GetAuthorId() int64 {
	if m != nil {
		return m.AuthorId
	}
	return 0
}

func (m *Task) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Task) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Task) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

type UpdateTaskRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Task                 *Task    `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateTaskRequest) Reset()         { *m = UpdateTaskRequest{} }
func (m *UpdateTaskRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateTaskRequest) ProtoMessage()    {}
func (*UpdateTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{1}
}

func (m *UpdateTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateTaskRequest.Unmarshal(m, b)
}
func (m *UpdateTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateTaskRequest.Marshal(b, m, deterministic)
}
func (m *UpdateTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateTaskRequest.Merge(m, src)
}
func (m *UpdateTaskRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateTaskRequest.Size(m)
}
func (m *UpdateTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateTaskRequest proto.InternalMessageInfo

func (m *UpdateTaskRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *UpdateTaskRequest) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

type DeleteTaskRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TaskId               int64    `protobuf:"varint,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteTaskRequest) Reset()         { *m = DeleteTaskRequest{} }
func (m *DeleteTaskRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteTaskRequest) ProtoMessage()    {}
func (*DeleteTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{2}
}

func (m *DeleteTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteTaskRequest.Unmarshal(m, b)
}
func (m *DeleteTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteTaskRequest.Marshal(b, m, deterministic)
}
func (m *DeleteTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteTaskRequest.Merge(m, src)
}
func (m *DeleteTaskRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteTaskRequest.Size(m)
}
func (m *DeleteTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteTaskRequest proto.InternalMessageInfo

func (m *DeleteTaskRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *DeleteTaskRequest) GetTaskId() int64 {
	if m != nil {
		return m.TaskId
	}
	return 0
}

type Error struct {
	StatusCode           int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{3}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetTaskRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TaskId               int64    `protobuf:"varint,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTaskRequest) Reset()         { *m = GetTaskRequest{} }
func (m *GetTaskRequest) String() string { return proto.CompactTextString(m) }
func (*GetTaskRequest) ProtoMessage()    {}
func (*GetTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{4}
}

func (m *GetTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTaskRequest.Unmarshal(m, b)
}
func (m *GetTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTaskRequest.Marshal(b, m, deterministic)
}
func (m *GetTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTaskRequest.Merge(m, src)
}
func (m *GetTaskRequest) XXX_Size() int {
	return xxx_messageInfo_GetTaskRequest.Size(m)
}
func (m *GetTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTaskRequest proto.InternalMessageInfo

func (m *GetTaskRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GetTaskRequest) GetTaskId() int64 {
	if m != nil {
		return m.TaskId
	}
	return 0
}

type GetTaskResponse struct {
	Error                *Error   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Task                 *Task    `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTaskResponse) Reset()         { *m = GetTaskResponse{} }
func (m *GetTaskResponse) String() string { return proto.CompactTextString(m) }
func (*GetTaskResponse) ProtoMessage()    {}
func (*GetTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{5}
}

func (m *GetTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTaskResponse.Unmarshal(m, b)
}
func (m *GetTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTaskResponse.Marshal(b, m, deterministic)
}
func (m *GetTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTaskResponse.Merge(m, src)
}
func (m *GetTaskResponse) XXX_Size() int {
	return xxx_messageInfo_GetTaskResponse.Size(m)
}
func (m *GetTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTaskResponse proto.InternalMessageInfo

func (m *GetTaskResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *GetTaskResponse) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

type GetTasksRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PageNumber           int32    `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	ResultsPerPage       int32    `protobuf:"varint,3,opt,name=results_per_page,json=resultsPerPage,proto3" json:"results_per_page,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTasksRequest) Reset()         { *m = GetTasksRequest{} }
func (m *GetTasksRequest) String() string { return proto.CompactTextString(m) }
func (*GetTasksRequest) ProtoMessage()    {}
func (*GetTasksRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{6}
}

func (m *GetTasksRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTasksRequest.Unmarshal(m, b)
}
func (m *GetTasksRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTasksRequest.Marshal(b, m, deterministic)
}
func (m *GetTasksRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTasksRequest.Merge(m, src)
}
func (m *GetTasksRequest) XXX_Size() int {
	return xxx_messageInfo_GetTasksRequest.Size(m)
}
func (m *GetTasksRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTasksRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTasksRequest proto.InternalMessageInfo

func (m *GetTasksRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GetTasksRequest) GetPageNumber() int32 {
	if m != nil {
		return m.PageNumber
	}
	return 0
}

func (m *GetTasksRequest) GetResultsPerPage() int32 {
	if m != nil {
		return m.ResultsPerPage
	}
	return 0
}

type GetTasksResponse struct {
	Error                *Error   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Tasks                []*Task  `protobuf:"bytes,2,rep,name=tasks,proto3" json:"tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTasksResponse) Reset()         { *m = GetTasksResponse{} }
func (m *GetTasksResponse) String() string { return proto.CompactTextString(m) }
func (*GetTasksResponse) ProtoMessage()    {}
func (*GetTasksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_904e753f03f3d787, []int{7}
}

func (m *GetTasksResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTasksResponse.Unmarshal(m, b)
}
func (m *GetTasksResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTasksResponse.Marshal(b, m, deterministic)
}
func (m *GetTasksResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTasksResponse.Merge(m, src)
}
func (m *GetTasksResponse) XXX_Size() int {
	return xxx_messageInfo_GetTasksResponse.Size(m)
}
func (m *GetTasksResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTasksResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTasksResponse proto.InternalMessageInfo

func (m *GetTasksResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *GetTasksResponse) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "data_service.Task")
	proto.RegisterType((*UpdateTaskRequest)(nil), "data_service.UpdateTaskRequest")
	proto.RegisterType((*DeleteTaskRequest)(nil), "data_service.DeleteTaskRequest")
	proto.RegisterType((*Error)(nil), "data_service.Error")
	proto.RegisterType((*GetTaskRequest)(nil), "data_service.GetTaskRequest")
	proto.RegisterType((*GetTaskResponse)(nil), "data_service.GetTaskResponse")
	proto.RegisterType((*GetTasksRequest)(nil), "data_service.GetTasksRequest")
	proto.RegisterType((*GetTasksResponse)(nil), "data_service.GetTasksResponse")
}

func init() {
	proto.RegisterFile("proto/data_service/data_service.proto", fileDescriptor_904e753f03f3d787)
}

var fileDescriptor_904e753f03f3d787 = []byte{
	// 507 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x4d, 0x8f, 0xd3, 0x30,
	0x14, 0xdc, 0xb4, 0x4d, 0x3f, 0x5e, 0xd0, 0xb2, 0x6b, 0x10, 0x44, 0xe1, 0xa3, 0x51, 0x24, 0x50,
	0xb8, 0xa4, 0x52, 0x39, 0x20, 0x8e, 0x74, 0xbb, 0x82, 0x0a, 0x09, 0xad, 0xa2, 0x72, 0xe1, 0x12,
	0xb9, 0xf5, 0x23, 0x44, 0xdb, 0xd6, 0xc1, 0x76, 0xf8, 0x5f, 0xfc, 0x29, 0x7e, 0x07, 0xb2, 0xd3,
	0x6c, 0x13, 0xda, 0xc2, 0xa2, 0x3d, 0xbe, 0x99, 0xe7, 0xf1, 0x78, 0x3a, 0x0d, 0xbc, 0xc8, 0x05,
	0x57, 0x7c, 0xc4, 0xa8, 0xa2, 0x89, 0x44, 0xf1, 0x23, 0x5b, 0x62, 0x63, 0x88, 0x0c, 0x4f, 0xee,
	0xd5, 0x31, 0x6f, 0x98, 0x72, 0x9e, 0xae, 0x70, 0x64, 0xb8, 0x45, 0xf1, 0x75, 0xa4, 0xb2, 0x35,
	0x4a, 0x45, 0xd7, 0x79, 0xb9, 0x1e, 0xfc, 0xb4, 0xa0, 0x33, 0xa7, 0xf2, 0x9a, 0x3c, 0x86, 0x9e,
	0xa2, 0xf2, 0x3a, 0xc9, 0x98, 0x6b, 0xf9, 0x56, 0xd8, 0x8e, 0xbb, 0x7a, 0x9c, 0x31, 0xf2, 0x04,
	0x06, 0xb4, 0x50, 0xdf, 0xb8, 0xd0, 0x54, 0xcb, 0x50, 0xfd, 0x12, 0x98, 0x31, 0xe2, 0x83, 0xc3,
	0x50, 0x2e, 0x45, 0x96, 0xab, 0x8c, 0x6f, 0xdc, 0xb6, 0x6f, 0x85, 0x83, 0xb8, 0x0e, 0x91, 0x47,
	0xd0, 0x95, 0x8a, 0xaa, 0x42, 0xba, 0x1d, 0x43, 0x6e, 0x27, 0xf2, 0x16, 0x60, 0x29, 0x90, 0x2a,
	0x64, 0x09, 0x55, 0xae, 0xed, 0x5b, 0xa1, 0x33, 0xf6, 0xa2, 0xd2, 0x6e, 0x54, 0xd9, 0x8d, 0xe6,
	0x95, 0xdd, 0x78, 0xb0, 0xdd, 0x7e, 0xa7, 0x82, 0x39, 0x9c, 0x7f, 0xce, 0x19, 0x55, 0xa8, 0x8d,
	0xc7, 0xf8, 0xbd, 0x40, 0xa9, 0xb4, 0xff, 0x42, 0xa2, 0xa8, 0xf9, 0xd7, 0xe3, 0x8c, 0x91, 0x97,
	0xd0, 0xd1, 0x2f, 0x31, 0xd6, 0x9d, 0x31, 0x89, 0x1a, 0x99, 0x19, 0x05, 0xc3, 0x07, 0x97, 0x70,
	0x3e, 0xc5, 0x15, 0xde, 0x52, 0xb5, 0x16, 0x57, 0xab, 0x1e, 0x57, 0x30, 0x01, 0xfb, 0x52, 0x08,
	0x2e, 0xc8, 0x10, 0x9c, 0xf2, 0xa9, 0xc9, 0x92, 0x33, 0x34, 0xc7, 0xed, 0x18, 0x4a, 0xe8, 0x82,
	0x33, 0x24, 0x2e, 0xf4, 0xd6, 0x28, 0x25, 0x4d, 0xd1, 0x48, 0x0c, 0xe2, 0x6a, 0x0c, 0x26, 0x70,
	0xfa, 0x1e, 0xd5, 0xdd, 0x7c, 0x30, 0xb8, 0x7f, 0xa3, 0x21, 0x73, 0xbe, 0x91, 0x48, 0x5e, 0x81,
	0x8d, 0xda, 0x9a, 0x91, 0x70, 0xc6, 0x0f, 0x9a, 0x51, 0x18, 0xd7, 0x71, 0xb9, 0x71, 0xeb, 0xd0,
	0x8a, 0x9b, 0x5b, 0xe4, 0x3f, 0xad, 0x0e, 0xc1, 0xc9, 0x69, 0x8a, 0xc9, 0xa6, 0x58, 0x2f, 0x50,
	0x18, 0x69, 0x3b, 0x06, 0x0d, 0x7d, 0x32, 0x08, 0x09, 0xe1, 0x4c, 0xa0, 0x2c, 0x56, 0x4a, 0x26,
	0x39, 0x8a, 0x44, 0x33, 0xa6, 0x51, 0x76, 0x7c, 0xba, 0xc5, 0xaf, 0x50, 0x5c, 0xe9, 0x80, 0x52,
	0x38, 0xdb, 0x5d, 0xfb, 0xff, 0xaf, 0x0b, 0xc1, 0xd6, 0xee, 0xa5, 0xdb, 0xf2, 0xdb, 0x47, 0x9e,
	0x57, 0x2e, 0x8c, 0x7f, 0xb5, 0xa0, 0xaf, 0xe7, 0x29, 0x55, 0x94, 0xbc, 0x01, 0xb8, 0x30, 0x25,
	0x34, 0x7f, 0x98, 0x03, 0xa7, 0xbc, 0x43, 0x97, 0x06, 0x27, 0x64, 0x0a, 0xb0, 0x2b, 0x2c, 0x19,
	0x36, 0x97, 0xf6, 0xaa, 0xfc, 0x17, 0x95, 0x5d, 0x41, 0xff, 0x54, 0xd9, 0xab, 0xee, 0x31, 0x95,
	0x0f, 0xd0, 0xdb, 0x46, 0x47, 0x9e, 0x36, 0x37, 0x9a, 0x95, 0xf3, 0x9e, 0x1d, 0x61, 0xcb, 0xb8,
	0x83, 0x13, 0xf2, 0x11, 0xfa, 0xd5, 0x8f, 0x40, 0x0e, 0x2f, 0x57, 0x9d, 0xf0, 0x9e, 0x1f, 0xa3,
	0x2b, 0xb1, 0xc9, 0xc3, 0x2f, 0x64, 0xff, 0xfb, 0xb6, 0xe8, 0x1a, 0xec, 0xf5, 0xef, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xe3, 0x8c, 0xf8, 0xf5, 0xfc, 0x04, 0x00, 0x00,
}