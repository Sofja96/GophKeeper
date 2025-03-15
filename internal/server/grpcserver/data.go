package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sofja96/GophKeeper.git/internal/models"
	"github.com/Sofja96/GophKeeper.git/internal/server/utils"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// CreateData создает новые данные для пользователя.
// Принимает запрос на создание данных и возвращает ответ с ID созданных данных.
func (s *gophKeeperServer) CreateData(ctx context.Context, req *proto.CreateDataRequest) (*proto.CreateDataResponse, error) {
	userName, ok := ctx.Value(models.ContextKeyUser).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user authentication")
	}

	userID, err := s.server.GetService().GetUserIDByUsername(ctx, userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user ID: %v", err)
	}

	dataType, err := models.GetModelType(req.DataType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid data type :%v", err)
	}

	data := &models.Data{
		UserID:      userID,
		DataType:    dataType,
		DataContent: req.DataContent,
		Metadata:    req.Metadata.AsMap(),
		FileName:    req.FileName,
	}

	dataId, err := s.server.GetService().CreateData(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create data: %v", err)
	}

	return &proto.CreateDataResponse{
		Message: "Data successfully created",
		DataId:  dataId,
	}, nil

}

// GetAllData получает все данные для текущего пользователя.
// Возвращает список всех данных пользователя.
func (s *gophKeeperServer) GetAllData(ctx context.Context, _ *proto.GetAllDataRequest) (*proto.GetAllDataResponse, error) {
	userName, ok := ctx.Value(models.ContextKeyUser).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user authentication")
	}

	userID, err := s.server.GetService().GetUserIDByUsername(ctx, userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user ID: %v", err)
	}

	data, err := s.server.GetService().GetData(ctx, userID)
	if err != nil {
		if errors.Is(err, utils.ErrUserDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "not found data for user %s", userName)
		}
		return nil, status.Errorf(codes.Internal, "failed to get data: %v", err)
	}

	responseData := make([]*proto.DataItem, 0, len(data))
	for i := range data {
		item := &data[i]

		protoDataType, err := models.ConvertModelDataTypeToProto(item.DataType)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert data type: %v", err)
		}

		protoMetadata, err := models.ConvertJSONBToStruct(item.Metadata)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert metadata: %v", err)
		}

		responseData = append(responseData, &proto.DataItem{
			DataId:      item.ID,
			DataType:    protoDataType,
			DataContent: item.DataContent,
			Metadata:    protoMetadata,
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.GetAllDataResponse{
		Data: responseData,
	}, nil

}

// DeleteData удаляет данные с указанным ID для текущего пользователя.
// Принимает ID данных для удаления и возвращает сообщение о результате.
func (s *gophKeeperServer) DeleteData(ctx context.Context, req *proto.DeleteDataRequest) (*proto.DeleteDataResponse, error) {
	userName, ok := ctx.Value(models.ContextKeyUser).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user authentication")
	}

	userID, err := s.server.GetService().GetUserIDByUsername(ctx, userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user ID: %v", err)
	}

	_, err = s.server.GetService().DeleteData(ctx, req.DataId, userID)
	if err != nil {
		if errors.Is(err, utils.ErrUserDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "данные с ID %d не найдены", req.DataId)
		}
		return nil, status.Errorf(codes.Internal, "failed delete data with ID %d", req.DataId)
	}

	return &proto.DeleteDataResponse{
		Message: fmt.Sprintf("Данные с ID %d успешно удалены", req.DataId),
	}, nil

}

// UpdateData обновляет данные с указанным ID для текущего пользователя.
// Принимает запрос на обновление данных и возвращает сообщение о результате.
func (s *gophKeeperServer) UpdateData(ctx context.Context, req *proto.UpdateDataRequest) (*proto.UpdateDataResponse, error) {
	userName, ok := ctx.Value(models.ContextKeyUser).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user authentication")
	}

	userID, err := s.server.GetService().GetUserIDByUsername(ctx, userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user ID: %v", err)
	}

	data := &models.Data{
		UserID:      userID,
		DataContent: req.DataContent,
		Metadata:    req.Metadata.AsMap(),
		FileName:    req.FileName,
		ID:          req.DataId,
	}

	err = s.server.GetService().UpdateData(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update data: %v", err)
	}

	return &proto.UpdateDataResponse{
		Message: "Data successfully updated",
	}, nil

}
