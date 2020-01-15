package utils

import (
	"io/ioutil"
	"net/http"
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
