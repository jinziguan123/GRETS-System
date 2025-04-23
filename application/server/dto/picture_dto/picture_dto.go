package picturedto

// UploadPictureDTO 上传图片请求参数
type UploadPictureDTO struct {
	PictureName string `json:"picture_name"`
	PictureURL  string `json:"picture_url"`
}
