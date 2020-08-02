package transaction

import (
	"reflect"
	"testing"
)

func TestMapRowToTransaction(t *testing.T) {
	type args struct {
		record []string
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		// TODO: Add test cases.
		{
			name: "Тестируем корректность маппинга для транзакций",
			args: args{
				record: []string{"xxxxx", "0001", "0002", "10000000", "15957988534"},
			},
			want: &Transaction{
				Id:      "xxxxx",
				From:    "0001",
				To:      "0002",
				Amount:  10_000_000,
				Created: 15957988534,
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapRowToTransaction(tt.args.record); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapRowToTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
