package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 文件上传最大尺寸
const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	// 获取当前用户
	userId := c.GetString("userId")
	if userId == "" {
		utils.ResponseUnauthorized(c, "未获取到用户信息")
		return
	}

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ResponseBadRequest(c, "获取上传文件失败: "+err.Error())
		return
	}
	defer file.Close()

	// 检查文件大小
	if header.Size > MaxFileSize {
		utils.ResponseBadRequest(c, "文件大小超过限制")
		return
	}

	// 读取文件内容
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		utils.ResponseInternalServerError(c, "读取文件内容失败: "+err.Error())
		return
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
		utils.ResponseInternalServerError(c, "创建上传目录失败: "+err.Error())
		return
	}

	// 生成文件ID和文件名
	fileID := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileExt := filepath.Ext(header.Filename)
	fileName := fileID + fileExt
	filePath := filepath.Join(uploadDir, fileName)

	// 保存文件
	if err := ioutil.WriteFile(filePath, fileData, 0644); err != nil {
		utils.ResponseInternalServerError(c, "保存文件失败: "+err.Error())
		return
	}

	// 获取文件类型
	fileType := header.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// 调用链码记录文件信息
	_, err = blockchain.DefaultFabricClient.Invoke("CreateFile",
		fileID,
		header.Filename,
		fileType,
		strconv.FormatInt(header.Size, 10),
		fileHash,
		fileName,
		userId)
	if err != nil {
		// 删除已上传文件
		os.Remove(filePath)
		utils.ResponseInternalServerError(c, "记录文件信息失败: "+err.Error())
		return
	}

	// 返回文件信息
	utils.ResponseWithData(c, gin.H{
		"fileId":   fileID,
		"fileName": header.Filename,
		"fileSize": header.Size,
		"fileHash": fileHash,
	})
}

// GetFile 获取文件
func GetFile(c *gin.Context) {
	// 获取文件ID
	fileId := c.Param("id")
	if fileId == "" {
		utils.ResponseBadRequest(c, "无效的文件ID")
		return
	}

	// 调用链码查询文件信息
	result, err := blockchain.DefaultFabricClient.Query("QueryFile", fileId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询文件信息失败")
		return
	}

	// 检查是否找到文件
	if len(result) == 0 {
		utils.ResponseNotFound(c, "文件不存在")
		return
	}

	// 解析文件数据
	var fileInfo struct {
		FileID       string `json:"fileId"`
		FileName     string `json:"fileName"`
		FileType     string `json:"fileType"`
		FileSize     int64  `json:"fileSize"`
		FileHash     string `json:"fileHash"`
		StorageName  string `json:"storageName"`
		UploadTime   string `json:"uploadTime"`
		UploaderID   string `json:"uploaderId"`
		UploaderName string `json:"uploaderName"`
	}
	if err := json.Unmarshal(result, &fileInfo); err != nil {
		utils.ResponseInternalServerError(c, "解析文件数据失败")
		return
	}

	// 确定文件路径
	uploadDir := viper.GetString("upload.path")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	filePath := filepath.Join(uploadDir, fileInfo.StorageName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ResponseNotFound(c, "文件不存在")
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		utils.ResponseInternalServerError(c, "打开文件失败: "+err.Error())
		return
	}
	defer file.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.FileName))
	c.Header("Content-Type", fileInfo.FileType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.FileSize, 10))

	// 发送文件内容
	io.Copy(c.Writer, file)
}
