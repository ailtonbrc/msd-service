package models

import "gorm.io/gorm"

// SystemLog representa um log do sistema
type SystemLog struct {
	gorm.Model

	UserID     *uint                  `json:"user_id"`
	User       *Usuario               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string                 `gorm:"size:100;not null" json:"action"`
	EntityType string                 `gorm:"size:50" json:"entity_type"`
	EntityID   string                 `json:"entity_id"`
	Details    map[string]interface{} `gorm:"type:jsonb" json:"details"`
	IPAddress  string                 `gorm:"size:45" json:"ip_address"`
}

// TableName especifica o nome da tabela
func (SystemLog) TableName() string {
	return "system_logs"
}
