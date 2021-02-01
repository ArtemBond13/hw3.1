package transaction

import (
	"reflect"
	"sync"
	"testing"
)

func TestService_MapRowToTransaction(t *testing.T) {
	type fields struct {
		mu           sync.Mutex
		transactions []*Transaction
	}
	type args struct {
		rows []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Ok", fields{mu: sync.Mutex{}, transactions: []*Transaction{}},
			args{rows: []string{"1", "0001", "0002", "1000_99", "1611473617"}},
			Transaction{"1", "0001", "0002", 1000_99, 1611473617}, false},
		{"Empty slice", fields{mu: sync.Mutex{}, transactions: []*Transaction{}},
			args{rows: []string{"", "", "", "", "", "", ""}}, Transaction{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				mu:           tt.fields.mu,
				transactions: tt.fields.transactions,
			}
			got, err := s.MapRowToTransaction(tt.args.rows)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapRowToTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapRowToTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
