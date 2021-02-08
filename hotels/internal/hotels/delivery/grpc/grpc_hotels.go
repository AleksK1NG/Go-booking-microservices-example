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
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/utils"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/proto/hotels"
)

// HotelsService
type HotelsService struct {
	hotelsUC hotels.UseCase
	logger   logger.Logger
	validate *validator.Validate
}

func NewHotelsService(hotelsUC hotels.UseCase, logger logger.Logger, validate *validator.Validate) *HotelsService {
	return &HotelsService{hotelsUC: hotelsUC, logger: logger, validate: validate}
}

// CreateHotel create new hotel
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
		Image:         &req.Image,
		Photos:        req.Photos,
		CommentsCount: int(req.CommentsCount),
		Latitude:      &req.Latitude,
		Longitude:     &req.Longitude,
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

// UpdateHotel update existing hotel
func (h *HotelsService) UpdateHotel(ctx context.Context, req *hotelsService.UpdateHotelReq) (*hotelsService.UpdateHotelRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsService.UpdateHotel")
	defer span.Finish()

	h.logger.Infof("request :%-v", req)

	hotelUUID, err := uuid.FromString(req.GetHotelID())
	if err != nil {
		h.logger.Errorf("uuid.FromString: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "uuid.FromString")
	}

	hotel := &models.Hotel{
		HotelID:       hotelUUID,
		Name:          req.GetName(),
		Email:         req.GetEmail(),
		Country:       req.GetCountry(),
		City:          req.GetCity(),
		Description:   req.GetDescription(),
		Image:         &req.Image,
		Photos:        req.Photos,
		CommentsCount: int(req.CommentsCount),
		Latitude:      &req.Latitude,
		Longitude:     &req.Longitude,
		Location:      req.Location,
		Rating:        req.GetRating(),
	}

	if err := h.validate.StructCtx(ctx, hotel); err != nil {
		h.logger.Errorf("validate.StructCtx: %v", err)
		return nil, grpc_errors.ErrorResponse(err, err.Error())
	}

	updatedHotel, err := h.hotelsUC.UpdateHotel(ctx, hotel)
	if err != nil {
		h.logger.Errorf("hotelsUC.UpdateHotel: %v", err)
		return nil, grpc_errors.ErrorResponse(err, err.Error())
	}

	return &hotelsService.UpdateHotelRes{Hotel: updatedHotel.ToProto()}, err
}

// GetHotelByID get hotel by uuid
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

// GetHotels
func (h *HotelsService) GetHotels(ctx context.Context, req *hotelsService.GetHotelsReq) (*hotelsService.GetHotelsRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsService.GetHotels")
	defer span.Finish()

	query := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))

	hotelsList, err := h.hotelsUC.GetHotels(ctx, query)
	if err != nil {
		h.logger.Errorf("hotelsUC.GetHotels: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "hotelsUC.GetHotels")
	}

	return &hotelsService.GetHotelsRes{
		TotalCount: int64(hotelsList.TotalCount),
		TotalPages: int64(hotelsList.TotalPages),
		Page:       int64(hotelsList.Page),
		Size:       int64(hotelsList.Size),
		HasMore:    hotelsList.HasMore,
		Hotels:     hotelsList.ToProto(),
	}, nil
}

// UploadImage
func (h *HotelsService) UploadImage(ctx context.Context, req *hotelsService.UploadImageReq) (*hotelsService.UploadImageRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsService.UploadImage")
	defer span.Finish()

	hotelUUID, err := uuid.FromString(req.GetHotelID())
	if err != nil {
		h.logger.Errorf("uuid.FromString: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "uuid.FromString")
	}

	if err := h.hotelsUC.UploadImage(ctx, &models.UploadHotelImageMsg{
		HotelID:     hotelUUID,
		Data:        req.GetData(),
		ContentType: req.GetContentType(),
	}); err != nil {
		h.logger.Errorf("hotelsUC.UploadImage: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "hotelsUC.UploadImage")
	}

	return &hotelsService.UploadImageRes{HotelID: hotelUUID.String()}, nil
}
