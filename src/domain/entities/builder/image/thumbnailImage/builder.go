package thumbnailImageBuilder

import (
	"go-gallery/src/commons/exception"
	validators "go-gallery/src/commons/utils/validations"
	thumbnailImageEntity "go-gallery/src/domain/entities/image/thumbnailImage"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
)

type ThumbnailImageBuilder struct {
	id          *string
	imageID     *string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
	imageSize   string
}

func NewThumbnailImageBuilder() *ThumbnailImageBuilder {
	return &ThumbnailImageBuilder{}
}

func (b *ThumbnailImageBuilder) FromDTO(dto *thumbnailImageDTO.ThumbnailImageDTO) *ThumbnailImageBuilder {
	id := dto.Id
	imageID := dto.ImageID
	b.id = id
	b.imageID = imageID
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = dto.ContentFile
	b.owner = dto.Owner
	b.size = dto.Size
	b.imageSize = dto.ImageSize

	return b
}

func (b *ThumbnailImageBuilder) SetId(id *string) *ThumbnailImageBuilder {
	b.id = id
	return b
}

func (b *ThumbnailImageBuilder) SetName(name string) *ThumbnailImageBuilder {
	b.name = name
	return b
}

func (b *ThumbnailImageBuilder) SetOwner(owner string) *ThumbnailImageBuilder {
	b.owner = owner
	return b
}
func (b *ThumbnailImageBuilder) SetExtension(extension string) *ThumbnailImageBuilder {
	b.extension = extension
	return b
}

func (b *ThumbnailImageBuilder) SetContentFile(contentFile string) *ThumbnailImageBuilder {
	b.contentFile = contentFile
	return b
}

func (b *ThumbnailImageBuilder) SetSize(size string) *ThumbnailImageBuilder {
	b.size = size
	return b
}

func (b *ThumbnailImageBuilder) SetImageID(imageID *string) *ThumbnailImageBuilder {
	b.imageID = imageID
	return b
}

func (b *ThumbnailImageBuilder) SetImageSize(imageSize string) *ThumbnailImageBuilder {
	b.imageSize = imageSize
	return b
}

func (b *ThumbnailImageBuilder) BuildNew() (*thumbnailImageEntity.ThumbnailImage, *exception.BuilderException) {
	err := b.validateCommons()
	if err != nil {
		return nil, err
	}

	return thumbnailImageEntity.NewThumbnailImage(nil, b.imageID, b.name, b.extension, b.contentFile, b.size, b.owner, b.imageSize), nil
}

func (b *ThumbnailImageBuilder) Build() (*thumbnailImageEntity.ThumbnailImage, *exception.BuilderException) {
	err := b.validateAll()
	if err != nil {
		return nil, err
	}

	return thumbnailImageEntity.NewThumbnailImage(b.id, b.imageID, b.name, b.extension, b.contentFile, b.size, b.owner, b.imageSize), nil
}

func (b *ThumbnailImageBuilder) validateAll() *exception.BuilderException {
	err := validators.ValidateNonEmptyStringField("id", b.id)
	if err != nil {
		return exception.NewBuilderException("id", err.Error())
	}

	errBuilder := b.validateCommons()
	if errBuilder != nil {
		return errBuilder
	}

	return nil
}

func (b *ThumbnailImageBuilder) validateCommons() *exception.BuilderException {
	if err := validators.ValidateNonEmptyStringField("imageID", b.imageID); err != nil {
		return exception.NewBuilderException("imageID", err.Error())
	}
	if err := validators.ValidateNonEmptyStringField("name", b.name); err != nil {
		return exception.NewBuilderException("name", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("extension", b.extension); err != nil {
		return exception.NewBuilderException("extension", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("contentFile", b.contentFile); err != nil {
		return exception.NewBuilderException("contentFile", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("owner", b.owner); err != nil {
		return exception.NewBuilderException("owner", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("size", b.size); err != nil {
		return exception.NewBuilderException("size", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("imageSize", b.imageSize); err != nil {
		return exception.NewBuilderException("imageSize", err.Error())
	}

	return nil
}
