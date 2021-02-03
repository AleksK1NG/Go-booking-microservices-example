package repository

type ImageAWSRepository struct {
}

func NewImageAWSRepository() *ImageAWSRepository {
	return &ImageAWSRepository{}
}

func (i *ImageAWSRepository) PutObject() {
	panic("implement me")
}
