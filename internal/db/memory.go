package db

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type memoryTable struct {
	name    string
	records map[string][]byte
}

type Memory struct {
	tables map[string]memoryTable
}

func (m *Memory) initDB() error {
	m.tables = map[string]memoryTable{}
	return nil
}

func (m Memory) LookupId(id string, val interface{}) error {
	docType := reflect.TypeOf(val).Elem().Name()
	table, found := m.tables[docType]

	// As we lazily initialize the memory tables they won't exist until the first store
	// therefore we need to return a not found error
	if !found {
		return ErrDBEntityDoesNotExist
	}

	lookup, found := table.records[id]

	if !found {
		return ErrDBEntityDoesNotExist
	}
	err := json.Unmarshal(lookup, &val)
	if err != nil {
		return ErrDBLookupFailed
	}
	return nil
}

func (m *Memory) Store(id string, val interface{}) error {
	docType := reflect.TypeOf(val).Name()
	_, found := m.tables[docType]

	// We lazily initialize tables, so create it if not found
	if !found {
		m.tables[docType] = memoryTable{name: docType, records: map[string][]byte{}}
	}

	jsonStr, err := json.Marshal(val)
	if err != nil {
		return ErrDBStoreFailed
	}
	m.tables[docType].records[id] = jsonStr
	return nil
}

// Debug function to dump the contents of the table
// This code is ugly and not performant but that's okay!
func (m Memory) dump() string {
	result := ""
	for name, table := range m.tables {
		result += name + "\n"
		for id, record := range table.records {
			result += fmt.Sprintf("    %s: %s\n", id, record)
		}
	}
	return result
}
