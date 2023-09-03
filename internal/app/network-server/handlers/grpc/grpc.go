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
	pb "github.com/gogapopp/shortener/internal/app/network-server/handlers/grpc/proto"
	"github.com/gogapopp/shortener/internal/app/storage"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type shortenerServer struct {
	pb.UnimplementedMultiServiceServer
	storage storage.Storage
	cfg     *config.Config
}

// NewGrpcServer создаём grpc сервер
func NewGrpcServer(cfg *config.Config, storage storage.Storage) *grpc.Server {
	// создаём gRPC-сервер без зарегистрированной службы
	grpcserver := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMultiServiceServer(grpcserver, &shortenerServer{cfg: cfg, storage: storage})
	return grpcserver
}

// SaveUrl сохраняем ссылку
func (s *shortenerServer) SaveURL(ctx context.Context, in *pb.UrlSaveRequest) (*pb.UrlSaveResponse, error) {
	var response pb.UrlSaveResponse
	// проверяем является ли ссылка переданная в body валидной
	if _, err := url.ParseRequestURI(in.LongURL); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	// делаем из обычной ссылки сжатую
	shortURL := urlshortener.ShortenerURL(s.cfg.BaseAddr)
	response.ShortURL = shortURL
	// сохраняем короткую ссылку
	err := s.storage.SaveURL(in.LongURL, shortURL, "0", "0")
	if err != nil {
		// если такая ссылка уже сохранена
		if errors.Is(postgres.ErrURLExists, err) {
			shortURL = s.storage.GetShortURL(in.LongURL)
			return nil, status.Error(codes.AlreadyExists, shortURL)
		}
		return nil, status.Error(codes.Internal, "something went wrond")
	}
	return &response, nil
}

// GetUrl получаем ссылку которая соответсвует сокращённой
func (s *shortenerServer) GetURL(ctx context.Context, in *pb.UrlGetRequest) (*pb.UrlGetResponse, error) {
	var response pb.UrlGetResponse
	// получает ссылку из хранилища
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

// Ping проверяет соединение с БД
func (s *shortenerServer) Ping(ctx context.Context, in *pb.Empty) (*pb.PingResponse, error) {
	var response pb.PingResponse
	_, err := s.storage.Ping()
	if err != nil {
		return nil, status.Error(codes.NotFound, "error ping db")
	}
	response.Result = "Pong"
	return &response, nil
}

// GetUrls получаем все ссылки пользователя
func (s *shortenerServer) GetURLs(ctx context.Context, in *pb.UrlsGetRequest) (*pb.UrlsGetResponse, error) {
	var response pb.UrlsGetResponse
	// получает ссылки из хранилища
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

// DeleteUrls принимает массив идентефикаторов сокращенных строк для удаления
func (s *shortenerServer) DeleteURLs(ctx context.Context, in *pb.UrlsToDeleteRequest) (*pb.UrlsToDeleteResponse, error) {
	var response pb.UrlsToDeleteResponse
	// получаем айди
	IDs := in.UrlIDs
	reqURL := fmt.Sprintf("http://%s", s.cfg.RunAddr)
	userID := in.UserID
	// отправляем айди на обработку
	go concurrency.ProcessIDs(IDs, reqURL, s.storage, userID)
	response.Result = "OK"
	return &response, nil
}

// BatchSave принимает на вход массив структур ссылок для сокращения и в ответ массив структур сокращённых ссылки
func (s *shortenerServer) BatchSave(ctx context.Context, in *pb.BatchUrlsRequest) (*pb.BatchUrlsResponse, error) {
	var response pb.BatchUrlsResponse
	result := make([]*pb.BatchUrlsResponse_UrlsResp, 0)
	var databaseResp []models.BatchDatabaseResponse
	// начинаем проходить по реквесту
	for k := range in.BatchUrlsReq {
		// проверяем является ли переданное значение ссылкой
		_, err := url.ParseRequestURI(in.BatchUrlsReq[k].LongURL)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid request 1")
		}
		// "сжимаем" ссылку
		BatchShortURL := urlshortener.ShortenerURL(s.cfg.BaseAddr)
		// собираем данные для отправки в бд
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

// Stats возвращает все сокращённые ссылки и кол-во юзеров
func (s *shortenerServer) Stats(ctx context.Context, in *pb.Empty) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	// получаем статистику из хранилища из хранилища
	shortURLcount, userIDcount, err := s.storage.GetStats()
	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	// формируем ответ
	response.Urls = int32(shortURLcount)
	response.Users = int32(userIDcount)
	return &response, nil
}
