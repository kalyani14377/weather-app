package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	http     *http.Client
	key      string
}


type Results struct {
	Location struct {
		Name           string `json:"name"`
		Country        string `json:"country"`
		Region         string `json:"region"`
		Lat            string `json:"lat"`
		Lon            string `json:"lon"`
		TimezoneID     string `json:"timezone_id"`
		Localtime      string `json:"localtime"`
		LocaltimeEpoch int    `json:"localtime_epoch"`
		UtcOffset      string `json:"utc_offset"`
	} 
}

type ForecastResults struct {
	Data []struct {
		Weather      struct {
			Icon        string `json:"icon"`
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"weather"`
		Datetime  string  `json:"datetime"`
		Temp      float64     `json:"temp"`
		Station   string  `json:"station"`
		AppTemp   float64     `json:"app_temp"`
	} `json:"data"`
}

func NewClient(httpClient *http.Client, key string) *Client {
	return &Client{httpClient, key}
}

func (c *Client) FetchWeather(query string) (*ForecastResults, error) {
	endpoint := fmt.Sprintf("http://api.weatherstack.com/forecast?access_key=%s&query=%s", c.key, url.QueryEscape(query))
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	resBody:=string(body)
	var result map[string]interface{}
	json.Unmarshal([]byte(resBody), &result)
	location := result["location"].(map[string]interface{})
	lat:=location["lat"]
	lng:=location["lon"]
	fmt.Println(lat)
	fmt.Println(lng)
	forecastApi := fmt.Sprintf("https://weatherbit-v1-mashape.p.rapidapi.com/current?lon=%s&lat=%s", lng ,lat)
	req, _ := http.NewRequest("GET", forecastApi, nil)

	req.Header.Add("x-rapidapi-key", "623bab92a7msh2474d3437ce4a9dp1bf32ajsn18777857c239")
	req.Header.Add("x-rapidapi-host", "weatherbit-v1-mashape.p.rapidapi.com")

	res1, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		return nil, err1
	}
	
	defer res1.Body.Close()
	body1, _ := ioutil.ReadAll(res1.Body)
	if err != nil {
		return nil, err
	}

	if res1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body1))
	}
	res := &ForecastResults{}
	return res, json.Unmarshal(body1, res)
}