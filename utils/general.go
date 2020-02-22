package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func HttpRequest(method string, url string, body string) (status int, responseBody string) {
	responseBody = ""
	request, _ := http.NewRequest(method, url, strings.NewReader(body))
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, _ := client.Do(request)
	defer response.Body.Close()

	status = response.StatusCode
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
