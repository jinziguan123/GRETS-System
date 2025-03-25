package controller

import (
	"encoding/json"
	"fmt"
	"grets_server/dao"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"grets_server/service"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 文件上传最大尺寸
const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

// 文件控制器结构体
type FileController struct {
	fileService service.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController() *FileController {
	return &FileController{
		fileService: service.NewFileService(dao.NewFileDAO()),
	}
}

// UploadFile 上传文件
func (ctrl *FileController) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ResponseBadRequest(c, "获取上传文件失败")
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := filepath.Ext(header.Filename)
	allowedExts := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
	valid := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			valid = true
			break
		}
	}
	if !valid {
		utils.ResponseBadRequest(c, "不支持的文件类型，仅支持PDF、Word和图片格式")
		return
	}

	// 获取关联的业务对象
	relationType := c.PostForm("relationType")
	relationID := c.PostForm("relationId")
	fileName := c.PostForm("fileName")
	fileType := c.PostForm("fileType")
	description := c.PostForm("description")

	// 如果没有提供文件名，则使用原始文件名
	if fileName == "" {
		fileName = header.Filename
	}

	// 调用服务上传文件
	fileID, err := ctrl.fileService.UploadFile(file, &service.FileUploadDTO{
		FileName:     fileName,
		FileSize:     header.Size,
		FileType:     fileType,
		RelationType: relationType,
		RelationID:   relationID,
		Description:  description,
	})
	if err != nil {
		utils.ResponseInternalServerError(c, err.Error())
		return
	}

	// 返回上传结果
	utils.ResponseSuccess(c, gin.H{
		"fileId":  fileID,
		"message": "文件上传成功",
	})
}

// 创建控制器实例
var File = NewFileController()

// 为兼容现有路由，提供这些函数
func UploadFile(c *gin.Context) {
	File.UploadFile(c)
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

// 创建全局文件控制器实例
var GlobalFileController *FileController

// 初始化文件控制器
func InitFileController() {
	GlobalFileController = NewFileController()
}
