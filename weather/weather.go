package weather

type Client struct {
	http     *http.Client
	key      string
}

func NewClient(httpClient *http.Client, key string) *Client {
	return &Client{httpClient, key}
}

func (c *Client) FetchWeather(query) (*Results, error) {
	endpoint := fmt.Sprintf("http://api.weatherstack.com/current?access_key =%s&query=%s", c.key, url.QueryEscape(query))
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

	res := &Results{}
	return res, json.Unmarshal(body, res)
}