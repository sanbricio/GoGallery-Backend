package dependency_dictionary

import (
	emailsender_repository "api-upload-photos/src/infrastructure/repository/emailSender"
	image_repository "api-upload-photos/src/infrastructure/repository/image"
	user_repository "api-upload-photos/src/infrastructure/repository/user"
)

func FindImageDependency(code string, args map[string]string) image_repository.ImageRepository {
	switch code {
	case image_repository.ImageMongoDBRepositoryKey:
		return image_repository.NewImageMongoDBRepository(args)
	default:
		return image_repository.NewImageMemoryRepository()
	}
}

func FindUserDependency(code string, args map[string]string) user_repository.UserRepository {
	switch code {
	case user_repository.UserMongoDBRepositoryKey:
		return user_repository.NewUserMongoDBRepository(args)
	default:
		return user_repository.NewUserPostgreSQLRepository(args)
	}

}

func FindEmailSenderDependency(code string, args map[string]string) emailsender_repository.EmailSenderRepository {
	switch code {
	case emailsender_repository.EmailSenderGoMailRepositoryKey:
		return emailsender_repository.NewEmailSenderGoMailRepository(args)
	default:
		panic("EmailSenderRepository not found")
	}
}
