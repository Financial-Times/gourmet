package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDecodeRequestFromQueryParameters(t *testing.T) {

	type normalReq struct {
		Param int `query:"param"`
	}

	type defaultRequest struct {
		Param int `query:"param" default:"456"`
	}

	type faultyReq struct {
		param int `query:"param"`
	}

	happyRequest := &http.Request{}
	happyURL, _ := url.Parse("https://example.com/test?param=123")
	happyRequest.URL = happyURL

	emptyRequest := &http.Request{}

	reusableRequest := &normalReq{}

	tests := []struct {
		name      string
		reqStruct interface{}
		request   *http.Request
		err       error
		assertRes func(interface{}) error
	}{
		{
			name:      "test that values are read from request query parameters",
			reqStruct: &normalReq{},
			request:   happyRequest,
			assertRes: func(req interface{}) error {
				res, ok := req.(*normalReq)
				if !ok {
					return fmt.Errorf("decoded request cannot be casted to original type")
				}
				if res.Param != 123 {
					return fmt.Errorf("data in decoded request is not correct")
				}
				return nil
			},
		},
		{
			name:      "",
			reqStruct: &defaultRequest{},
			request:   emptyRequest,
			assertRes: func(req interface{}) error {
				res, ok := req.(*defaultRequest)
				if !ok {
					return fmt.Errorf("decoded request cannot be casted to original type")
				}
				if res.Param != 456 {
					return fmt.Errorf("data in decoded request is not correct")
				}
				return nil
			},
		},
		{
			name:      "test that original struct is not modified",
			reqStruct: reusableRequest,
			request:   happyRequest,
			assertRes: func(req interface{}) error {
				res, ok := req.(*normalReq)
				if !ok {
					return fmt.Errorf("decoded request cannot be casted to original type")
				}
				if res.Param != 123 {
					return fmt.Errorf("data in decoded request is not correct")
				}
				if reusableRequest.Param != 0 {
					return fmt.Errorf("original struct was modified")
				}
				return nil
			},
		},
		{
			name:      "test loader error is handled",
			reqStruct: &faultyReq{},
			request:   happyRequest,
			err:       fmt.Errorf("error decoding request: field 'param' cannot be set"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoder := DecodeRequestFromQueryParameters(test.reqStruct)
			res, err := decoder(context.TODO(), test.request)
			if err != nil {
				if test.err == nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if test.err.Error() != err.Error() {
					t.Fatalf("exptected err: %s; got: %s", test.err, err)
				}
				return
			}

			err = test.assertRes(res)
			if err != nil {
				t.Fatalf("incorrectly decoded request: %s", err)
			}
		})
	}
}

func TestDecodeRequestFromJSONBody(t *testing.T) {

	type normalReq struct {
		Param int `query:"param"`
	}

	type defaultRequest struct {
		Param int `query:"param" default:"456"`
	}

	req1 := &http.Request{}
	req1.Body = ioutil.NopCloser(strings.NewReader("{\"param\": 789}"))

	reusableRequest := &normalReq{}
	req2 := &http.Request{}
	req2.Body = ioutil.NopCloser(strings.NewReader("{\"param\": 789}"))

	tests := []struct {
		name      string
		reqStruct interface{}
		request   *http.Request
		err       error
		assertRes func(interface{}) error
	}{
		{
			name:      "test that props are set",
			reqStruct: &normalReq{},
			request:   req1,
			assertRes: func(req interface{}) error {
				res, ok := req.(*normalReq)
				if !ok {
					return fmt.Errorf("decoded request cannot be casted to original type")
				}
				if res.Param != 789 {
					return fmt.Errorf("data in decoded request1 is not correct")
				}
				return nil
			},
		},
		{
			name:      "test that original struct is not modified",
			reqStruct: reusableRequest,
			request:   req2,
			assertRes: func(req interface{}) error {
				if reusableRequest.Param != 0 {
					return fmt.Errorf("original struct was modified")
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoder := DecodeRequestFromJSONBody(test.reqStruct)
			res, err := decoder(context.TODO(), test.request)
			if err != nil {
				if test.err == nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if test.err.Error() != err.Error() {
					t.Fatalf("exptected err: %s; got: %s", test.err, err)
				}
				return
			}

			err = test.assertRes(res)
			if err != nil {
				t.Fatalf("incorrectly decoded request: %s", err)
			}
		})
	}
}
