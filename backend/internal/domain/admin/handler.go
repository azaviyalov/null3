package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, adminJWT echo.MiddlewareFunc) {
	e.POST("/api/admin/auth/login", handler.Login)
	e.POST("/api/admin/auth/logout", handler.Logout, adminJWT)
	e.GET("/api/admin/auth/me", handler.Me, adminJWT)
	e.POST("/api/admin/invites", handler.CreateInvite, adminJWT)
}

type Handler struct {
	accountService *account.Service
	adminService   *Service
	config         Config
	sessionConfig  session.Config
}

func NewHandler(accountService *account.Service, adminService *Service, config Config, sessionConfig session.Config) *Handler {
	return &Handler{
		accountService: accountService,
		adminService:   adminService,
		config:         config,
		sessionConfig:  sessionConfig,
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	token, err := h.adminService.Authenticate(req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return newHTTPError(http.StatusUnauthorized, "Incorrect admin credentials.", err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	session.SetAdminCookie(c, h.sessionConfig, token, adminAccessTokenTTL)
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Logout(c echo.Context) error {
	session.ClearAdminCookie(c, h.sessionConfig)
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Me(c echo.Context) error {
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) CreateInvite(c echo.Context) error {
	rawToken, invite, err := h.accountService.CreateInvite(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	resp := account.InviteResponse{
		InviteURL: h.frontendURL("/invite/" + rawToken),
		ExpiresAt: invite.ExpiresAt,
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) frontendURL(path string) string {
	baseURL := strings.TrimRight(h.config.FrontendURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:4200"
	}
	return baseURL + path
}

func newHTTPError(status int, message string, internal error) error {
	httpError := echo.NewHTTPError(status, message)
	httpError.Internal = internal
	return httpError
}

func emptyJSON(c echo.Context, status int) error {
	return c.JSON(status, struct{}{})
}
