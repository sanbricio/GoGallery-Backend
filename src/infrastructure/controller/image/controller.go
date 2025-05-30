package imageController

import (
	"fmt"
	"go-gallery/src/commons/exception"
	validators "go-gallery/src/commons/utils/validations"
	imageService "go-gallery/src/service/image"
	userService "go-gallery/src/service/user"
	"strconv"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	log "go-gallery/src/infrastructure/logger"

	imageHandler "go-gallery/src/infrastructure/controller/image/handler"
	userDTO "go-gallery/src/infrastructure/dto/user"

	"github.com/gofiber/fiber/v2"
)

const (
	INVALID_AUTHENTIFICATION_MSG string = "User not authenticated"
	IMAGE_ID_REQUIRED_MSG        string = "Image ID is required"
	DEFAULT_PAGE_SIZE            int64  = 10
)

var logger log.Logger

type ImageController struct {
	imageService *imageService.ImageService
	userService  *userService.UserService
}

func NewImageController(imageService *imageService.ImageService, userService *userService.UserService) *ImageController {
	logger = log.Instance()
	return &ImageController{
		imageService: imageService,
		userService:  userService,
	}
}

func (c *ImageController) SetUpRoutes(router fiber.Router) {
	//Image
	router.Get("/getImage/:id", c.getImage)
	router.Post("/uploadImage", c.uploadImage)
	router.Put("/updateImage", c.updateImage)
	router.Delete("/deleteImage/", c.deleteImage)

	// Thumbnail
	router.Get("/getThumbnailImages", c.getThumbnailImages)
}

//	@Summary		Obtiene una imagen por su identificador
//	@Description	Obtiene una imagen específica del usuario según el identificador proporcionado
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Identificador de la imagen"
//	@Security		CookieAuth
//	@Success		200	{object}	imageDTO.ImageDTO
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/getImage/{id} [get]
func (c *ImageController) getImage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		logger.Error(IMAGE_ID_REQUIRED_MSG)
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, IMAGE_ID_REQUIRED_MSG))
	}
	logger.Info("GET /getImage called with id: " + id)

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(INVALID_AUTHENTIFICATION_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	dtoFindImage := &imageDTO.ImageDTO{
		Id:    &id,
		Owner: claims.Username,
	}

	image, err := c.imageService.Find(dtoFindImage)
	if err != nil {
		logger.Error(fmt.Sprintf("Error finding image with id %s : %s", id, err.Message))
		return ctx.Status(err.Status).JSON(err)
	}

	logger.Info("Image successfully retrieved with id: " + id)
	return ctx.Status(fiber.StatusOK).JSON(image)
}

//	@Summary		Persiste una imagen
//	@Description	Permite a un usuario autenticado persistir una imagen
//	@Tags			image
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Archivo de imagen a subir (jpeg, jpg, png, webp)"
//	@Security		CookieAuth
//	@Success		200	{object}	imageDTO.ImageDTO		"Imagen subida correctamente"
//	@Failure		400	{object}	exception.ApiException	"Error al procesar la imagen"
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		409	{object}	exception.ApiException	"La imagen ya existe"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/uploadImage [post]
func (c *ImageController) uploadImage(ctx *fiber.Ctx) error {
	logger.Info("POST /uploadImage called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(INVALID_AUTHENTIFICATION_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	fileInput, errForm := ctx.FormFile("file")
	if errForm != nil {
		logger.Error("Failed to get file from form data caused by:" + errForm.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(exception.NewApiException(fiber.StatusNotFound, "Error getting image from form"))
	}

	logger.Info("Processing image upload for user: " + claims.Username)
	dtoInsertImage, errFile := imageHandler.ProcessImageFile(fileInput, claims.Username)
	if errFile != nil {
		logger.Error("Error processing image file: " + errFile.Message)
		return ctx.Status(errFile.Status).JSON(errFile)
	}

	dto, errInsert := c.imageService.Insert(dtoInsertImage)
	if errInsert != nil {
		logger.Error("Error inserting image: " + errInsert.Message)
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	logger.Info("Image successfully uploaded by user: " + claims.Username)
	return ctx.Status(fiber.StatusOK).JSON(dto)
}

//	@Summary		Elimina una imagen
//	@Description	Borra una imagen específica del usuario autentificado
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			request	body	imageDTO.ImageDeleteRequestDTO	true	"Datos para eliminar la imagen"
//	@Security		CookieAuth
//	@Success		200	{object}	dto.MessageResponseDTO	"Imagen eliminada correctamente"
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/deleteImage [delete]
func (c *ImageController) deleteImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(INVALID_AUTHENTIFICATION_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	request := new(imageDTO.ImageDeleteRequestDTO)
	err := ctx.BodyParser(request)
	if err != nil {
		errorMessage := "Invalid JSON in update request"
		logger.Error(errorMessage)
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, errorMessage))
	}

	// Validate required fields
	if err := validators.ValidateNonEmptyStringField("id", request.Id); err != nil {
		logger.Error("Image ID is required for update")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, err.Error()))
	}

	if err := validators.ValidateNonEmptyStringField("thumbnail_id", request.ThumbnailID); err != nil {
		logger.Error("Thumbnail ID is required for update")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, err.Error()))
	}

	request.Owner = claims.Username

	logger.Info("DELETE /deleteImage called with id: " + request.Id)

	response, errDelete := c.imageService.Delete(request)
	if errDelete != nil {
		logger.Error("Error deleting image with id " + request.Id + ": " + errDelete.Message)
		return ctx.Status(errDelete.Status).JSON(err)
	}

	logger.Info(fmt.Sprintf("Image and thumbnail successfully deleted with image id: %s and thumbnail id: %s", request.Id, request.ThumbnailID))
	return ctx.Status(fiber.StatusOK).JSON(response)
}

