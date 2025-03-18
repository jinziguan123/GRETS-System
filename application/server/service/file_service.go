package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"grets_server/api/constants"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// 文件上传DTO
type FileUploadDTO struct {
	FileName     string `json:"fileName"`
	FileSize     int64  `json:"fileSize"`
	FileType     string `json:"fileType"`
	RelationType string `json:"relationType"`
	RelationID   string `json:"relationId"`
	Description  string `json:"description"`
}

// FileService 文件服务接口
type FileService interface {
	// UploadFile 上传文件
	UploadFile(file multipart.File, dto *FileUploadDTO) (string, error)
}

// fileService 文件服务实现
type fileService struct{}

// NewFileService 创建文件服务实例
func NewFileService() FileService {
	return &fileService{}
}

// UploadFile 上传文件
func (s *fileService) UploadFile(file multipart.File, dto *FileUploadDTO) (string, error) {
	// 设置最大文件大小
	const MaxFileSize = 10 * 1024 * 1024 // 10MB

	// 检查文件大小
	if dto.FileSize > MaxFileSize {
		return "", fmt.Errorf("文件大小超过限制")
	}

	// 读取文件内容
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 计算文件哈希值
	hash := sha256.Sum256(fileData)
	fileHash := hex.EncodeToString(hash[:])

	// 创建上传目录
	uploadDir := viper.GetString("upload.path")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 生成文件ID和文件名
	fileID := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileExt := filepath.Ext(dto.FileName)
	storedFileName := fileID + fileExt
	filePath := filepath.Join(uploadDir, storedFileName)

	// 保存文件
	if err := ioutil.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("Failed to save file: %v", err)
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 调用链码记录文件信息
	args := []string{
		fileID,
		dto.FileName,
		dto.FileType,
		strconv.FormatInt(dto.FileSize, 10),
		fileHash,
		storedFileName,
	}

	// 添加关联信息
	if dto.RelationType != "" && dto.RelationID != "" {
		args = append(args, dto.RelationType, dto.RelationID)
	}

	// 添加描述
	if dto.Description != "" {
		args = append(args, dto.Description)
	}

	contract, err := blockchain.GetContract(constants.AgencyOrganization)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("获取合约失败: %v", err))
		return "", fmt.Errorf("获取合约失败: %v", err)
	}

	// 调用链码创建文件记录
	_, err = contract.SubmitTransaction("CreateFile", args...)
	if err != nil {
		// 删除已上传文件
		os.Remove(filePath)
		log.Printf("Failed to create file record: %v", err)
		return "", fmt.Errorf("记录文件信息失败: %v", err)
	}

	return fileID, nil
}
