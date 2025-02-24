package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type SearXNG struct {
	url string
}

func NewSearXNG(url string) *SearXNG {
	return &SearXNG{
		url: url,
	}
}

func (s *SearXNG) Search(query string) (Result, error) {
	if s.url[len(s.url)-1] == '/' {
		s.url = s.url[:len(s.url)-1]
	}
	url := s.url + "/search?q=" + query + "&format=json"

	fmt.Println(url)

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Cookie", `categories=general; language=auto; locale=zh-Hans-CN; autocomplete=; favicon_resolver=; image_proxy=0; method=POST; safesearch=0; theme=simple; results_on_new_tab=0; doi_resolver=oadoi.org; simple_style=auto; center_alignment=0; advanced_search=0; query_in_title=0; infinite_scroll=0; search_on_category_select=1; hotkeys=default; url_formatting=pretty; disabled_engines=; enabled_engines=bing__general; disabled_plugins=; enabled_plugins=; tokens=`)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return Result{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		var result Result
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return Result{}, err
		}
		result.Results = result.Results[:5]
		return result, nil
	}
	return Result{}, nil
}

type Result struct {
	Query           string `json:"query"`
	NumberOfResults int    `json:"number_of_results"`
	Results         []struct {
		URL           string   `json:"url"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		PublishedDate string   `json:"publishedDate"`
		Thumbnail     string   `json:"thumbnail"`
		Engine        string   `json:"engine"`
		Template      string   `json:"template"`
		ParsedURL     []string `json:"parsed_url"`
		Engines       []string `json:"engines"`
		Positions     []int    `json:"positions"`
		Score         float64  `json:"score"`
		Category      string   `json:"category"`
	} `json:"results"`
}

func (r *Result) ToString() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (r *Result) ToResultsContent() *ResultsContent {

	contents := make([]struct {
		Url     string `json:"url"`
		Content string `json:"content"`
	}, len(r.Results))

	for i, result := range r.Results {
		contents[i].Url = result.URL
	}
	return &ResultsContent{Contents: contents}
}

type ResultsContent struct {
	Contents []struct {
		Url     string `json:"url"`
		Content string `json:"content"`
	}
}

func (r *ResultsContent) GetContents() *ResultsContent {
	wg := &sync.WaitGroup{}
	for i, _ := range r.Contents {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			// fmt.Println(index)
			res, err := http.Get(r.Contents[index].Url)
			if err != nil {
				r.Contents[index].Content = err.Error()
				return
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				r.Contents[index].Content = fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				r.Contents[index].Content = err.Error()
				return
			}

			doc.Find(`
				script, style, noscript, 
				link[rel='stylesheet'], 
				meta[name^='generator'],
				img[loading='lazy'],
				iframe, embed, object,
				form, button, input,
				svg, canvas 
			`).Remove()

			doc.Find("*").FilterFunction(func(i int, s *goquery.Selection) bool {
				if style, exists := s.Attr("style"); exists {
					return strings.Contains(strings.ToLower(style), "display:none") ||
						strings.Contains(strings.ToLower(style), "visibility:hidden")
				}
				return false
			}).Remove()

			r.Contents[index].Content = doc.Find("body").Text()
			r.Contents[index].Content = strings.ReplaceAll(r.Contents[index].Content, "\n", " ")
			r.Contents[index].Content = strings.ReplaceAll(r.Contents[index].Content, "\t", " ")
			r.Contents[index].Content = strings.ReplaceAll(r.Contents[index].Content, "  ", " ")

			// fmt.Println(index)
			wg.Done()
		}(i, wg)
	}
	wg.Wait()
	// fmt.Println(r)
	return r
}

func (r *ResultsContent) ToString() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
