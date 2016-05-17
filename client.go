package googlesearch

import (
  "github.com/google/go-querystring/query"
  "github.com/PuerkitoBio/goquery"
  "net/http"
  "net/url"
  "log"
)

const (
  GoogleSearchURL = "https://www.google.com/search"
  Max_Results = 10
)


type Client struct {
  Header http.Header
  HTTPClient *http.Client
  transport *http.Transport
  proxyURL *url.URL
}

type Query struct {
  Query string `url:"q"`
}


func NewClient() *Client {
  client := &Client{
    Header: http.Header{}, 
    HTTPClient: &http.Client{},
    transport:  &http.Transport{},
  }
  return client
}


func (c *Client) SetProxy(proxy string) *Client {
  if pURL, err := url.Parse(proxy); err == nil {
    c.proxyURL = pURL
  } else {
    log.Printf("ERROR [%v]", err)
  }

  return c
}

func (c *Client) SetHeader(header, value string) *Client {
  c.Header.Set(header, value)
  return c
} 

func (c *Client) GetSearchResults(query string) (SearchResultList, error) {
  
  results := make([]SearchResult, Max_Results)
  
  if req, err := http.NewRequest("GET", createSearchURL(query), nil); err == nil {
    req.Header = c.Header
    
    if c.proxyURL != nil {
      c.transport.Proxy = http.ProxyURL(c.proxyURL)
    }

    c.HTTPClient.Transport = c.transport

    resp, err := c.HTTPClient.Do(req)

    if err != nil {
      return SearchResultList{}, err
    }

    if document, err := goquery.NewDocumentFromResponse(resp); err != nil {
      return SearchResultList{}, err
    } else {
      
      document.Find("div.g").Each(func(i int, s *goquery.Selection) {
        result_link, exists := s.Find(".r").Find("a").Attr("href")
        result_description := s.Find(".s").Find(".st").Text()
        
        if exists {
          results[i] = SearchResult{result_description, result_link}
        }
      })

    }
  } else {
    return SearchResultList{}, err
  }

  return SearchResultList{results}, nil
}

func createSearchURL(search string) string {
  v, _ := query.Values(&Query{search})
  return GoogleSearchURL + "?" + v.Encode()
}
