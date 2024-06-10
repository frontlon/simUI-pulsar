package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const HttpTimeout = 3 * time.Second

// get请求
func HttpGet(uri string, headers map[string]string) ([]byte, error) {
	// 表单数据

	// 创建一个新的请求对象
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 设置请求头部
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// 发送请求
	client := &http.Client{Timeout: HttpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bodyBytes, nil
}

// post form 请求
func HttpPostForm(uri string, formData, headers map[string]string) ([]byte, error) {
	// 表单数据
	data := url.Values{}
	if formData != nil && len(formData) > 0 {
		for k, v := range formData {
			data.Set(k, v)
		}
	}

	// 创建一个缓冲区来存储表单编码后的数据
	buf := bytes.NewBufferString(data.Encode())

	// 创建一个新的请求对象
	req, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 设置请求头部
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// 发送请求
	client := &http.Client{Timeout: HttpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bodyBytes, nil
}

// 下载文件到本地
func DownloadFile(httpUrl string, localPath string) error {

	//下载文件
	response, err := http.Get(httpUrl)
	if err != nil {
		return err
	}

	dir := GetFilePath(localPath)

	if err = CreateDir(dir); err != nil {
		return err
	}

	f, err := os.Create(localPath)
	defer f.Close()

	if err != nil {
		return err
	}

	if _, err := io.Copy(f, response.Body); err != nil {
		return err
	}

	return nil
}
