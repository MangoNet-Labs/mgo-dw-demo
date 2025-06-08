package model

type MgoCheckpoint struct {
	SequenceNumber           string `gorm:"primaryKey;column:sequence_number;type:bigint(20);not null"`
	CheckpointDigest         string `gorm:"column:checkpoint_digest;type:varchar(44);not null"`
	Epoch                    string `gorm:"column:epoch;type:bigint(20);not null"`
	Transactions             string `gorm:"column:transactions;type:mediumtext;not null"`
	PreviousCheckpointDigest string `gorm:"column:previous_checkpoint_digest;type:varchar(44);default:null"`
	EndOfEpoch               bool   `gorm:"column:end_of_epoch;type:tinyint(1);not null"`
	TotalTransactionBlocks   int    `gorm:"column:total_transaction_blocks;type:bigint(20);not null"`
	NetworkTotalTransactions string `gorm:"column:network_total_transactions;type:bigint(20);not null"`
	TimestampMs              string `gorm:"column:timestamp_ms;type:bigint(20);not null"`
	Status                   int    `gorm:"column:status;type:tinyint(1);default:0"`
}

func (MgoCheckpoint) TableName() string {
	return "mgo_checkpoints"
}
