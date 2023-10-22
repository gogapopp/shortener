package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/concurrency"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/storage"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
	pb "github.com/gogapopp/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type shortenerServer struct {
	pb.UnimplementedMultiServiceServer
	storage storage.Storage
	cfg     *config.Config
}

// NewGrpcServer creating an rpc server
func NewGrpcServer(cfg *config.Config, storage storage.Storage) *grpc.Server {
	// creating an grpc server without a registered service
	grpcserver := grpc.NewServer()
	// registering the service
	pb.RegisterMultiServiceServer(grpcserver, &shortenerServer{cfg: cfg, storage: storage})
	return grpcserver
}

// SaveURL saving the link
func (s *shortenerServer) SaveURL(ctx context.Context, in *pb.UrlSaveRequest) (*pb.UrlSaveResponse, error) {
	var response pb.UrlSaveResponse
	// check whether the link passed to the body is valid
	if _, err := url.ParseRequestURI(in.LongURL); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	// making a compressed link from a regular link
	shortURL := urlshortener.ShortenerURL(s.cfg.BaseAddr)
	response.ShortURL = shortURL
	// saving a short link
	err := s.storage.SaveURL(in.LongURL, shortURL, "0", "0")
	if err != nil {
		// if such a link has already been saved
		if errors.Is(postgres.ErrURLExists, err) {
			shortURL = s.storage.GetShortURL(in.LongURL)
			return nil, status.Error(codes.AlreadyExists, shortURL)
		}
		return nil, status.Error(codes.Internal, "something went wrond")
	}
	return &response, nil
}

// GetURL we get a link that corresponds to the abbreviated
func (s *shortenerServer) GetURL(ctx context.Context, in *pb.UrlGetRequest) (*pb.UrlGetResponse, error) {
	var response pb.UrlGetResponse
	// gets a link from the repository
	isDelete, longURL, err := s.storage.GetURL(in.ShortURL, "0")
	if isDelete {
		return nil, status.Error(codes.InvalidArgument, "url was deleted")
	}
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "url not found")
	}
	response.LongURL = longURL
	return &response, nil
}

// Ping checks the connection to the database
func (s *shortenerServer) Ping(ctx context.Context, in *pb.Empty) (*pb.PingResponse, error) {
	var response pb.PingResponse
	_, err := s.storage.Ping()
	if err != nil {
		return nil, status.Error(codes.NotFound, "error ping db")
	}
	response.Result = "Pong"
	return &response, nil
}

// GetURLs getting all the user's links
func (s *shortenerServer) GetURLs(ctx context.Context, in *pb.UrlsGetRequest) (*pb.UrlsGetResponse, error) {
	var response pb.UrlsGetResponse
	// gets links from the repository
	userURLs, err := s.storage.GetUserURLs(in.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "url not found")
	}
	if len(userURLs) == 0 {
		return nil, status.Error(codes.NotFound, "no content")
	}
	result := make([]*pb.UrlsGetResponse_UserUrl, len(userURLs))
	for _, urls := range userURLs {
		result = append(result, &pb.UrlsGetResponse_UserUrl{
			OriginalURL: urls.OriginalURL,
			ShortURL:    urls.ShortURL,
		})
	}
	response.Urls = result
	return &response, nil
}

// DeleteURLs accepts an array of abbreviated string IDs to delete
func (s *shortenerServer) DeleteURLs(ctx context.Context, in *pb.UrlsToDeleteRequest) (*pb.UrlsToDeleteResponse, error) {
	var response pb.UrlsToDeleteResponse
	// getting an id
	IDs := in.UrlIDs
	reqURL := fmt.Sprintf("http://%s", s.cfg.RunAddr)
	userID := in.UserID
	// we send the ID for processing
	go concurrency.ProcessIDs(IDs, reqURL, s.storage, userID)
	response.Result = "OK"
	return &response, nil
}

// BatchSave takes as input an array of reference structures for shortening and in response an array of abbreviated reference structures
func (s *shortenerServer) BatchSave(ctx context.Context, in *pb.BatchUrlsRequest) (*pb.BatchUrlsResponse, error) {
	var response pb.BatchUrlsResponse
	result := make([]*pb.BatchUrlsResponse_UrlsResp, 0)
	var databaseResp []models.BatchDatabaseResponse
	// we are starting to go through the request
	for k := range in.BatchUrlsReq {
		// check whether the passed value is a reference
		_, err := url.ParseRequestURI(in.BatchUrlsReq[k].LongURL)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid request 1")
		}
		// "compress" the link
		BatchShortURL := urlshortener.ShortenerURL(s.cfg.BaseAddr)
		// collecting data to send to the database
		databaseResp = append(databaseResp, models.BatchDatabaseResponse{
			ShortURL:      BatchShortURL,
			OriginalURL:   in.BatchUrlsReq[k].LongURL,
			CorrelationID: in.BatchUrlsReq[k].CorrelationID,
			UserID:        in.UserID,
		})
		result = append(result, &pb.BatchUrlsResponse_UrlsResp{
			CorrelationID: in.BatchUrlsReq[k].CorrelationID,
			ShortURL:      BatchShortURL,
		})

	}
	err := s.storage.BatchInsertURL(databaseResp, in.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request 2")
	}
	response.BatchUrlsResp = result
	return &response, nil
}

// Stats returns all abbreviated links and the number of users
func (s *shortenerServer) Stats(ctx context.Context, in *pb.Empty) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	// getting statistics from the repository from the repository
	shortURLcount, userIDcount, err := s.storage.GetStats()
	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	// forming a response
	response.Urls = int32(shortURLcount)
	response.Users = int32(userIDcount)
	return &response, nil
}
