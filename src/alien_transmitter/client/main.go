package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"flag"

	"google.golang.org/grpc"
	transmitter "greaterm/alien_detector/gen/go"
)

var (
	port       = flag.Int("port", 8888, "The server port")
)

func main( ) {

}