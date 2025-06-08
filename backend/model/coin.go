package model

type MgoCoin struct {
	ID                 int64   `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	Address            string  `gorm:"column:address;type:varchar(66);not null" json:"address"`
	Module             string  `gorm:"column:module;type:varchar(255);not null" json:"module"`
	Struct             string  `gorm:"column:struct;type:varchar(255);not null" json:"struct"`
	Type               string  `gorm:"column:type;type:varchar(700);not null" json:"type"`
	Decimals           int64   `gorm:"column:decimals;type:int(11);not null" json:"decimals"`
	Name               string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Symbol             string  `gorm:"column:symbol;type:varchar(255);not null" json:"symbol"`
	Description        *string `gorm:"column:description;type:text" json:"description,omitempty"`
	IconURL            *string `gorm:"column:icon_url;type:text" json:"icon_url,omitempty"`
	TotalSupply        string  `gorm:"column:total_supply;type:decimal(60,0);default:0" json:"total_supply"`
	MetadataID         *string `gorm:"column:metadata_id;type:varchar(66)" json:"metadata_id,omitempty"`
	MetadataVersion    *int64  `gorm:"column:metadata_version;type:bigint(20)" json:"metadata_version,omitempty"`
	TreasuryCapID      *string `gorm:"column:treasury_cap_id;type:varchar(66)" json:"treasury_cap_id,omitempty"`
	TreasuryCapVersion *int64  `gorm:"column:treasury_cap_version;type:bigint(20)" json:"treasury_cap_version,omitempty"`
	Creator            *string `gorm:"column:creator;type:varchar(66)" json:"creator,omitempty"`
	CreateTimeMs       *int64  `gorm:"column:create_time_ms;type:bigint(20)" json:"create_time_ms,omitempty"`
	CreateDigest       *string `gorm:"column:create_digest;type:varchar(66)" json:"create_digest,omitempty"`
}

func (MgoCoin) TableName() string {
	return "mgo_coins"
}
