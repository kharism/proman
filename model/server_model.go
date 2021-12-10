package model

import (
	"github.com/kharism/proman/pkg/module/confaccessor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	ID           string                    `bson:"_id" json:"ID"`
	Name         string                    `bson:"Name" json:"Name"`
	Address      string                    `bson:"Address" json:"Address"`
	Path         string                    `bson:"Path" json:"Path"`
	Url          string                    `bson:"Url" json:"Url"`
	AccessorType string                    `bson:"AccessorType" json:"AccessorType"`
	Accessor     confaccessor.ConfAccessor `bson:"-" json:"-"`
}

func (m *Server) GetAddress() string {
	return m.Address
}
func (m *Server) GetPath() string {
	return m.Path
}
func (m *Server) TableName() string {
	return "Server"
}
func (m *Server) ReadData() (string, error) {
	return m.Accessor.ReadData(m)
}
func (m *Server) WriteData(s string) error {
	return m.Accessor.WriteData(m, s)
}
func (m *Server) SetID(values []interface{}) {
	id := values[0]
	if v, ok := id.(string); ok {
		m.ID = v
	} else if v, ok := id.(primitive.ObjectID); ok {
		m.ID = v.Hex()
	}
}

// GetID get model id
func (m *Server) GetID() ([]string, []interface{}) {
	return []string{"ID"}, []interface{}{m.ID}
}
