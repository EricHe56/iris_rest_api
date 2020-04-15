package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

func HttpRequest(method string, url string, body string) (status int, responseBody string) {
	responseBody = ""
	request, _ := http.NewRequest(method, url, strings.NewReader(body))
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("HttpRequest Error: ", err)
		status = -1
		responseBody = err.Error()
		return
	}
	defer response.Body.Close()

	status = response.StatusCode
	responseBytes, _ := ioutil.ReadAll(response.Body)
	responseBody = string(responseBytes)
	//StdPrint.Info(method, url, body, responseBody)
	return
}

func HttpRequestWithHeader(method string, url string, body string, headers map[string]string) (status int, responseBody string, responseHeaders http.Header) {
	responseBody = ""
	request, _ := http.NewRequest(method, url, strings.NewReader(body))
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("HttpRequestWithHeader Error: ", err)
		status = -1
		responseBody = err.Error()
		return
	}
	defer response.Body.Close()

	status = response.StatusCode
	responseHeaders = response.Header
	responseBytes, _ := ioutil.ReadAll(response.Body)
	responseBody = string(responseBytes)
	//StdPrint.Info(method, url, body, responseBody)
	return
}

func ReadFileInString(fileName string) (txt string, err error) {
	var txtBytes []byte
	txtBytes, err = ioutil.ReadFile(fileName)
	if err == nil {
		txt = string(txtBytes)
	}
	return
}

func WriteFile(fileName, txt string) (err error) {
	data := []byte(txt)
	err = ioutil.WriteFile(fileName, data, 0644)
	return
}

func GetFilesAndDirs(dirPath string) (files []string, dirs []string, err error) {
	dir, e := ioutil.ReadDir(dirPath)
	if e != nil {
		return nil, nil, e
	}

	PathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPath+PathSep+fi.Name())
			f, d, e1 := GetFilesAndDirs(dirPath + PathSep + fi.Name())
			if e1 != nil {
				return nil, nil, e1
			}
			files = append(files, f...)
			dirs = append(dirs, d...)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPath+PathSep+fi.Name())
			}
		}
	}
	return
}

func SortMapByKeys(mapInput map[string]string) (queryStr string) {
	keys := make([]string, 0, len(mapInput))
	for k, _ := range mapInput {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var i = 0
	for k1 := range keys {
		var v = mapInput[keys[k1]]
		var k = keys[k1]
		//mapInput[keys[k1]] = mapInput[keys[k1]]
		if v != "" {
			if i > 0 {
				queryStr += "&"
			}
			queryStr += k + "=" + v
			i++
		}
	}
	return
}
