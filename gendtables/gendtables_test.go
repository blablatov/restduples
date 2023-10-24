package main

import (
	"fmt"
	"reflect"
	"testing"
)

type Generater interface {
	WriteString(string) (interface{}, error)
}

type TestBuffer struct {
	Generater
	success bool
}

func (gn *TestBuffer) WriteString(string) (interface{}, error) {
	if gn.success {
		return nil, nil
	}
	return nil, fmt.Errorf("WriteString test error")
}

func TestWriteString(t *testing.T) {
	type args struct {
		ws Generater
	}
	tests := []struct {
		strw    string
		args    args
		wantInt int
		wantErr error
	}{
		{
			strw: "strw exists",
			args: args{
				ws: &TestBuffer{success: true},
			},
			wantInt: 1,
			wantErr: nil,
		}, {
			strw: "strw not exists",
			args: args{
				ws: &TestBuffer{success: false},
			},
			wantInt: 0,
			wantErr: fmt.Errorf("strw test error"),
		},
	}
	var f TestBuffer

	for _, tt := range tests {
		t.Run(tt.strw, func(t *testing.T) {
			gotInt, gotErr := f.WriteString("test")
			if gotInt == tt.wantInt {
				t.Errorf("Check func WriteString() gotInt = %v, want %v", gotInt, tt.wantInt)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func WriteString() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
