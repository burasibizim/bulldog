package service

type CloudService struct {
	CodebuildService *CodebuildService
	StorageService   *StorageService
}

func NewCloudService(codebuildService *CodebuildService, storageService *StorageService) *CloudService {
	return &CloudService{
		CodebuildService: codebuildService,
		StorageService:   storageService,
	}
}
