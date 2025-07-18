package handlers

import (
	"net/http"
	"onfly-api/internal/application/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service}
}

// RegisterInput representa os dados de registro de usuário
type RegisterInput struct {
	Nome      string `json:"nome" example:"João"`
	Email     string `json:"email" example:"joao@email.com"`
	Senha     string `json:"senha" example:"123456"`
	TipoConta string `json:"tipo" example:"empresa"` // "administrador" ou "empresa"
}

// Register godoc
// @Summary Cria um novo usuário
// @Description Registra um novo usuário com nome, e-mail, senha e tipo
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterInput true "Dados para registrar usuário"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	u, err := h.service.Register(input.Nome, input.Email, input.Senha)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":    u.ID,
		"nome":  u.Nome,
		"email": u.Email,
		"tipo":  u.Tipo,
	})
}

// LoginInput representa os dados de login
type LoginInput struct {
	Email string `json:"email" example:"joao@email.com"`
	Senha string `json:"senha" example:"123456"`
}

// Login godoc
// @Summary Realiza login
// @Description Autentica o usuário e retorna um token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginInput true "Credenciais de login"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	token, err := h.service.Login(input.Email, input.Senha)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

// ResetPasswordInput representa os dados para redefinição de senha
type ResetPasswordInput struct {
	Email     string `json:"email" example:"joao@email.com"`
	NovaSenha string `json:"nova_senha" example:"novasenha123"`
}

// ResetPassword godoc
// @Summary Solicita reset de senha
// @Description Envia uma solicitação para redefinir a senha do usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordInput true "Dados para reset de senha"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var input ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	err := h.service.ResetPassword(input.Email, input.NovaSenha)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Senha redefinida com sucesso"})
}
