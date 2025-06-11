package model

type MgoTransaction struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Digest      string `gorm:"type:varchar(44);not null;column:digest" json:"digest"`
	From        string `gorm:"type:varchar(66);not null;column:from" json:"from"`
	To          string `gorm:"type:varchar(66);not null;column:to" json:"to"`
	Amount      string `gorm:"type:varchar(255);column:amount" json:"amount,omitempty"`
	FromAmount  string `gorm:"type:varchar(255);column:from_amount" json:"from_amount,omitempty"`
	Checkpoint  string `gorm:"type:varchar(255);column:checkpoint" json:"checkpoint,omitempty"`
	CoinType    string `gorm:"type:varchar(255);column:coin_type" json:"coin_type,omitempty"`
	GasOwner    string `gorm:"type:varchar(255);column:gas_owner" json:"gas_owner,omitempty"`
	GasPrice    string `gorm:"type:varchar(255);column:gas_price" json:"gas_price,omitempty"`
	GasBudget   string `gorm:"type:varchar(255);column:gas_budget" json:"gas_budget,omitempty"`
	TimestampMs string `gorm:"column:timestamp_ms" json:"timestamp_ms,omitempty"`
}

type SolTransaction struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Digest      string `gorm:"type:varchar(44);not null;column:digest" json:"digest"`
	From        string `gorm:"type:varchar(66);not null;column:from" json:"from"`
	To          string `gorm:"type:varchar(66);not null;column:to" json:"to"`
	Amount      string `gorm:"type:varchar(255);column:amount" json:"amount,omitempty"`
	FromAmount  string `gorm:"type:varchar(255);column:from_amount" json:"from_amount,omitempty"`
	Checkpoint  string `gorm:"type:varchar(255);column:checkpoint" json:"checkpoint,omitempty"`
	CoinType    string `gorm:"type:varchar(255);column:coin_type" json:"coin_type,omitempty"`
	GasOwner    string `gorm:"type:varchar(255);column:gas_owner" json:"gas_owner,omitempty"`
	GasPrice    string `gorm:"type:varchar(255);column:gas_price" json:"gas_price,omitempty"`
	GasBudget   string `gorm:"type:varchar(255);column:gas_budget" json:"gas_budget,omitempty"`
	TimestampMs string `gorm:"column:timestamp_ms" json:"timestamp_ms,omitempty"`
	Status      int    `gorm:"column:status;type:tinyint(1);default:0"`
}

func (MgoTransaction) TableName() string {
	return "mgo_transactions"
}

func (SolTransaction) TableName() string {
	return "sol_transactions"
}
