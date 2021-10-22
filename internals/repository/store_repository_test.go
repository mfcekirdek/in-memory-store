//+build unit

package repository

import (
	"reflect"
	"testing"
)

func TestNewStoreRepository(t *testing.T) {
	tests := []struct {
		name string
		want StoreRepository
	}{
		{"New storeRepository instance", &storeRepository{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRepository_Flush(t *testing.T) {
	type fields struct {
		Store map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"Empty store", fields{Store: map[string]string{}}, map[string]string{}},
		{"Not Empty store", fields{Store: map[string]string{"for": "bar"}}, map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeRepository{
				store: tt.fields.Store,
			}
			if got := s.Flush(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRepository_Get(t *testing.T) {
	type fields struct {
		Store map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"Item does not exist", fields{Store: map[string]string{}}, args{key: "foo"}, ""},
		{"Item does exist", fields{Store: map[string]string{"foo": "bar"}}, args{key: "foo"}, "bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeRepository{
				store: tt.fields.Store,
			}
			if got := s.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRepository_GetStore(t *testing.T) {
	type fields struct {
		Store map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"Empty store", fields{Store: map[string]string{}}, map[string]string{}},
		{"Not Empty store", fields{Store: map[string]string{"foo": "bar"}}, map[string]string{"foo": "bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeRepository{
				store: tt.fields.Store,
			}
			if got := s.GetStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRepository_LoadStore(t *testing.T) {
	type fields struct {
		Store map[string]string
	}
	type args struct {
		store map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Empty store", fields{Store: map[string]string{}}, args{store: map[string]string{}}},
		{"Not Empty store", fields{Store: map[string]string{"foo": "bar"}}, args{store: map[string]string{"foo": "bar"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeRepository{
				store: tt.fields.Store,
			}
			s.LoadStore(tt.args.store)
			if got := s.store; !reflect.DeepEqual(got, tt.args.store) {
				t.Errorf("LoadStore() = %v, want %v", got, tt.args.store)
			}
		})
	}
}

func TestStoreRepository_Set(t *testing.T) {
	type fields struct {
		Store map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Nil store", fields{Store: nil}, args{
			key:   "foo",
			value: "newValue",
		}, false},
		{"Empty store", fields{Store: map[string]string{}}, args{
			key:   "foo",
			value: "newValue",
		}, false},
		{"Key does not exist", fields{Store: map[string]string{"key": "value"}}, args{
			key:   "foo",
			value: "newValue",
		}, false},
		{"Key exists", fields{Store: map[string]string{"foo": "bar"}}, args{
			key:   "foo",
			value: "newValue",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeRepository{
				store: tt.fields.Store,
			}
			if got := s.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}
