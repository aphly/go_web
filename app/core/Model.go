package core

type Model struct {
	CreatedAt int64 `gorm:"autoUpdateTime"`
	UpdatedAt int64 `gorm:"autoCreateTime"`
}

type ModelId struct {
	ID uint `gorm:"primarykey"`
	Model
}

type ModelDelete struct {
	ModelId
	DeletedAt int64 `gorm:"index"`
}
