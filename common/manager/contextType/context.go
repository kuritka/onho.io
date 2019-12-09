// Package contextType implements context contextType set/get functionality.
package contextType

import (
	"context"
)

// Custom context setting/getting values
const serviceTypeKey = "service-type"

var (
	// AcceptLanguage context typed contextType for storing accept language.
	AcceptLanguage = Type("accept-language")
	// Service context typed contextType for storing service name.
	Service = Type(serviceTypeKey)
	// Environment context typed contextType for storing environment information.
	Environment = Type("environment")
	// Request represents rest request type.
	Request = struct {
		ID Type
	}{
		ID: Type("x-request-id"),
	}
)

// Type models the unique contextType type
type Type string

// String returns the string representation of a Type instance.
func (t Type) String() string {
	return string(t)
}

// SetInContext maps value to the regular and grpc contextType.
func (t Type) SetInContext(ctxt context.Context, value string) context.Context {
	ctx := context.WithValue(ctxt, t, value)
	return ctx
	//return metadata.NewOutgoingContext(ctx, metadata.Pairs(string(t), value))
}

// GetFromContext extracts the value keyed by this type instance
func (t Type) GetFromContext(ctxt context.Context) string {
	//// Check grpc metadata firstly.
	//if md, ok := metadata.FromIncomingContext(ctxt); ok {
	//	if services, ok := md[serviceTypeKey]; ok {
	//		return services[0]
	//	}
	//}
	// If metadata not found, takes from the general contextType.
	if id, ok := ctxt.Value(t).(string); ok {
		return id
	}

	return "<undefined>"
}
