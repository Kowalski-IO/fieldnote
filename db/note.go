package db

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"database/sql"
)

type Note struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Owner  int64  `db:"owner_id" json:"owner_id"`
	Public bool   `db:"visible" json:"public"`
	Parts  []Part `json:"parts"`
}

type Part struct {
	ID      int64  `json:"id"`
	Parent  int64  `db:"parent_id" json:"parent"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	Content string `json:"content"`
}

func (n Note) Create() (Note, error) {
	var err error
	tx, _ := conn.Begin()

	res, err := tx.Exec(`INSERT INTO notes (title, owner_id, visible) VALUES (?, ?, ?)`,
		n.Title, n.Owner, n.Public)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when inserting new note")
		return Note{}, errors.Wrap(err, "Error when inserting new note")
	}

	id, _ := res.LastInsertId()
	n.ID = id

	err = insertParts(n, tx)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when inserting new note part")
		return Note{}, errors.Wrap(err, "Error when inserting new note part")
	}

	tx.Commit()

	return n, nil
}

func (n Note) Update() (Note, error) {
	var err error
	tx, _ := conn.Begin()

	_, err = tx.Exec(`UPDATE notes SET title = ?, visible = ? WHERE id = ?`,
		n.Title, n.Public, n.ID)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when updating note")
		return Note{}, errors.Wrap(err, "Error when updating note")
	}

	_, err = tx.Exec(`DELETE FROM parts WHERE parent_id = ?`, n.ID)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when deleting existing note parts")
		return Note{}, errors.Wrap(err, "Error when deleting existing note parts")
	}

	err = insertParts(n, tx)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when updating note part")
		return Note{}, errors.Wrap(err, "Error when updating note part")
	}

	tx.Commit()

	return n, nil
}

func FetchNotes() ([]Note, error) {
	var err error
	var notes []Note

	err = conn.Select(&notes, "SELECT * FROM notes")

	if err != nil {
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Unable to select all notes")
		return []Note{}, errors.Wrap(err, "Unable to select all notes")
	}

	for i := range notes {
		err = conn.Select(&notes[i].Parts, "SELECT * FROM parts WHERE parent_id = ?", notes[i].ID)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Unable to select parts")
		return []Note{}, errors.Wrap(err, "Unable to select parts")
	}

	return notes, nil
}

func insertParts(n Note, tx *sql.Tx) (error) {
	var err error

	for _, v := range n.Parts {
		_, err = tx.Exec(`INSERT INTO parts (parent_id, title, kind, content) VALUES (?, ?, ?, ?)`,
			n.ID, v.Title, v.Kind, v.Content)

		if err != nil {
			break
		}
	}

	return err
}
