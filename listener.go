package main

import "net"

func newListener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
