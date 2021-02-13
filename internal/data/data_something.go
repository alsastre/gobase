package data

import (
	"errors"

	"github.com/upper/db/v4"
	"go.uber.org/zap"
)

// Data ..
type Data struct {
	DBSession db.Session
	Logger    *zap.Logger
}

// Something structure that represents something in the database
type Something struct {
	ID    string `json:"id" db:"id,omitempty"`
	Name  string `json:"name" db:"name"`
	Value int    `json:"value" db:"value"`
}

// Collection were the data will be stored
const collection string = "somethings"

// Dummy data to populate the DB
var dummyListOfSomethings = []*Something{
	{ID: "", Name: "Africanus", Value: 11},
	{ID: "", Name: "I, Julia", Value: 21},
	{ID: "", Name: "Red Rising", Value: 31},
	{ID: "", Name: "Golden Son", Value: 41},
	{ID: "", Name: "Harry Potter", Value: 51},
}

// SQL to create the table of the dummy data
var dummyTableSQL = `CREATE TABLE "` + collection + `" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "value" INTEGER
);`

// NewData returns a new Data object
func NewData(dbSession db.Session, logger *zap.Logger) *Data {
	return &Data{DBSession: dbSession, Logger: logger}
}

// ListSomethings function that retrieves the full List of Somethings from the database
func (d *Data) ListSomethings() ([]*Something, error) {
	// Consult DB
	var some []*Something
	err := d.DBSession.Collection(collection).Find().All(&some)
	// Error is handled outside data package
	return some, err
}

// GetSomething obtains a Something with the given ID
func (d *Data) GetSomething(ID string) (*Something, error) {
	var some *Something
	// Find Something by ID
	err := d.DBSession.Collection(collection).Find(db.Cond{"id": ID}).One(&some)
	// Error is handled outside data package
	return some, err
}

// UpdateSomething updates a Something
func (d *Data) UpdateSomething() (*Something, error) {
	// Update Something
	return nil, errors.New("Not implemented")
}

// DeleteSomething deletes a Something from the DB
func (d *Data) DeleteSomething() (*Something, error) {
	// Delete it and return it
	return nil, errors.New("Not implemented")
}

// FillSession fills the database with dummy values and creates the collection if it does not exists
func (d *Data) FillSession() {
	// Create SQL table to intialize postgresql with dummy values.
	// TODO REMOVE This should not be here since it is not generic for every DB
	if exits, _ := d.DBSession.Collection(collection).Exists(); !exits {
		_, err := d.DBSession.SQL().Exec(dummyTableSQL)
		if err != nil {
			d.Logger.Error("Could not execute", zap.String("Query", dummyTableSQL))
		}
	} else {
		d.Logger.Info("Table already exists")
	}

	// Fill table with data (valid for any DB)
	collection := d.DBSession.Collection(collection)
	for _, dummy := range dummyListOfSomethings {
		_, err := collection.Insert(dummy)
		if err != nil {
			d.Logger.Error("Could not insert", zap.Any("dummy", dummy))
		}
	}
}
