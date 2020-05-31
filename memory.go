package gosql

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
)

type Cell interface {
	AsText() string
	AsInt() int32
}

type Results struct {
	Columns []struct {
		Type ColumnType
		Name string
	}
	Rows [][]Cell
}

var (
	ErrTableDoesNotExist  = errors.New("Table does not exist")
	ErrColumnDoesNotExist = errors.New("Column does not exist")
	ErrInvalidSelectItem  = errors.New("Select item is not valid")
	ErrInvalidDatatype    = errors.New("Invalid datatype")
	ErrMissingValues      = errors.New("Missing values")
)

type BackEnd interface {
	CreateTable(*CreateTableStatement) error
	Insert(*InsertStatement) error
	Select(*SelectStatement) (*Results, error)
}

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
	var i int32
	if err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i); err != nil {
		panic(err)
	}

	return i
}

func (mc MemoryCell) AsText() string {
	return string(mc)
}

type table struct {
	columns     []string
	columnTypes []ColumnType
	rows        [][]MemoryCell
}

type MemoryBackend struct {
	tables map[string]*table
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{tables: map[string]*table{}}
}
