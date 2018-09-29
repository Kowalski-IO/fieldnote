package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Note struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Owner  uuid.UUID `db:"owner_id" json:"owner_id"`
	Title  string    `db:"title" json:"title"`
	Public bool      `db:"visible" json:"public"`
	Tags   []string  `db:"tags" json:"tags"`
	Parts  []Part    `db:"-" json:"parts"`
}

type Part struct {
	ID       uuid.UUID `db:"id" json:"id"`
	NoteID   uuid.UUID `db:"note_id" json:"note_id"`
	Position int64     `db:"position" json:"position"`
	Title    string    `db:"title" json:"title"`
	Kind     string    `db:"kind" json:"kind"`
	Content  string    `db:"content" json:"content"`
}

func (n Note) Upsert() (Note, error) {
	tx, _ := conn.Begin()

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	_, err := tx.Exec("UPSERT INTO notes (id, owner_id, title, visible, tags) VALUES ($1, $2, $3, $4, $5)",
		n.ID, n.Owner, n.Title, n.Public, pq.Array(n.Tags))

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when upserting note")
		return n, errors.Wrap(err, "Error when upserting note")
	}

	var ext []string

	for _, v := range n.Parts {
		ext = append(ext, v.ID.URN())
	}

	fmt.Println(ext)

	_, err = tx.Exec("DELETE FROM parts WHERE note_id = $1 AND id NOT IN ($2)", n.ID, ext)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when deleting parts that have been removed")
		return n, errors.Wrap(err, "Error when deleting parts that have been removed")
	}

	for i, v := range n.Parts {
		if v.ID == uuid.Nil {
			v.ID = uuid.New()
		}

		v.NoteID = n.ID
		_, err = tx.Exec("UPSERT INTO parts (id, note_id, position, title, kind, content) VALUES ($1, $2, $3, $4, $5, $6)",
			v.ID, n.ID, n.Parts[i].Position, n.Parts[i].Title, n.Parts[i].Kind, n.Parts[i].Content)

		if err != nil {
			break
		}

		n.Parts[i] = v
	}

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Error when upserting note")
		return n, errors.Wrap(err, "Error when upserting parts")
	}

	tx.Commit()

	return n, nil
}

func FetchNotes() ([]Note, error) {
	var err error
	var n []Note

	rows, err := conn.Queryx("SELECT id, owner_id, title, visible, tags FROM notes")

	if err != nil {
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Error("Unable to select all notes")
		return []Note{}, errors.Wrap(err, "Unable to select all notes")
	}

	defer rows.Close()

	for rows.Next() {
		t := Note{}

		err = rows.Scan(&t.ID, &t.Owner, &t.Title, &t.Public, pq.Array(&t.Tags))

		if err != nil {
			log.WithFields(log.Fields{
				"cause": err.Error(),
			}).Error("Unable to populate note")
			return []Note{}, errors.Wrap(err, "Unable to populate note")
		}

		err := conn.Select(&t.Parts, "SELECT id, note_id, position, title, kind, content FROM parts WHERE note_id = $1", t.ID)

		if err != nil {
			if err != nil {
				log.WithFields(log.Fields{
					"cause": err.Error(),
				}).Error("Unable to populate parts")
				return []Note{}, errors.Wrap(err, "Unable to populate parts")
			}
		}

		n = append(n, t)
	}

	return n, nil
}
