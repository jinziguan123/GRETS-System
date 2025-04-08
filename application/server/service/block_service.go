package service

var GlobalBlockService BlockService

func InitBlockService() {
	GlobalBlockService = NewBlockService()
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

type BlockService interface {
}
