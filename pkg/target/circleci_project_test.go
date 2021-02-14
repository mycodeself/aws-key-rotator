package target

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mycodeself/aws-key-rotator/pkg/circleci"
	"github.com/mycodeself/aws-key-rotator/pkg/mock"
)

func TestCircleciProjectTargetRotate(t *testing.T) {
	target := &CircleciProjectTarget{
		circleci: &circleci.Client{},
	}
	accessKeyId := "ABCDEF123456789"
	secretAccessKey := "k6Fz8yjggX8tUdMyV9TJykjBrGKYvr9V8mDXE3QE"
	key := &types.AccessKey{
		AccessKeyId:     &accessKeyId,
		SecretAccessKey: &secretAccessKey,
	}
	config := struct {
		project_slug               string
		access_key_id_var_name     string
		secret_access_key_var_name string
	}{}

	target.circleci.HTTPClient = &mock.HTTPClientDoMock{
		DoFn: func(req *http.Request) (*http.Response, error) {
			reqJson := &struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{}

			err := json.NewDecoder(req.Body).Decode(reqJson)
			if err != nil {
				t.Error(err)
			}

			resBody, err := json.Marshal(reqJson)
			if err != nil {
				t.Error(err)
			}

			return &http.Response{
				Status: http.StatusText(201),
				Body:   ioutil.NopCloser(bytes.NewReader(resBody)),
			}, nil
		},
	}
	err := target.Rotate(context.TODO(), key, config)
	if err != nil {
		t.Error(err)
	}

}
