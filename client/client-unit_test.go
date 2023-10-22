// Модульное (unit) тестирование клиента.
// Определены тестовые сигнатуры реальных методов клиента.
// Объявлены тестовые типы соответствующие интерфейсу, как его экземпляры.
// go test -v client-unit_test.go

package main

import (
	"fmt"
	"reflect"
	"testing"
)

type Clienter interface {
	LoadX509KeyPair(...interface{}) (interface{}, error)
	Get(string) (interface{}, error)
	Set(...string)
	NewRequest(interface{}, string, interface{}) (interface{}, error)
	Do(interface{}) (interface{}, error)
	ReadAll(interface{}) (interface{}, error)
}

type TestClient struct {
	Clienter
	success bool
}

func (tс *TestClient) LoadX509KeyPair(...interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("LoadX509KeyPair test error")
}

func TestLoadX509KeyPair(t *testing.T) {
	type args struct {
		key Clienter
	}
	tests := []struct {
		crtFile    string
		keyFile    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			crtFile: "crtFile exists",
			args: args{
				key: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			crtFile: "crtFile not exists",
			args: args{
				key: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("crtFile test error"),
			wantExists: false,
		},

		{
			keyFile: "keyFile exists",
			args: args{
				key: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			keyFile: "keyFile not exists",
			args: args{
				key: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("keyFile test error"),
			wantExists: false,
		},
	}

	var f TestClient

	for _, tt := range tests {
		t.Run(tt.crtFile, func(t *testing.T) {
			gotExists, gotErr := f.LoadX509KeyPair(tt.args.key, "crtFile", "keyFile")
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func LoadX509KeyPair() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func LoadX509KeyPair() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestClient) NewRequest(interface{}, string, interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("NewRequest test error")
}

func TestNewRequest(t *testing.T) {
	type args struct {
		req Clienter
	}
	tests := []struct {
		apiUrl     string
		method     string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			apiUrl: "apiUrl exists",
			args: args{
				req: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			apiUrl: "apiUrl not exists",
			args: args{
				req: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("apiUrl test error"),
			wantExists: false,
		}, {
			method: "method exists",
			args: args{
				req: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			method: "method not exists",
			args: args{
				req: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("method test error"),
			wantExists: false,
		},
	}

	var f TestClient

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			gotExists, gotErr := f.NewRequest("GET", "apiUrl", nil)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func NewRequest() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func NewRequest() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestClient) Do(interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("Do test error")
}

func TestDo(t *testing.T) {
	type args struct {
		do Clienter
	}
	tests := []struct {
		request    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			request: "request exists",
			args: args{
				do: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			request: "request not exists",
			args: args{
				do: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("request do test error"),
			wantExists: false,
		},
	}

	var f TestClient

	for _, tt := range tests {
		t.Run(tt.request, func(t *testing.T) {
			gotExists, gotErr := f.Do(tt.args.do)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func Do() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func Do() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestClient) ReadAll(interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("ReadAll test error")
}

func TestReadAll(t *testing.T) {
	type args struct {
		rall Clienter
	}
	tests := []struct {
		read       string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			read: "read exists",
			args: args{
				rall: &TestClient{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			read: "read not exists",
			args: args{
				rall: &TestClient{success: false},
			},
			wantErr:    fmt.Errorf("read do test error"),
			wantExists: false,
		},
	}

	var f TestClient

	for _, tt := range tests {
		t.Run(tt.read, func(t *testing.T) {
			gotExists, gotErr := f.ReadAll(tt.args.rall)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func ReadAll() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func ReadAll() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
