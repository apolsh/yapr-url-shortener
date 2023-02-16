package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	pb "github.com/apolsh/yapr-url-shortener/internal/app/handler/grpc/proto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/apolsh/yapr-url-shortener/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const ownerID string = "userId"
const sessionID string = "sessionId"
const realIPHeader string = "X-Real-IP"

var log = logger.LoggerOfComponent("grpc-router")

type urlShortenerServer struct {
	pb.UnimplementedURLShortenerServer
	shortenService service.URLShortenerService
	trustedSubnet  *net.IPNet
}

// GetServerStarter возвращает grpc server и функцию, стартующую сервер
func GetServerStarter(addr string, shortenService service.URLShortenerService, cryptoProvider crypto.CryptographicProvider, trustedSubnet *net.IPNet) (*grpc.Server, func()) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(newAuthInterceptor(cryptoProvider)))

	pb.RegisterURLShortenerServer(s, &urlShortenerServer{shortenService: shortenService, trustedSubnet: trustedSubnet})

	return s, func() {
		err = s.Serve(listen)
		if err != nil {
			log.Error(err)
		}
	}
}

// PingDB проверяет работу хранилища URL
func (u *urlShortenerServer) PingDB(ctx context.Context, _ *emptypb.Empty) (*pb.PingDBResponse, error) {
	ok := u.shortenService.PingDB(ctx)
	response := &pb.PingDBResponse{IsAlive: ok}

	return response, nil
}

