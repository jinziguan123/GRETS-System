package dao

import (
	"fmt"
	"grets_server/db"
	"time"

	"gorm.io/gorm"
)

// 文件模型已经移动到db.File中，这里删除重复定义

// FileDAO 文件数据访问对象
type FileDAO struct {
	mysqlDB *gorm.DB
	boltDB  *db.BoltDB
}

// 创建新的FileDAO实例
func NewFileDAO() *FileDAO {
	return &FileDAO{
		mysqlDB: db.GlobalMysql,
		boltDB:  db.GlobalBoltDB,
	}
}

// SaveFile 保存文件信息到数据库
func (dao *FileDAO) SaveFile(file *db.File) error {
	// 保存到MySQL数据库
	if err := dao.mysqlDB.Create(file).Error; err != nil {
		return fmt.Errorf("保存文件信息失败: %v", err)
	}

	// 保存状态到BoltDB
	fileState := map[string]interface{}{
		"id":            file.ID,
		"resource_id":   file.ResourceID,
		"resource_type": file.ResourceType,
		"uploaded_at":   time.Now(),
	}
	if err := dao.boltDB.Put("file_states", file.ID, fileState); err != nil {
		return fmt.Errorf("保存文件状态失败: %v", err)
	}

	return nil
}

// GetFileByID 根据ID获取文件
func (dao *FileDAO) GetFileByID(id string) (*db.File, error) {
	var file db.File
	if err := dao.mysqlDB.First(&file, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("根据ID查询文件失败: %v", err)
	}
	return &file, nil
}

// GetFilesByResourceID 根据资源ID获取文件列表
func (dao *FileDAO) GetFilesByResourceID(resourceID, resourceType string) ([]*db.File, error) {
	var files []*db.File
	query := dao.mysqlDB.Model(&db.File{})

	if resourceID != "" {
		query = query.Where("resource_id = ?", resourceID)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, fmt.Errorf("查询文件列表失败: %v", err)
	}

	return files, nil
}

// UpdateFileIPFS 更新文件的IPFS哈希
func (dao *FileDAO) UpdateFileIPFS(id, ipfsHash string) error {
	if err := dao.mysqlDB.Model(&db.File{}).Where("id = ?", id).Update("ipfs_hash", ipfsHash).Error; err != nil {
		return fmt.Errorf("更新文件IPFS哈希失败: %v", err)
	}
	return nil
}

// DeleteFile 删除文件记录
func (dao *FileDAO) DeleteFile(id string) error {
	if err := dao.mysqlDB.Delete(&db.File{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %v", err)
	}

	// 删除BoltDB中的状态
	if err := dao.boltDB.Delete("file_states", id); err != nil {
		return fmt.Errorf("删除文件状态失败: %v", err)
	}

	return nil
}
