package main

import (
	"fmt"
	"os"
	"os/signal"

	"capture-feed/scrapper"
	"capture-feed/capture"
	"capture-feed/feed"
)

func errorAndDie(code int, err string) {
	fmt.Println(err)
	os.Exit(code)
}

func waitOnFeed(channel chan feed.FeedItem, db *capture.DatabaseContext) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//Prepare Scrapper Context
	scContext := scrapper.SCContextNew()
	scArticleChannel := make(chan *scrapper.SCNewsArticle, 100)

	for {
		select {
		case <-c:
			fmt.Println("Killing gracefully")
			os.Exit(0)

			//All the scrapped article will be recieved here
		case articleItem := <-scArticleChannel: 
			{
				if err := capture.CaptureUpdateArticle(db, articleItem); err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
				}
			}

		case feedItem := <-channel:
			{
				feedItem.FeedItemPrintStatus()
				for _, feed := range feedItem.Items {
					capture := capture.CaptureItemNewFromFeedItem(feed)
					exists, _ := capture.Exits(db)
					if !exists {
						capture.Save(db)
						capture.Log()
						//Push to Article Channel
						if err := scContext.GetArticleAsync(&scrapper.SCNewsArticleRequest{
							Hash: capture.Hash,
							Headline: capture.Title,
							URL: capture.URL,
							PubDate: capture.PubDate,
						}, scArticleChannel); err != nil {
							fmt.Fprint(os.Stderr, err.Error)
						}

						continue
					}

					//Capture already in Database
				}
			}


		}
	}
}



func main() {

	feedsData := `[
		{
			"name":"HT Top Stories",
			"url":"https://www.hindustantimes.com/rss/topnews/rssfeed.xml",
			"interval":10,
			"kind":0
		},
		{
			"name":"HT World",
			"url":"https://www.hindustantimes.com/rss/world/rssfeed.xml",
			"interval":10,
			"kind":0
		},
		{
			"name":"HT Hollywood",
			"url":"https://www.hindustantimes.com/rss/hollywood/rssfeed.xml",
			"interval":10,
			"kind":0
		},
		{
			"name":"HT Bollywood",
			"url":"https://www.hindustantimes.com/rss/bollywood/rssfeed.xml",
			"interval":10,
			"kind":0
		},
		{
			"name":"HT Education",
			"url":"https://www.hindustantimes.com/rss/education/rssfeed.xml",
			"interval":10,
			"kind":0
		},
		{
			"name":"HT India",
			"url":"https://www.hindustantimes.com/rss/india/rssfeed.xml",
			"interval":10,
			"kind":0
		}
		
	]
	`

	db := capture.DatabaseContext{
		Uri: "mongodb://localhost:27017",
	}

	if err := db.Connect(); err != nil {
		errorAndDie(1, err.Error())
	}

	channel := make(chan feed.FeedItem, 100)
	if err := feed.FeedSourceJSONStart(feedsData, channel); err != nil {
		errorAndDie(1, err.Error())
	}

	// go fetchFeedWithInterval("http://www.gsmarena.com/rss-news-reviews.php3", channel, 5)
	// go fetchFeedWithInterval("https://www.thehindubusinessline.com/feeder/default.rss", channel, 5)
	// go fetchFeedWithInterval("https://www.thehindubusinessline.com/news/national/feeder/default.rss", channel, 5)
	// go fetchFeedWithInterval("http://feeds.feedburner.com/ndtvnews-top-stories?format=xml", channel, 5)
	// go fetchFeedWithInterval("http://timesofindia.indiatimes.com/rssfeeds/1221656.cms", channel, 5)

	waitOnFeed(channel, &db)
}



// feedsData := `[
// 		{
// 			"name":"GSMArena",
// 			"url":"http://www.gsmarena.com/rss-news-reviews.php3",
// 			"interval":5,
// 			"kind":0
// 		},
// 		{
// 			"name":"TheHindu Buissness",
// 			"url":"https://www.thehindubusinessline.com/feeder/default.rss",
// 			"interval":5,
// 			"kind":0
// 		},
// 		{
// 			"name":"NDTV Top Stories",
// 			"url":"http://feeds.feedburner.com/ndtvnews-top-stories?format=xml",
// 			"interval":7,
// 			"kind":0
// 		},
// 		{
// 			"name":"Times of India",
// 			"url":"http://timesofindia.indiatimes.com/rssfeeds/1221656.cms",
// 			"interval":10,
// 			"kind":0
// 		}

// 	]
// 	`
