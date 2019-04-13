package dao

import (
	"fmt"

	"proofn/task5/models"

	"github.com/go-pg/pg"
)

type Message struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (d *Message) Connect() error {
	var n int

	//conn string
	db = pg.Connect(&pg.Options{
		User:     d.User,
		Password: d.Password,
		Addr:     fmt.Sprintf("%s:%s", d.Host, d.Port),
		Database: d.Database,
	})

	//Check our connection
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}

func (d *Message) Close() error {
	err := db.Close()
	return err
}

func (d *Message) FindAll() ([]models.Message, error) {
	var messages []models.Message

	//Go get the messages
	err := db.Model(&messages).Select()
	if err != nil {
		return []models.Message{}, err
	}

	return messages, nil
}

func (d *Message) DeleteAll() error {
	var ids []int

	//Find the messages ids
	err := db.Model(&Message{}).Column("id").Select(&ids)
	if err != nil {
		return err
	}

	//Delete the message ids if we have results
	if len(ids) > 0 {
		pgids := pg.In(ids)
		_, err := db.Model(&Message{}).Where("id IN (?)", pgids).Delete()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Message) Insert(message models.Message) (models.Message, error) {
	err := db.Insert(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}
