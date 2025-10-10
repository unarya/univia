package queries

import (
	"github.com/google/uuid"
	sessions "github.com/unarya/univia/internal/api/modules/session/model"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/pkg/types"
	"github.com/unarya/univia/pkg/utils"
)

func InsertIntoSessionByValidSession(sessionID, userID uuid.UUID, refreshTokenID uuid.UUID, meta types.SessionMetadata) error {
	db := mysql.DB
	session := sessions.UserSession{
		SessionID:      sessionID,
		UserID:         userID,
		IP:             meta.IP,
		UserAgent:      meta.UserAgent,
		RefreshTokenID: refreshTokenID,
		LastActive:     utils.NowPtr(),
	}
	if err := db.Create(&session).Error; err != nil {
		return err
	}
	return nil
}

func InsertNewSessionByUserID(userID uuid.UUID, refreshTokenID uuid.UUID, meta types.SessionMetadata) (uuid.UUID, error) {
	db := mysql.DB
	session := sessions.UserSession{
		SessionID:      uuid.New(),
		UserID:         userID,
		IP:             meta.IP,
		UserAgent:      meta.UserAgent,
		RefreshTokenID: refreshTokenID,
		LastActive:     utils.NowPtr(),
	}
	if err := db.Create(&session).Error; err != nil {
		return uuid.Nil, err
	}
	return session.SessionID, nil
}
