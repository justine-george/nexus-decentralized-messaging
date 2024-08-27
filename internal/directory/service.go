package directory

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/justine-george/nexus-decentralized-messaging/proto"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedDirectoryServiceServer
	peers map[string]*pb.PeerInfo
	mu    sync.RWMutex
}

func NewService() *Service {
	return &Service{
		peers: make(map[string]*pb.PeerInfo),
	}
}

func (s *Service) RegisterPeer(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers[req.Id] = &pb.PeerInfo{Id: req.Id, Address: req.Address}
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *Service) GetPeers(ctx context.Context, req *pb.GetPeersRequest) (*pb.GetPeersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var peerList []*pb.PeerInfo
	for _, peer := range s.peers {
		peerList = append(peerList, peer)
	}
	return &pb.GetPeersResponse{Peers: peerList}, nil
}

func (s *Service) Serve() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDirectoryServiceServer(grpcServer, s)
	log.Println("Directory service started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
