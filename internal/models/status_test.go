package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
)

func TestEnvelopIntoHttpResponse(t *testing.T) {
	type args struct {
		obj     any
		objName string
		code    int
	}
	testChat := Chat{
		ChatID: uuid.New(),
		Name:   "test",
	}
	testChatList := []Chat{
		testChat,
	}
	tests := []struct {
		name string
		args args
		want HttpResponse
	}{
		{
			name: "envelop chat body",
			args: args{
				obj:     testChat,
				objName: "chat",
				code:    http.StatusCreated,
			},
			want: HttpResponse{
				Message: http.StatusText(http.StatusCreated),
				Properties: map[string]any{
					"chat": testChat,
				},
			},
		},
		{
			name: "envelop chat list body",
			args: args{
				obj:     testChatList,
				objName: "chat_list",
				code:    http.StatusOK,
			},
			want: HttpResponse{
				Message: http.StatusText(http.StatusOK),
				Properties: map[string]any{
					"chat_list": testChatList,
				},
			},
		},
		{
			name: "envelop error",
			args: args{
				obj:     "error text",
				objName: "error",
				code:    http.StatusNotFound,
			},
			want: HttpResponse{
				Message: http.StatusText(http.StatusNotFound),
				Properties: map[string]any{
					"error": "error text",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnvelopIntoHttpResponse(tt.args.obj, tt.args.objName, tt.args.code)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvelopIntoHttpResponse() = %v, want %v", got, tt.want)
			}
			body, _ := json.Marshal(got)
			t.Log(string(body))
		})
	}
}
