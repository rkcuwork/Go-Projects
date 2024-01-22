package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	photoapi = "https://api.pexels.com/v1/"
	videoapi = "https://api.pexels.com/videos/"
)

type Client struct {
	token          string
	hc             http.Client
	remainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{
		token: token,
		hc:    c,
	}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string   `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	ID              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"Photographer"`
	PhorographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Potrait   string `json:"potrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	Photos   []Photo `json:"photos"`
	NextPage string   `json:"next_page"`
}

func (c *Client) SearchPhoto(query string, perPage int32, page int32) (*SearchResult, error) {
	url := fmt.Sprintf(photoapi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		log.Fatalf("main :: SearchPhoto :: Error while c.requestDoWithAuth() %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("main :: SearchPhoto :: Error while ioutil.ReadAll() %v", err)
	}
	// fmt.Printf("data = %s\n", data)

	var result SearchResult
	err = json.Unmarshal(data, &result)
	return &result, err

}

func (c *Client) CuratedPhotos(perPage, Page int32) (*CuratedResult, error) {
	url := fmt.Sprintf(photoapi+"/curated?per_page=%d&page=%d", perPage, Page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		log.Fatalf("main :: CuratedPhotos :: Error while c.requestDoWithAuth() :: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("main :: CuratedPhotos :: Error while ioutil.ReadAll() :: %v", err)
	}
	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}


func (c *Client) GetPhotos(id int32) (*Photo,error){
	url := fmt.Sprintf(photoapi+"/photos/%d", id)
    resp, err := c.requestDoWithAuth("GET", url)
    if err!= nil {
		log.Fatalf("main :: GetPhotos :: Error while c.requestDoWithAuth() :: %v", err)
    }
    defer resp.Body.Close()
	
    data, err := ioutil.ReadAll(resp.Body)
    if err!= nil {
		log.Fatalf("main :: GetPhotos :: Error while ioutil.ReadAll() :: %v", err)
    }
    var result Photo
    err = json.Unmarshal(data, &result)
    return &result,err
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	res,err := c.CuratedPhotos(1,int32(randNum))
	if err!= nil {
		log.Fatalf("main :: GetRandomPhoto :: Error while c.CuratedPhotos() :: %v", err)
		} else if (len(res.Photos) ==1){
		return &res.Photos[0],nil
	}
	return nil,err
	
}

func (c *Client) requestDoWithAuth(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("main :: requestDoWithAuth :: Error while http.NewRequest() %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s", c.token))
	res, err := c.hc.Do(req)
	if err != nil {
		log.Fatalf("main :: requestDoWithAuth :: Error while c.hc.Do() %v", err)
	}

	// fmt.Printf("status code : %d\n", res.StatusCode)
	// fmt.Println("Headers:")
	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("%s: %s\n", key, value)
	// 	}
	// }

	timesStr := res.Header.Get("X-Ratelimit-Remaining")
	fmt.Printf("timestr : %s \n", timesStr)
	times, err := strconv.Atoi(timesStr)
	if err != nil {
		log.Fatalf("main :: requestDoWithAuth :: Error while strconv.Atoi() %v", err)
	}
	c.remainingTimes = int32(times)
	return res, err
}

func (c *Client) GetRemainingRequestForThisMonth() (int32){
	return c.remainingTimes
}

func main() {
	fmt.Println("Welcome to RK-world!")
	
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("main :: main :: Error while godotenv.Load() :: %v", err)
	}
	
	apiKey := os.Getenv("PEXELS_API_KEY")
	
	var c = NewClient(apiKey)
	
	// result, err := c.SearchPhoto("cars", 1, 5)
	// result,err := c.CuratedPhotos(2,3)
	result,err := c.GetPhotos(19193788)
	if err != nil {
		fmt.Println(fmt.Errorf("main :: main :: Error while c.SearchPhoto() :: %v \n", err))
	}
	// if result.Page == 0 {
		// 	fmt.Errorf("main :: main :: page = 0, while c.SearchPhoto()")
	// }

	fmt.Println(result)
	
}
