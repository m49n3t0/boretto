package bot

//import (
//	"bytes"
//	"crypto/tls"
//	"encoding/json"
//	"fmt"
//	"io"
//	"io/ioutil"
//	"mime/multipart"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strings"
//	"time"
//
//	"github.com/mitchellh/mapstructure"
//	"github.com/ovh/venom"
//	"github.com/ovh/venom/executors"
//)
//
//// Headers represents header HTTP for Request
//type Headers map[string]string
//
//// Executor struct. Json and yaml descriptor are used for json output
//type Executor struct {
//	Method            string      `json:"method" yaml:"method"`
//	URL               string      `json:"url" yaml:"url"`
//	Path              string      `json:"path" yaml:"path"`
//	Body              string      `json:"body" yaml:"body"`
//	Headers           Headers     `json:"headers" yaml:"headers"`
//	IgnoreVerifySSL   bool        `json:"ignore_verify_ssl" yaml:"ignore_verify_ssl" mapstructure:"ignore_verify_ssl"`
//	BasicAuthUser     string      `json:"basic_auth_user" yaml:"basic_auth_user" mapstructure:"basic_auth_user"`
//	BasicAuthPassword string      `json:"basic_auth_password" yaml:"basic_auth_password" mapstructure:"basic_auth_password"`
//}
//
//// Run execute TestStep
//func (Executor) Run(testCaseContext venom.TestCaseContext, l venom.Logger, step venom.TestStep) (venom.ExecutorResult, error) {
//
//	// transform step to Executor Instance
//	var t Executor
//
//	r := Result{Executor: t}
//
//	req, err := t.getRequest()
//	if err != nil {
//		return nil, err
//	}
//
//	for k, v := range t.Headers {
//		req.Header.Set(k, v)
//	}
//
//    tr := &http.Transport{
//        TLSClientConfig: &tls.Config{InsecureSkipVerify: t.IgnoreVerifySSL},
//    }
//    client := &http.Client{Transport: tr}
//
//	start := time.Now()
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	elapsed := time.Since(start)
//	r.TimeSeconds = elapsed.Seconds()
//	r.TimeHuman = fmt.Sprintf("%s", elapsed)
//
//	var bb []byte
//	if resp.Body != nil {
//		defer resp.Body.Close()
//		var errr error
//		bb, errr = ioutil.ReadAll(resp.Body)
//		if errr != nil {
//			return nil, errr
//		}
//		r.Body = string(bb)
//
//		bodyJSONArray := []interface{}{}
//		if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
//			bodyJSONMap := map[string]interface{}{}
//			if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
//				r.BodyJSON = bodyJSONMap
//			}
//		} else {
//			r.BodyJSON = bodyJSONArray
//		}
//	}
//
//	r.StatusCode = resp.StatusCode
//
//	return executors.Dump(r)
//}
//
//// getRequest returns the request correctly set for the current executor
//func (e Executor) getRequest() (*http.Request, error) {
//	path := fmt.Sprintf("%s%s", e.URL, e.Path)
//	method := e.Method
//	body := &bytes.Buffer{}
//	if e.Body != "" {
//		body = bytes.NewBuffer([]byte(e.Body))
//	}
//	req, err := http.NewRequest(method, path, body)
//	if err != nil {
//		return nil, err
//	}
//	if len(e.BasicAuthUser) > 0 || len(e.BasicAuthPassword) > 0 {
//	    req.SetBasicAuth(e.BasicAuthUser, e.BasicAuthPassword)
//	}
//	return req, err
//}
//
//
