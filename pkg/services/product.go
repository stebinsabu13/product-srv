package services

import (
	"context"
	"net/http"

	"github.com/stebin13/product-srv/pkg/db"
	"github.com/stebin13/product-srv/pkg/models"
	"github.com/stebin13/product-srv/pkg/pb"
)

type Server struct {
	H db.Handler
	pb.ProductServiceServer
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var product models.Product
	product.Name = req.Name
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.H.DB.Create(&product).Error; err != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to create product",
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product
	if err := s.H.DB.Where("id=?", req.Id).First(&product).Error; err != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  "product doesn't exsist",
		}, nil
	}
	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}
	return &pb.FindOneResponse{
		Status: http.StatusFound,
		Data:   data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	var product models.Product
	if err := s.H.DB.Where("id=?", req.Id).Find(&product).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  "failed to find the product",
		}, nil
	}
	// if product.Stock <= 0 {
	// 	return &pb.DecreaseStockResponse{
	// 		Status: http.StatusUnprocessableEntity,
	// 		Error:  "insufficient stock",
	// 	}, nil
	// }
	var log models.StockDecreaseLog
	if err := s.H.DB.Where("order_id=?", req.OrderId).Find(&log).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "stock already decreased",
		}, nil
	}
	product.Stock -= 1
	if err := s.H.DB.Save(&product).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to decrease stock",
		}, nil
	}
	log.OrderId = req.OrderId
	log.ProductRefer = req.Id
	if err := s.H.DB.Create(&log).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to create log entry",
		}, nil
	}
	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
