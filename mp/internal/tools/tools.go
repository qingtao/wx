// Package tools 实验性分离常用工具型代码，减小代码的重复率,暂时放置于internal目录下
package tools

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// ReadResponse 读取HTTP响应的内容
func ReadResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return b, nil
}

// UnmarshalJSON 读取HTTP响应的内容并使用json反序列化数据到v
func UnmarshalJSON(resp *http.Response, v interface{}) error {
	b, err := ReadResponse(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// UnmarshalXML 读取HTTP响应的内容并使用xml反序列化数据到v
func UnmarshalXML(resp *http.Response, v interface{}) error {
	b, err := ReadResponse(resp)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, v)
}
