package brickset

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AuthSuccess struct {
}

type BricksetTheme struct {
	Status  string `json:"status"`
	Matches int    `json:"matches"`
	Themes  []struct {
		Theme         string `json:"theme"`
		Setcount      int    `json:"setCount"`
		Subthemecount int    `json:"subthemeCount"`
		Yearfrom      int    `json:"yearFrom"`
		Yearto        int    `json:"yearTo"`
	} `json:"themes"`
}

type BricksetSubtheme struct {
	Status    string `json:"status"`
	Matches   int    `json:"matches"`
	Subthemes []struct {
		Theme    string `json:"theme"`
		Subtheme string `json:"subtheme"`
		Setcount int    `json:"setCount"`
		Yearfrom int    `json:"yearFrom"`
		Yearto   int    `json:"yearTo"`
	} `json:"subthemes"`
}

func sendRequest(apiKey string, url string, body string) (*resty.Response, error) {

	client := resty.New()

	//client.SetDebug(true)

	baseUrl := "https://brickset.com/api/v3.asmx"
	fullUrl := fmt.Sprint(baseUrl, url)

	resp, _ := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContentLength(true).
		SetBody(body).
		SetResult(AuthSuccess{}).
		Post(fullUrl)

	return resp, nil
}

func GetThemes(apiKey string) (BricksetTheme, error) {

	var bricksetResponse BricksetTheme
	body := fmt.Sprint("apiKey=", apiKey)

	resp, err := sendRequest(apiKey, "/getThemes", body)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	return bricksetResponse, nil
}

func GetSubthemes(apiKey string, theme string) (BricksetSubtheme, error) {

	var bricksetResponse BricksetSubtheme
	body := fmt.Sprint("apiKey=", apiKey, "&theme=", theme)

	resp, err := sendRequest(apiKey, "/getSubthemes", body)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	return bricksetResponse, nil
}
