package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/grpc_errors"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/types"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsService.CreateHotel")
	defer span.Finish()

	h.logger.Infof("request :%-v", req)

	hotel := &models.Hotel{
		Name:          req.GetName(),
		Email:         req.GetEmail(),
		Country:       req.GetCountry(),
		City:          req.GetCity(),
		Description:   req.GetDescription(),
		Image:         types.NullString{String: req.GetImage(), Valid: true},
		Photos:        req.Photos,
		CommentsCount: int(req.CommentsCount),
		Latitude:      types.NullFloat64{Float64: req.GetLatitude(), Valid: true},
		Longitude:     types.NullFloat64{Float64: req.GetLongitude(), Valid: true},
		Location:      req.Location,
		Rating:        req.GetRating(),
	}

	if err := h.validate.StructCtx(ctx, hotel); err != nil {
		h.logger.Errorf("validate.StructCtx: %v", err)
		return nil, grpc_errors.ErrorResponse(err, err.Error())
	}

	createdHotel, err := h.hotelsUC.CreateHotel(ctx, hotel)
	if err != nil {
		h.logger.Errorf("hotelsUC.CreateHotel: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "userUC.GetByID")
	}

	h.logger.Infof("CREATED HOTEL: %-v", createdHotel)
	h.logger.Infof("CREATED HOTEL createdHotel.ToProto(): %-v", createdHotel.ToProto())

	return &hotelsService.CreateHotelRes{Hotel: createdHotel.ToProto()}, nil
}

func (h *HotelsService) UpdateHotel(ctx context.Context, req *hotelsService.UpdateHotelReq) (*hotelsService.UpdateHotelRes, error) {
	panic("implement me")
}

func (h *HotelsService) GetHotelByID(ctx context.Context, req *hotelsService.GetByIDReq) (*hotelsService.GetByIDRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsService.GetHotelByID")
	defer span.Finish()

	hotelUUID, err := uuid.FromString(req.GetHotelID())
	if err != nil {
		h.logger.Errorf("uuid.FromString: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "uuid.FromString")
	}

	hotel, err := h.hotelsUC.GetHotelByID(ctx, hotelUUID)
	if err != nil {
		h.logger.Errorf("hotelsUC.GetHotelByID: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "hotelsUC.GetHotelByID")
	}

	return &hotelsService.GetByIDRes{Hotel: hotel.ToProto()}, nil
}
