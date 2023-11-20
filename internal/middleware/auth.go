package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/proto/session"
)

type SessionGetter interface {
	GetSession(session models.Session) (models.Session, models.StatusCode)
}

type DefaultGRPCSessionGetter struct {
	Client session.SessionClient
}

func NewDefaultGRPCSessionGetter(Client session.SessionClient) SessionGetter {
	return DefaultGRPCSessionGetter{
		Client: Client,
	}
}

func (g DefaultGRPCSessionGetter) GetSession(s models.Session) (models.Session, models.StatusCode) {
	resp, err := g.Client.GetSession(context.Background(),
		&session.GetSessionRequest{SessionID: s.ID.String()})
	if err != nil {
		slog.Error(err.Error())
		return models.Session{}, models.NotFound
	}

	sID, err := uuid.Parse(resp.SessionID)
	if err != nil {
		return models.Session{}, models.BadRequest
	}

	uID, err := uuid.Parse(resp.UserID)
	if err != nil {
		return models.Session{}, models.BadRequest
	}

	return models.Session{
		ID:     sID,
		UserID: uID,
		Type:   resp.Type,
	}, models.OK
}

type AuthMiddlewareHandler struct {
	SessionGetter       SessionGetter
	RequiredSessionType string
}

// Auth is the middleware function that .
func (h AuthMiddlewareHandler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(*pkg.JwtCustomClaims)
		sessionIDStr := claims.SessionID
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return pkg.UnauthorizedResponse(c, err)
		}

		s, status := h.SessionGetter.GetSession(models.Session{ID: sessionID})
		if status != models.OK {
			return pkg.UnauthorizedResponse(c, err)
		}

		c.Set("session_id", sessionID)
		c.Set("user_id", s.UserID)
		return next(c)
	}
}
