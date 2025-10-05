package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"url-shortener/internal/dto"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

const aliasLength = 6

type UrlShortener struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}
type Service interface {
	CreateAlias(ctx *fiber.Ctx) error
	GetURL(ctx *fiber.Ctx) error
	DeleteUrl(ctx *fiber.Ctx) error
}
type service struct {
	repo storage.Repository
	log  *zap.SugaredLogger
}

func New(repo storage.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

func (s *service) CreateAlias(ctx *fiber.Ctx) error {
	s.log.Info("CreateAlias")
	var req UrlShortener
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	alias := req.Alias
	if alias == "" {
		alias = random.NewRandomString(aliasLength)
	}
	aliasShortener := UrlShortener{
		Url:   req.Url,
		Alias: alias,
	}
	aliasShortenerID, err := s.repo.SaveURL(ctx.Context(), aliasShortener.Url, aliasShortener.Alias)
	if err != nil {
		s.log.Error("Error creating alias", zap.Error(err))
	}

	response := dto.Response{
		Status: "success",
		Data:   map[string]int64{"alias": aliasShortenerID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetURL(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	if alias == "" {
		s.log.Error("Empty alias")
		return ctx.Status(fiber.StatusBadRequest).SendString("Alias is required")
	}

	targetURL, err := s.repo.GetURL(ctx.Context(), alias)
	if err != nil {
		s.log.Error("Alias not found",
			zap.String("alias", alias),
			zap.Error(err))
		return ctx.Status(fiber.StatusNotFound).SendString("Short link not found")
	}

	s.log.Info("Short link redirected successfully",
		zap.String("from", ctx.Request().URI().String()),
		zap.String("to", targetURL))

	return ctx.Redirect(targetURL, fiber.StatusFound)
}

func (s *service) DeleteUrl(ctx *fiber.Ctx) error {
	alias := ctx.Params("alias")
	if alias == "" {
		s.log.Error("Empty alias")
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Alias is required")
	}
	err := s.repo.DeleteUrl(ctx.Context(), alias)
	if err != nil {
		s.log.Error("Failed to delete URL", zap.String("alias", alias), zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Failed to delete alias")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
