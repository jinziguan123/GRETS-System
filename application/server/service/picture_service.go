package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// GlobalPictureService 全局图片服务实例
var GlobalPictureService PictureService

// InitPictureService 初始化图片服务
func InitPictureService() {
	GlobalPictureService = NewPictureService()
}

// PictureService 图片服务接口
type PictureService interface {
	UploadPicture(file *multipart.File) (string, error)
	getLskyProToken() (string, error)
	uploadPicture(token string, file *multipart.File) (string, error)
}

// pictureService 图片服务实现
type pictureService struct {
}

// NewPictureService 创建图片服务实例
func NewPictureService() PictureService {
	return &pictureService{}
}

// UploadPicture 上传图片
func (s *pictureService) UploadPicture(file *multipart.File) (string, error) {
	// 首先调用lsky-pro的接口获取token
	token, err := s.getLskyProToken()
	if err != nil {
		return "", err
	}
	// 然后调用lsky-pro的接口上传图片
	url, err := s.uploadPicture(token, file)
	if err != nil {
		return "", err
	}
	return url, nil
}

// getLskyProToken 获取lsky-pro的token
func (s *pictureService) getLskyProToken() (string, error) {
	// 调用lsky-pro的接口获取token
	// 地址为 http://localhost:8089/api/v1/tokens POST请求
	// 请求体为：
	// {
	// 	"email": "18917950960@163.com",
	// 	"password": "Jinziguan123"
	// }
	// 返回体为：
	// {
	// 	"status": true,
	// 	"message": "success",
	// 	"data": {
	// 		"token": "1|5p2z39ziDXZrBHcxhSNyweRaSCXCkhCGNmzyOAe9"
	// 	}
	// }
	// 使用http.Post请求
	requestBody := map[string]string{
		"email":    "18917950960@163.com",
		"password": "Jinziguan123",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", "http://localhost:8089/api/v1/tokens", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}
	if responseBody["status"].(bool) {
		return responseBody["data"].(map[string]interface{})["token"].(string), nil
	}
	return "", fmt.Errorf("获取token失败: %v", responseBody["message"])
}

// uploadPicture 上传图片
func (s *pictureService) uploadPicture(token string, file *multipart.File) (string, error) {
	// 调用lsky-pro的接口上传图片
	// 地址为 http://localhost:8089/api/v1/upload POST请求
	// 请求体为：
	// {
	// 	"file": "图片文件"
	// }
	// 请求头为：
	// Authorization: Bearer 1|5p2z39ziDXZrBHcxhSNyweRaSCXCkhCGNmzyOAe9
	// 返回体为：
	// {
	// 	"status": true,
	// 	"message": "上传成功",
	// 	"data": {
	// 		"key": "CA9PHo",
	// 		"name": "6808ff391b3ab.jpg",
	// 		"pathname": "2025/04/23/6808ff391b3ab.jpg",
	// 		"origin_name": "虹夏妈妈.jpg",
	// 		"size": 30.38671875,
	// 		"mimetype": "image/jpeg",
	// 		"extension": "jpg",
	// 		"md5": "be8a38fb0ac56354df0063eb2b9ac9a8",
	// 		"sha1": "be1b099ef2945b50e00ce40a1e5a292d880913b6",
	// 		"links": {
	// 			"url": "http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg",
	// 			"html": "&lt;img src=\"http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg\" alt=\"虹夏妈妈.jpg\" title=\"虹夏妈妈.jpg\" /&gt;",
	// 			"bbcode": "[img]http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg[/img]",
	// 			"markdown": "![虹夏妈妈.jpg](http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg)",
	// 			"markdown_with_link": "[![虹夏妈妈.jpg](http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg)](http://localhost:8089/i/2025/04/23/6808ff391b3ab.jpg)",
	// 			"thumbnail_url": "http://localhost:8089/thumbnails/be8a38fb0ac56354df0063eb2b9ac9a8.png"
	// 		}
	// 	}
	// }
	// 使用http.Post请求
	// 创建一个buffer用于构建multipart请求
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 创建文件表单字段
	part, err := writer.CreateFormFile("file", "file.jpg") // 文件名可以根据需要调整
	if err != nil {
		return "", err
	}

	// 复制文件内容到表单字段
	_, err = io.Copy(part, *file)
	if err != nil {
		return "", err
	}

	// 关闭writer以完成请求内容
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// 创建请求
	request, err := http.NewRequest("POST", "http://localhost:8089/api/v1/upload", &requestBody)
	if err != nil {
		return "", err
	}

	// 设置请求头
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// 处理响应
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		return "", err
	}

	if responseBody["status"].(bool) {
		return responseBody["data"].(map[string]interface{})["links"].(map[string]interface{})["url"].(string), nil
	}

	return "", fmt.Errorf("上传图片失败: %v", responseBody["message"])
}
