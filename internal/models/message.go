package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	UserId    string `gorm:"type:varchar(255);not null;index"`
	EventType string `gorm:"type:varchar(25);not null;index"`
	Payload   string `gorm:"type:text;not null;"`
	Read      bool
	Notified  bool
}

func (m *Message) Save() error {
	return db.Create(&m).Error
}

func MessagesSave(messages []Message) error {
	return db.Create(messages).Error
}

func MessagesUpdate(messages []Message) error {
	return db.Save(messages).Error
}

func MessagesMarkAsNotified(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return db.
		Debug().
		Model(&Message{}).
		Where("id IN ?", ids).
		Update("notified", true).
		Error
}

func MessageListNotNotified(userId string, pagination *Pagination) ([]Message, error) {
	var m []Message
	err := db.
		Scopes(pagination.GetScope).
		Where("user_id = ? AND notified = false", userId).
		Find(&m).
		Error

	if err != nil {
		return nil, err
	}
	return m, nil
}
