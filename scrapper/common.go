package scrapper

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type SCInfo struct {
	Name 	string
	Logo	string
	Domains 	[]string
}

type SCContext struct {
	scrappers []interface{}
}

//Represents a single image
type SCNewsImage struct {
	URL     string
	Caption string
}

// Represents a single paragraph
type SCNewsParagraph struct {
	Content string
}

//Represents the whole article
type SCNewsArticle struct {
	Hash			string

	SourceName		string	
	SourceLogoURL	string //Max 128x128
	SourceURL		string

	Headline    string
	SubHeadline string
	PubDate		time.Time
	Posters      []*SCNewsImage
	Paragraphs  []*SCNewsParagraph
}

type SCNewsArticleRequest struct {
	Hash 		string
	URL			string
	Headline	string
	PubDate		time.Time
}

type SCNewsApi interface {
	GetArticle(request *SCNewsArticleRequest) *SCNewsArticle
	MatchDomain(domain string) bool
}

//Print onto the console
func (sc *SCNewsArticle) Print() {
	fmt.Printf("*************\n%s\n*************\n%s\n*************\n", sc.Headline, sc.SubHeadline)
	for _, p := range sc.Paragraphs {
		fmt.Printf("\n%s\n", p.Content)
	}

	for _, p := range sc.Posters {
		fmt.Printf("\n%s\n", p.URL)
		fmt.Printf("%s\n", p.Caption)
	}
}



func SCContextNew() *SCContext {
	c := SCContext {}

	//Hindustan Times
	c.scrappers = append(c.scrappers, ScrapperHindustandTimesNew())
	
	return &c

}

func (c *SCContext)GetArticleAsync(request *SCNewsArticleRequest, channel chan *SCNewsArticle) error {
	urlInfo, err := url.Parse(request.URL)
	if err != nil {
		return err
	}

	for _, scrapper := range c.scrappers {
		if scrap, ok := scrapper.(SCNewsApi); ok {
			if scrap.MatchDomain(urlInfo.Host) {
				go func() {
					article := scrap.GetArticle(request)
					channel <- article
				}()

				return nil
			}
		}
	}


	return errors.New("No scrapper for domain " + urlInfo.Host)
}


func (info *SCInfo) MatchDomain(domain string) bool {
	for _, d := range info.Domains {
		if d == domain {
			return true
		}
	}
	return false
}


