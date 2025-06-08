package model

import "time"

type Address struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	Username         string     `gorm:"unique" json:"username"`
	Password         string     `json:"password"`
	MgoAddress       string     `gorm:"unique" json:"mgo_address"`
	MgoPrivateKey    string     `json:"mgo_private_key"`
	SolanaAddress    string     `json:"solana_address"`
	SolanaPrivateKey string     `json:"solana_private_key"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at"`
}

func (Address) TableName() string {
	return "address"
}
