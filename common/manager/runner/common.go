package runner

import (
	"context"
	"fmt"
	"net"
)

type singleRunner interface {
	run(ctx context.Context) error
	String() string
}


var withAddressListenerFactory = func(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return listener, fmt.Errorf("failed to start listening on port %s: %v", port, err)
	}
	return listener, err
}

