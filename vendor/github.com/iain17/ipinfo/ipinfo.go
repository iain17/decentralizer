package ipinfo

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type ipInfo struct {
	Ip string
	CountryCode string
}

func GetIpInfo() (*ipInfo, error) {
	var info *ipInfo
	var err error
	info, err = getIpWithOptionA()
	if err == nil {
		return info, nil
	}
	info, err = getIpWithOptionB()
	return info, err
}

type optionA struct {
	CountryCode string  `json:"countryCode"`
	Query       string  `json:"query"`
}

//Option A
func getIpWithOptionA() (*ipInfo, error) {
	res, err := http.Get("http://www.ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response optionA
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &ipInfo{
		Ip: response.Query,
		CountryCode: response.CountryCode,
	}, nil
}

type optionB struct {
	Ip string  `json:"ip"`
	Country  string  `json:"country"`
}

func getIpWithOptionB() (*ipInfo, error) {
	res, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response optionB
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &ipInfo{
		Ip: response.Ip,
		CountryCode: response.Country,
	}, nil
}