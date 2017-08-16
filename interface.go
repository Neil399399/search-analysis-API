package main

import (
	"net"
)

type StoppableListener struct {
	*net.TCPListener          //Wrapped listener
	stop             chan int //Channel used only to indicate listener should shutdown
}