//	@Summary		Actualiza el nombre de una imagen
//	@Description	Actualiza una imagen específica del usuario autentificado
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			request	body	imageDTO.ImageUpdateRequestDTO	true	"Datos de actualización para la imagen"
//	@Security		CookieAuth
//	@Success		200	{object}	imageDTO.ImageUpdateResponseDTO	"Imagen actualizada correctamente"
//	@Failure		400	{object}	exception.ApiException			"JSON invalido"
//	@Failure		401	{object}	exception.ApiException			"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException			"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException			"Ha ocurrido un error inesperado"
//	@Router			/image/updateImage [put]
func (c *ImageController) updateImage(ctx *fiber.Ctx) error {
	logger.Info("UPDATE /updateImage called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(INVALID_AUTHENTIFICATION_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	request := new(imageDTO.ImageUpdateRequestDTO)
	err := ctx.BodyParser(request)
	if err != nil {
		errorMessage := "Invalid JSON in update request"
		logger.Error(errorMessage)
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, errorMessage))
	}
	// Validate required fields
	if err := validators.ValidateNonEmptyStringField("id", request.Id); err != nil {
		logger.Error("Image ID is required for update")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, err.Error()))
	}

	if err := validators.ValidateNonEmptyStringField("name", request.Name); err != nil {
		logger.Error("Image name is required for update")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, err.Error()))
	}

	if err := validators.ValidateNonEmptyStringField("thumbnail_id", request.ThumbnailID); err != nil {
		logger.Error("Thumbnail ID is required for update")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, err.Error()))
	}

	request.Owner = claims.Username

	result, errUpdate := c.imageService.Update(request)
	if errUpdate != nil {
		logger.Error("Error updating image with id " + request.Id + ": " + errUpdate.Message)
		return ctx.Status(errUpdate.Status).JSON(errUpdate)
	}

	logger.Info(fmt.Sprintf("Image and thumbnail successfully updated with image id: %s and thumbnail id: %s", request.Id, request.ThumbnailID))
	return ctx.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Listar imágenes en miniatura (thumbnails)
//	@Description	Obtiene una lista paginada de imágenes en miniatura del usuario autenticado, usando paginación por cursor (lastId y pageSize).
//	@Tags			thumbnail
//	@Accept			json
//	@Produce		json
//	@Param			lastID		query	string	false	"Último ID recibido para la paginación"
//	@Param			pageSize	query	int		false	"Cantidad de miniaturas a devolver (por defecto 10)"
//	@Security		CookieAuth
//	@Success		200	{object}	thumbnailImageDTO.ThumbnailImageCursorDTO	"Lista de miniaturas con el último id para poder realizar paginacione"
//	@Failure		401	{object}	exception.ApiException						"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException						"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException						"No se encontraron thumbnails"
//	@Failure		500	{object}	exception.ApiException						"Error inesperado"
//	@Router			/image/getThumbnailImages [get]
func (c *ImageController) getThumbnailImages(ctx *fiber.Ctx) error {
	lastID := ctx.Query("lastID")
	pageSizeParam := ctx.Query("pageSize")

	logger.Info("GET /getThumbnailImages called with lastID: " + lastID + ", pageSize: " + pageSizeParam)

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(INVALID_AUTHENTIFICATION_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	// PagesSize validation (it must be positive by default 10)
	pageSize := int64(DEFAULT_PAGE_SIZE)
	if pageSizeParam != "" {
		if parsedPageSize, err := strconv.ParseInt(pageSizeParam, 10, 64); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	thumbnails, errThumb := c.imageService.FindAllThumbnails(claims.Username, lastID, pageSize)
	if errThumb != nil {
		logger.Error("Error retrieving thumbnails: " + errThumb.Message)
		return ctx.Status(errThumb.Status).JSON(errThumb)
	}

	logger.Info("Thumbnails successfully retrieved for user: " + claims.Username)
	return ctx.Status(fiber.StatusOK).JSON(thumbnails)
}
