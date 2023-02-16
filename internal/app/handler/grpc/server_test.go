package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	pb "github.com/apolsh/yapr-url-shortener/internal/app/handler/grpc/proto"
	"github.com/apolsh/yapr-url-shortener/internal/app/mocks"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServerSuite struct {
	suite.Suite
	shorts *mocks.MockURLShortenerService
	server *grpc.Server
	ctrl   *gomock.Controller
	client pb.URLShortenerClient
}

func TestGRPCServer(t *testing.T) {
	suite.Run(t, new(GRPCServerSuite))
}

var (
	cryptoProvider = crypto.NewAESCryptoProvider("secret")
	shortURL1      = "http://shorturl1.com/123"
	shortURL2      = "http://shorturl2.com/456"
	longURL1       = "http://longurl1.com"
	longURL2       = "http://longurl2.com"
)

func (s *GRPCServerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.ctrl = ctrl
	s.shorts = mocks.NewMockURLShortenerService(ctrl)
	server, starter := StartGRPCServer(":3333", s.shorts, cryptoProvider, &net.IPNet{IP: net.IPv4(0, 0, 0, 0)})
	s.server = server
	go starter()
	conn, err := grpc.Dial(":3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	s.client = pb.NewURLShortenerClient(conn)
}

func (s *GRPCServerSuite) TearDownTest() {
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *GRPCServerSuite) TestPingDB() {
	s.shorts.EXPECT().PingDB(gomock.Any()).Return(true)
	res, err := s.client.PingDB(context.Background(), &emptypb.Empty{})
	fmt.Println(res)
	assert.True(s.T(), res.IsAlive)
	assert.NoError(s.T(), err)
}

func (s *GRPCServerSuite) TestGetShortenURLByIDWithSuccess() {
	id := "some_id"
	originalURL := "http://rediercted.com/url"
	s.shorts.EXPECT().GetURLByID(gomock.Any(), id).Return(originalURL, nil)
	res, err := s.client.GetShortenURLByID(context.Background(), &pb.GetShortenURLByIDRequest{UrlID: id})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), originalURL, res.GetOriginalURL())
}

func (s *GRPCServerSuite) TestGetShortenURLByIDNotFound() {
	id := "some_id"
	s.shorts.EXPECT().GetURLByID(gomock.Any(), id).Return("", repository.ErrorItemNotFound)
	res, err := s.client.GetShortenURLByID(context.Background(), &pb.GetShortenURLByIDRequest{UrlID: id})
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), codes.NotFound, status.Code(err))
}

func (s *GRPCServerSuite) TestGetShortenURLByIDItemDeleted() {
	id := "some_id"
	s.shorts.EXPECT().GetURLByID(gomock.Any(), id).Return("", service.ErrorItemIsDeleted)
	res, err := s.client.GetShortenURLByID(context.Background(), &pb.GetShortenURLByIDRequest{UrlID: id})
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), "item is marked as deleted", status.Convert(err).Message())
}

func (s *GRPCServerSuite) TestGetShortenURLsByUserSomeFound() {
	s.shorts.EXPECT().GetURLsByOwnerID(gomock.Any(), gomock.Any()).Return([]dto.URLPair{
		{ShortURL: shortURL1, OriginalURL: longURL1},
		{ShortURL: shortURL2, OriginalURL: longURL2},
	}, nil)

	res, err := s.client.GetShortenURLsByUser(context.Background(), &emptypb.Empty{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(res.UrlPairs))
}

func (s *GRPCServerSuite) TestGetShortenURLsByUserSomeError() {
	s.shorts.EXPECT().GetURLsByOwnerID(gomock.Any(), gomock.Any()).Return(nil, service.ErrorItemIsDeleted)

	res, err := s.client.GetShortenURLsByUser(context.Background(), &emptypb.Empty{})
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), codes.Internal, status.Code(err))
}

func (s *GRPCServerSuite) TestSaveShortenURLNewURLSaved() {
	s.shorts.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).Return("123", nil)
	s.shorts.EXPECT().GetShortenURLFromID(gomock.Any(), "123").Return(shortURL1)

	res, err := s.client.SaveShortenURL(context.Background(), &pb.SaveShortenURLRequest{OriginalURL: longURL1})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shortURL1, res.ShortenedURL)
}

func (s *GRPCServerSuite) TestSaveShortenURLAlreadySaved() {
	s.shorts.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).Return("", repository.ErrorURLAlreadyStored)
	s.shorts.EXPECT().GetByOriginalURL(gomock.Any(), longURL1).Return(entity.ShortenedURLInfo{ID: "123"}, nil)
	s.shorts.EXPECT().GetShortenURLFromID(gomock.Any(), "123").Return(shortURL1)

	res, err := s.client.SaveShortenURL(context.Background(), &pb.SaveShortenURLRequest{OriginalURL: longURL1})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shortURL1, res.ShortenedURL)
}

func (s *GRPCServerSuite) TestSaveShortenURLsInBatchWithSuccess() {
	s.shorts.EXPECT().AddNewURLsInBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return([]dto.ShortenInBatchResponseItem{
		{CorrelationID: "1", ShortURL: shortURL1},
		{CorrelationID: "2", ShortURL: shortURL2},
	}, nil)

	req := &pb.SaveShortenURLsInBatchRequest{
		Items: []*pb.ShortenInBatchRequestItem{
			{CorrelationID: "1", OriginalURL: longURL1},
			{CorrelationID: "2", OriginalURL: longURL2},
		}}
	res, err := s.client.SaveShortenURLsInBatch(context.Background(), req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(res.Items))
}

func (s *GRPCServerSuite) TestDeleteShortenURLsInBatchWithSuccess() {
	s.shorts.EXPECT().DeleteURLsInBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	_, err := s.client.DeleteShortenURLsInBatch(context.Background(), &pb.DeleteShortenURLsInBatchRequest{
		Items: []string{"123", "456"},
	})

	assert.NoError(s.T(), err)
}

func (s *GRPCServerSuite) TestDeleteShortenURLsInBatchWithError() {
	s.shorts.EXPECT().DeleteURLsInBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.ErrorItemIsDeleted)

	_, err := s.client.DeleteShortenURLsInBatch(context.Background(), &pb.DeleteShortenURLsInBatchRequest{
		Items: []string{"123"},
	})

	assert.Equal(s.T(), codes.Internal, status.Code(err))
}
