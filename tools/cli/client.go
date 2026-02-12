package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type Client struct {
	clientId     string
	clientSecret string
	httpClient   *http.Client
	accessToken  *oauth2.Token
}

func NewClient(clientId string, clientSecret string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   httpClient,
	}
}

func (client *Client) GetAccessToken(clientId string, clientSecret string) (*oauth2.Token, error) {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientId)
	formData.Set("client_secret", clientSecret)
	httpRequest, err := http.NewRequest("POST", "https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResponse.Body)

	var token oauth2.Token

	jsonDecoder := json.NewDecoder(httpResponse.Body)
	err = jsonDecoder.Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (client *Client) getAccessToken() (*oauth2.Token, error) {
	if !client.accessToken.Valid() {
		token, err := client.GetAccessToken(client.clientId, client.clientSecret)
		if err != nil {
			return nil, err
		}
		client.accessToken = token
	}
	return client.accessToken, nil
}

func (client *Client) GetAllStates(request AllStatesRequest) (*AllStatesResponse, error) {
	u, err := url.Parse("https://opensky-network.org/api/states/all")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if !request.Time.IsZero() {
		q.Add("time", strconv.FormatInt(request.Time.Unix(), 10))
	}
	if request.ICAO24 != "" {
		q.Add("icao24", request.ICAO24)
	}
	if request.Extended {
		q.Add("extended", "1")
	}
	if request.MinLat != 0 {
		q.Add("lamin", strconv.FormatFloat(request.MinLat, 'f', -1, 64))
	}
	if request.MinLon != 0 {
		q.Add("lomin", strconv.FormatFloat(request.MinLon, 'f', -1, 64))
	}
	if request.MaxLat != 0 {
		q.Add("lamax", strconv.FormatFloat(request.MaxLat, 'f', -1, 64))
	}
	if request.MaxLon != 0 {
		q.Add("lomax", strconv.FormatFloat(request.MaxLon, 'f', -1, 64))
	}
	u.RawQuery = q.Encode()

	accessToken, err := client.getAccessToken()
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken.AccessToken))

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResponse.Body)

	var response AllStatesResponse

	jsonDecoder := json.NewDecoder(httpResponse.Body)
	err = jsonDecoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type AllStatesRequest struct {
	Time     time.Time
	ICAO24   string
	Extended bool
	MinLat   float64
	MinLon   float64
	MaxLat   float64
	MaxLon   float64
}

type AllStatesResponse struct {
	Time   int64   `json:"time"`
	States [][]any `json:"states"`
}

func (response *AllStatesResponse) WriteToCSV(w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	err := csvWriter.Write([]string{
		"icao24",
		"callsign",
		"origin_country",
		"time_position",
		"last_contact",
		"longitude",
		"latitude",
		"baro_altitude",
		"on_ground",
		"velocity",
		"true_track",
		"vertical_rate",
		"sensors",
		"geo_altitude",
		"squawk",
		"spi",
		"position_source",
		"category",
	})
	if err != nil {
		return err
	}

	for _, row := range response.States {
		var cells []string
		for _, col := range row {
			cells = append(cells, fmt.Sprintf("%v", col))
		}
		err = csvWriter.Write(cells)
		if err != nil {
			return err
		}
	}

	return nil
}
