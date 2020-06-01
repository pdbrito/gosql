package gosql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
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

func (mb *MemoryBackend) CreateTable(crt *CreateTableStatement) error {
	t := table{}
	mb.tables[crt.name.value] = &t
	if crt.cols == nil {
		return nil
	}

	for _, col := range *crt.cols {
		t.columns = append(t.columns, col.name.value)

		var dt ColumnType
		switch col.datatype.value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		default:
			return ErrInvalidDatatype
		}

		t.columnTypes = append(t.columnTypes, dt)
	}

	return nil
}

func (mb *MemoryBackend) Insert(inst *InsertStatement) error {
	table, ok := mb.tables[inst.table.value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if inst.values == nil {
		return nil
	}

	row := []MemoryCell{}

	if len(*inst.values) != len(table.columns) {
		return ErrMissingValues
	}

	for _, val := range *inst.values {
		if val.kind != literalKind {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, mb.tokenToCell(val.literal))
	}

	table.rows = append(table.rows, row)
	return nil
}

func (mb *MemoryBackend) tokenToCell(t *token) MemoryCell {
	if t.kind == numericKind {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.value)
		if err != nil {
			panic(err)
		}

		if err := binary.Write(buf, binary.BigEndian, int32(i)); err != nil {
			panic(err)
		}
		return MemoryCell(buf.Bytes())
	}
	if t.kind == stringKind {
		return MemoryCell(t.value)
	}
	return nil
}
