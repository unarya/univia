package notifications

import (
	Users "github.com/deva-labs/univia-api/api/gin-gonic/src/modules/user/models"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SenderID   uuid.UUID  `gorm:"type:uuid;not null"`
	ReceiverID uuid.UUID  `gorm:"type:uuid;not null"`
	Sender     Users.User `gorm:"foreignKey:SenderID;references:ID"`
	Receiver   Users.User `gorm:"foreignKey:ReceiverID;references:ID"`
	Message    string     `gorm:"type:text;default:null"`
	IsSeen     bool       `gorm:"default:false"`
	NotiType   string     `gorm:"type:varchar(50);default:null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}
