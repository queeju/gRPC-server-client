package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	rnd "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	transmitter "greaterm/alien_detector/gen/go"
)

var (
	port       = flag.Int("port", 8888, "The server port")
	normalDist distuv.Normal
)

type server struct {
	transmitter.UnimplementedTransmitterServiceServer
}

func (s *server) Transmit(
	req *transmitter.Request,
	stream transmitter.TransmitterService_TransmitServer) error {
	doMath() // generate normal distribution
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(time.Second)
			uuid, freq := getMsg()
			res := &transmitter.Response{
				SessionId: uuid,
				Frequency: freq,
				Time:      timestamppb.Now(),
			}
			if err := stream.SendMsg(res); err != nil {
				return err
			}
		}
	}
} 

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	transmitterServer := &server{}
	transmitter.RegisterTransmitterServiceServer(grpcServer, transmitterServer)
	grpcServer.Serve(lis)
}

func doMath() {
	// mean from [-10, 10] interval
	mean := rand.Float64()*20 - 10
	// standard deviation from [0.3, 1.5].
	sd := rand.Float64()*1.2 + 0.3

	// Create a normal distribution with the specified mean and standard deviation
	normalDist = distuv.Normal{
		Mu:    mean,
		Sigma: sd,
		Src:   rnd.NewSource(uint64(time.Now().UnixNano())),
	}

	fmt.Println("mean:", mean)
	fmt.Println("sd:  ", sd)
}

func getMsg() (string, float64) {
	// Generate a sample from the normal distribution
	freq := normalDist.Rand()
	uuid := uuid.New().String()
	fmt.Println("freq:", freq)
	fmt.Println("uuid:", uuid)
	return uuid, freq
}
