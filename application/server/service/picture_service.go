package service

import (
	pictureDTO "grets_server/dto/picture_dto"
)

// GlobalPictureService 全局图片服务实例
var GlobalPictureService PictureService

// InitPictureService 初始化图片服务
func InitPictureService() {
	GlobalPictureService = NewPictureService()
}

// PictureService 图片服务接口
type PictureService interface {
	UploadPicture(req *pictureDTO.UploadPictureDTO) error
}

// pictureService 图片服务实现
type pictureService struct {
}

// NewPictureService 创建图片服务实例
func NewPictureService() PictureService {
	return &pictureService{}
}

// UploadPicture 上传图片
func (s *pictureService) UploadPicture(req *pictureDTO.UploadPictureDTO) error {
	return nil
}
