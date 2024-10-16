package pm3

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPiForm_Handler(t *testing.T) {

	type TestBody struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	type TestResponse struct {
		Id    int    `json:"id"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	api := CreateApiWidthDoc(http.MethodPut, "/api/test/:id", []string{"context", ":id", "body"}, []string{"response"}, func(ctx *gin.Context, id int, body *TestBody) (*TestResponse, error) {
		return &TestResponse{
			Id:    id,
			Name:  body.Name,
			Value: body.Value,
		}, nil
	})
	engine := gin.New()
	engine.Handle(api.Method(), api.Path(), api.Handler)

	rw := httptest.NewRecorder()

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&TestBody{
		Name:  "test",
		Value: time.Now().Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error(err)

		return
	}
	req, err := http.NewRequest(api.Method(), "/api/test/1", buf)
	if err != nil {
		t.Error(err)
		return
	}

	engine.ServeHTTP(rw, req)

	w := rw.Result()

	data, err := io.ReadAll(w.Body)
	if err != nil {
		return
	}
	err = w.Body.Close()
	if err != nil {
		return
	}
	t.Log(string(data))
}
