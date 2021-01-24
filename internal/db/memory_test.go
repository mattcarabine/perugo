package db

import (
	"errors"
	"reflect"
	"testing"
)

func TestMemoryDB(t *testing.T) {
	// Internal structs to use in the tests
	type foo struct {
		Id string `json:"id"`
	}

	type bar struct {
		Name string `json:"name"`
	}

	var tests = []struct {
		descr     string
		storageFn func(Memory) (interface{}, error)
		wantVal   interface{}
		wantErr   error
	}{
		{"Set then get same type", func(memory Memory) (interface{}, error) {
			fooObj := foo{"matt"}
			err := memory.Store("key1", fooObj)
			if err != nil {
				return nil, err
			}
			var result foo
			err = memory.LookupId("key1", &result)
			return result, err
		}, foo{"matt"}, nil},
		{"Set then get different type", func(memory Memory) (interface{}, error) {
			fooObj := foo{"matt"}
			err := memory.Store("key1", fooObj)
			if err != nil {
				return nil, err
			}
			var result bar
			err = memory.LookupId("key1", &result)
			return result, err
		}, bar{}, ErrDBEntityDoesNotExist},
		{"Get without set", func(memory Memory) (interface{}, error) {
			var result foo
			err := memory.LookupId("key1", &result)
			return result, err
		}, foo{}, ErrDBEntityDoesNotExist},
		{"Multiple sets",func(memory Memory) (interface{}, error) {
			fooObj := foo{"matt"}
			err := memory.Store("key1", fooObj)
			if err != nil {
				return nil, err
			}

			fooObj = foo{"PJ"}
			err = memory.Store("key1", fooObj)
			if err != nil {
				return nil, err
			}

			var result foo
			err = memory.LookupId("key1", &result)
			return result, err
		}, foo{"PJ"}, nil},
		{"Set then get different key",func(memory Memory) (interface{}, error) {
			fooObj := foo{"matt"}
			err := memory.Store("key1", fooObj)
			if err != nil {
				return nil, err
			}

			var result foo
			err = memory.LookupId("key3", &result)
			return result, err
		}, foo{}, ErrDBEntityDoesNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			mem := Memory{}
			err := mem.initDB()
			if err != nil {
				t.Fatalf("failed to initialize mem: %v", err)
			}

			val, err := tt.storageFn(mem)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v but got %v", tt.wantErr, err)
				t.Logf("Dump: %s", mem.dump())
				return
			}

			if val != tt.wantVal {
				valType := reflect.TypeOf(val).Name()
				wantValType := reflect.TypeOf(tt.wantVal).Name()
				t.Errorf("expected value %s%v but got %s%v", wantValType, tt.wantVal, valType, val)
				t.Logf("Dump:\n %s", mem.dump())
			}
		})
	}

}