// GetShortenURLByID производит редирект на сохраненный ранее в хранилище URL
func (u *urlShortenerServer) GetShortenURLByID(ctx context.Context, request *pb.GetShortenURLByIDRequest) (*pb.GetShortenURLByIDResponse, error) {
	url, err := u.shortenService.GetURLByID(ctx, request.UrlID)
	if err != nil {
		log.Error(err)
		if errors.Is(repository.ErrorItemNotFound, err) {
			return nil, status.Error(codes.NotFound, "")
		}
		if errors.Is(repository.ErrorItemNotFound, err) {
			return nil, err
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &pb.GetShortenURLByIDResponse{OriginalURL: url}

	return response, nil
}

// GetShortenURLsByUser возвращает список пар (короткий + длинный) URL пользователя
func (u *urlShortenerServer) GetShortenURLsByUser(ctx context.Context, _ *emptypb.Empty) (*pb.GetShortenURLsByUserResponse, error) {
	owner := extractSingleValueFromContext(ctx, ownerID)

	urlPairs, err := u.shortenService.GetURLsByOwnerID(ctx, owner)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbURLPairs := make([]*pb.URLPair, 0, len(urlPairs))
	for _, urlPair := range urlPairs {
		pbURLPairs = append(pbURLPairs, &pb.URLPair{OriginalURL: urlPair.OriginalURL, ShortURL: urlPair.ShortURL})
	}

	response := &pb.GetShortenURLsByUserResponse{UrlPairs: pbURLPairs}
	return response, nil
}

// SaveShortenURL принимает запрос в виде простого текста, сохраняет URL в хранилище
func (u *urlShortenerServer) SaveShortenURL(ctx context.Context, request *pb.SaveShortenURLRequest) (*pb.SaveShortenURLResponse, error) {
	owner := extractSingleValueFromContext(ctx, ownerID)

	urlID, err := u.shortenService.AddNewURL(ctx, *entity.NewUnstoredShortenedURLInfo(owner, request.OriginalURL))
	if err != nil {
		log.Error(err)
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := u.shortenService.GetByOriginalURL(ctx, request.OriginalURL)
			urlID = info.GetID()
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}
	shortenURL := u.shortenService.GetShortenURLFromID(ctx, urlID)

	response := &pb.SaveShortenURLResponse{ShortenedURL: shortenURL}

	return response, nil
}

// GetAppStats получить статистику приложения
func (u *urlShortenerServer) GetAppStats(ctx context.Context, _ *emptypb.Empty) (*pb.GetAppStatsResponse, error) {
	stringIP := extractSingleValueFromContext(ctx, realIPHeader)

	ip := net.ParseIP(stringIP)
	if ip == nil {
		log.Error(fmt.Errorf("unable to parse %s to IP", stringIP))
		return nil, status.Error(codes.PermissionDenied, "")
	}

	if !u.trustedSubnet.Contains(ip) {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	statistic, err := u.shortenService.GetAppStatistic(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetAppStatsResponse{URLs: int64(statistic.URLs), Users: int64(statistic.URLs)}

	return response, nil
}

// SaveShortenURLsInBatch сохраняет сразу несколько URL в хранилище за один запрос
func (u *urlShortenerServer) SaveShortenURLsInBatch(ctx context.Context, request *pb.SaveShortenURLsInBatchRequest) (*pb.SaveShortenURLsInBatchResponse, error) {
	owner := extractSingleValueFromContext(ctx, ownerID)

	batch := make([]dto.ShortenInBatchRequestItem, 0, len(request.Items))

	for _, item := range request.Items {
		batch = append(batch, dto.ShortenInBatchRequestItem{CorrelationID: item.CorrelationID, OriginalURL: item.OriginalURL})
	}

	saveBatchResponse, err := u.shortenService.AddNewURLsInBatch(ctx, owner, batch)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseBatch := make([]*pb.ShortenInBatchResponseItem, 0, len(request.Items))

	for _, item := range saveBatchResponse {
		responseBatch = append(responseBatch, &pb.ShortenInBatchResponseItem{CorrelationID: item.CorrelationID, ShortURL: item.ShortURL})
	}

	response := &pb.SaveShortenURLsInBatchResponse{Items: responseBatch}

	return response, nil
}

// DeleteShortenURLsInBatch помечает URL в хранилище как удаленный
func (u *urlShortenerServer) DeleteShortenURLsInBatch(ctx context.Context, request *pb.DeleteShortenURLsInBatchRequest) (*emptypb.Empty, error) {
	owner := extractSingleValueFromContext(ctx, ownerID)

	err := u.shortenService.DeleteURLsInBatch(ctx, owner, request.Items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func extractSingleValueFromContext(ctx context.Context, key string) string {
	meta, ok := metadata.FromIncomingContext(ctx)

	var value string

	if ok {
		values := meta.Get(key)
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			value = values[0]
		}
	}
	return value
}

func newAuthInterceptor(cryptoProvider crypto.CryptographicProvider) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			fillWithNewUserAndToken(&ctx, md, cryptoProvider)
			return handler(ctx, req)
		}

		values := md.Get(sessionID)
		if len(values) == 0 {
			fillWithNewUserAndToken(&ctx, md, cryptoProvider)

			return handler(ctx, req)
		}

		token := values[0]

		if len(token) == 0 {
			fillWithNewUserAndToken(&ctx, md, cryptoProvider)
			return handler(ctx, req)
		}

		sessionIDBytes, err := hex.DecodeString(token)
		if err != nil {
			fillWithNewUserAndToken(&ctx, md, cryptoProvider)
			return handler(ctx, req)
		} else {
			owner, err := cryptoProvider.Decrypt(sessionIDBytes)
			if err != nil {
				fillWithNewUserAndToken(&ctx, md, cryptoProvider)
				return handler(ctx, req)
			}
			md.Append(ownerID, owner)
			header := metadata.Pairs(sessionID, token)
			err = grpc.SendHeader(ctx, header)
			if err != nil {
				log.Error(err)
			}
			metadata.NewIncomingContext(ctx, md)
		}
		return handler(ctx, req)
	}
}

func fillWithNewUserAndToken(ctx *context.Context, md metadata.MD, cryptoProvider crypto.CryptographicProvider) {
	userUUID := uuid.New()
	token := cryptoProvider.Encrypt(userUUID[:])
	userID := userUUID.String()
	md.Append(ownerID, userID)

	header := metadata.Pairs(sessionID, token)
	err := grpc.SendHeader(*ctx, header)
	if err != nil {
		log.Error(err)
	}

	metadata.NewIncomingContext(*ctx, md)
}
