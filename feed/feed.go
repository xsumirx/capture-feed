package feed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedType int

const (
	FeedTypeRSS = iota
	FeedTypeScrapper
)

type FeedContext struct {
	name           string
	url            string
	kind           FeedType
	repeatInterval int

	

	//Reply Channel
	channel chan FeedItem
}

type FeedItem struct {
	Code int
	Err  string

	Name  string
	Items []*gofeed.Item
}

func (feedItem *FeedItem) send(channel chan FeedItem) {
	channel <- *feedItem
}

func (feedItem *FeedItem) FeedItemPrintStatus() {
	if feedItem.Code != 0 {
		fmt.Printf("[FAIL]\t%s\t%s\n", feedItem.Name, feedItem.Err)
	} else {
		fmt.Printf("[OK]\t%d\t%s\n", len(feedItem.Items), feedItem.Name)
	}
}

func (feed *FeedContext) start() {
	go feed.feedFetchWithInterval()
}

func feedSendErrorOnChannel(code int, err string, channel chan FeedItem) {
	errFeedItem := FeedItem{
		Code:  1,
		Err:   err,
		Items: nil,
	}

	errFeedItem.send(channel)
}

func feedSendOnChannel(feed *gofeed.Feed, channel chan FeedItem) {
	errFeedItem := FeedItem{
		Code:  0,
		Err:   "Success",
		Name:  feed.Title,
		Items: feed.Items,
	}

	errFeedItem.send(channel)
}

func (feed *FeedContext) feedFetch() {
	//Fetch the details of the url

	response, err := http.Get(feed.url)
	if err != nil {
		feedSendErrorOnChannel(1, err.Error(), feed.channel)
		return
	}

	//Ready all the content
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		feedSendErrorOnChannel(1, err.Error(), feed.channel)
		return
	}

	//Parse
	fp := gofeed.NewParser()
	feedObj, err := fp.ParseString(string(body))
	if err != nil {
		feedSendErrorOnChannel(1, err.Error(), feed.channel)
		return
	}

	feedSendOnChannel(feedObj, feed.channel)

}

func (feed *FeedContext) feedFetchWithInterval() {
	for {
		feed.feedFetch()
		time.Sleep(time.Second * time.Duration(feed.repeatInterval))
	}
}

type FeedConfig struct {
	Name     string
	Url      string
	Kind     int
	Interval int
}

func FeedSourceJSONStart(feedConfigsJSON string, channel chan FeedItem) error {

	var feedConfigs []FeedConfig
	if err := json.Unmarshal([]byte(feedConfigsJSON), &feedConfigs); err != nil {
		return err
	}

	for _, feedConfig := range feedConfigs {
		//Make a new FeedContext and run it
		feedCtx := FeedContext{
			name:           feedConfig.Name,
			url:            feedConfig.Url,
			kind:           FeedType(feedConfig.Kind),
			repeatInterval: feedConfig.Interval,
			channel:        channel,
		}

		feedCtx.start()
	}

	return nil

}
