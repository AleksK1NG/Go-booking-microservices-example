package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/proto/hotels"
)

type HotelsService struct {
	hotelsUC hotels.UseCase
	logger   logger.Logger
	validate *validator.Validate
}

func NewHotelsService(hotelsUC hotels.UseCase, logger logger.Logger, validate *validator.Validate) *HotelsService {
	return &HotelsService{hotelsUC: hotelsUC, logger: logger, validate: validate}
}

func (h *HotelsService) CreateHotel(ctx context.Context, req *hotelsService.CreateHotelReq) (*hotelsService.CreateHotelRes, error) {
	panic("implement me")
}

func (h *HotelsService) UpdateHotel(ctx context.Context, req *hotelsService.UpdateHotelReq) (*hotelsService.UpdateHotelRes, error) {
	panic("implement me")
}

func (h *HotelsService) GetHotelByID(ctx context.Context, req *hotelsService.GetByIDReq) (*hotelsService.GetByIDRes, error) {
	panic("implement me")
}
