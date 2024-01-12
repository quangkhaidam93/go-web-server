package chat

import "context"

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	return &Message{Body: "Hello from the server!"}, nil
}

func (s *Server) mustEmbedUnimplementedChatServiceServer() {}
