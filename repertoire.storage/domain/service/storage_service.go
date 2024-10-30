package service

type StorageService interface {
}

type storageService struct{}

func NewStorageService() StorageService {
	return storageService{}
}
