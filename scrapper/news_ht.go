package scrapper

import (
	"capture-feed/utility"
	"github.com/gocolly/colly"
)

type ScrapperHindustandTimes struct {
	SCInfo
}

func ScrapperHindustandTimesNew() *ScrapperHindustandTimes {
	return &ScrapperHindustandTimes{
		SCInfo{
			Name: "Hindustan Times",
			Logo: "https://www.hindustantimes.com/images/app-images/ht/htlogo.png",
			Domains: []string{"hindustantimes.com", "www.hindustantimes.com"},
		},
	}
}

func (ht *ScrapperHindustandTimes) GetArticle(request *SCNewsArticleRequest) *SCNewsArticle {
	article := ht.Scrap(request.URL)
	article.SourceLogoURL = ht.Logo
	article.SourceName = ht.Name
	article.SourceURL = request.URL
	article.Hash = utility.GetHash([]byte(request.URL))
	article.PubDate = request.PubDate

	return article
}


func (ht *ScrapperHindustandTimes) Scrap(url string) *SCNewsArticle {
	// Instantiate default collector
	c := colly.NewCollector()

	article := SCNewsArticle{}

	// On every a element which has href attribute call callback
	c.OnHTML("div.storyArea", func(e *colly.HTMLElement) {

		article.Headline = e.ChildText("h1")
		article.SubHeadline = e.ChildText("h2")

		e.ForEach("div.storyDetail p", func(index int, e *colly.HTMLElement){
			para := SCNewsParagraph{}
			para.Content = e.Text
			article.Paragraphs = append(article.Paragraphs, &para)
		})

		e.ForEach("figure", func(index int, e *colly.HTMLElement){
			img := SCNewsImage{}
			img.URL = e.ChildAttr("img", "src")
			img.Caption = e.ChildText("figcaption")
			article.Posters = append(article.Posters, &img)
		})

	})

	c.Visit(url)

	return &article
}
