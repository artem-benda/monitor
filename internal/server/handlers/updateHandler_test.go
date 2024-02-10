package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string
		requestPath string
		want        want
	}{
		{
			name:        "Receive valid counter metric",
			requestPath: "/update/counter/testcounter/3",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name:        "Receive valid gauge metric",
			requestPath: "/update/gauge/testcounter/3",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name:        "Invalid metric kind",
			requestPath: "/update/invalid/testcounter/3",
			want: want{
				code:        400,
				response:    "Metric type not supported\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Empty counter metric name",
			requestPath: "/update/counter/",
			want: want{
				code:        404,
				response:    "Metric name cannot be empty\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Empty gauge metric name",
			requestPath: "/update/gauge/",
			want: want{
				code:        404,
				response:    "Metric name cannot be empty\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Non numeric counter metric value",
			requestPath: "/update/counter/testmetric",
			want: want{
				code:        400,
				response:    "Bad metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Non numeric gauge metric value",
			requestPath: "/update/gauge/testmetric",
			want: want{
				code:        400,
				response:    "Bad metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Not a number counter metric value",
			requestPath: "/update/counter/testmetric/badval",
			want: want{
				code:        400,
				response:    "Bad metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Not a number gauge metric value",
			requestPath: "/update/gauge/testmetric/badval",
			want: want{
				code:        400,
				response:    "Bad metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Для каждого теста нужно новое хранилище, чтобы результаты не зависили от порядка выполнения
			store := storage.NewStorage()
			handler := MakeUpdateHandler(store)

			request := httptest.NewRequest(http.MethodPost, test.requestPath, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
