// Code generated by protoc-gen-go-binary. DO NOT EDIT.
// source: agent/agentpb/acl/acl.proto

package acl

import (
	"github.com/golang/protobuf/proto"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (msg *ACLLink) MarshalBinary() ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (msg *ACLLink) UnmarshalBinary(b []byte) error {
	return proto.Unmarshal(b, msg)
}