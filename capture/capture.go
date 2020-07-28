package capture

import (
	"context"	
	"fmt"
	"time"

	
	"capture-feed/scrapper"
	"capture-feed/utility"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
)

type CaptureStateType int

const (
	CaptureStateTypeCooking = iota
	CaptureStateTypeReady
)

type CaptureItem struct {
	Title       string           `json:"title"`
	URL         string           `json:"url"`
	Description string           `json:"description"`
	ImgURL      string           `json:"img_url"`
	PubDate     time.Time        `json:"pub_date"`
	State       CaptureStateType `json:"state"`
	Hash        string           `json:"hash"`
}

func CaptureItemNewFromFeedItem(feed *gofeed.Item) CaptureItem {
	var imgURL string
	if feed.Image != nil {
		imgURL = feed.Image.URL
	}
	capture := CaptureItem{
		Title:       feed.Title,
		URL:         feed.Link,
		Description: feed.Description,
		ImgURL:      imgURL,
		PubDate:     *feed.PublishedParsed,
		State:       CaptureStateTypeCooking,
	}

	capture.hashUpdate()
	return capture
}

func CaptureUpdateArticle(db *DatabaseContext, article *scrapper.SCNewsArticle) error {

	var des string
	var img_url string
	for _, val := range article.Paragraphs {
		des = val.Content
		break
	}

	for _, val := range article.Posters {
		img_url = val.URL
		break
	}



	collection := db.client.Database("capture").Collection("news")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.D{
			{"hash", article.Hash},
		},
		bson.D{
			{"$set", bson.D{
				{"title", article.Headline},
				{"description", des},
				{"state", CaptureStateTypeReady},
				{"pub_date", article.PubDate}, 
				{"img_url", img_url},
			}},
			{"$currentDate", bson.D{
				{"lastModified", true},
			}},
		},
	)

	return err
}


func (item *CaptureItem) hashUpdate() {
	item.Hash = utility.GetHash([]byte(item.URL))
}

func (item *CaptureItem) Log() {
	fmt.Printf("%s\t\t%s\n", item.Hash, item.Title)
}

func (item *CaptureItem) Process() {
	// Process
	
}

func (item *CaptureItem) Exits(db *DatabaseContext) (bool, error) {

	var capture CaptureItem

	collection := db.client.Database("capture").Collection("news")

	filter := bson.M{"hash": item.Hash}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&capture)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (item *CaptureItem) Save(db *DatabaseContext) error {
	collection := db.client.Database("capture").Collection("news")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, bson.M{
		"title":       item.Title,
		"url":         item.URL,
		"description": item.Description,
		"img_url":     item.ImgURL,
		"pub_date":    item.PubDate,
		"hash":        item.Hash,
		"state":       item.State,
	})
	if err != nil {
		return err
	}

	return nil
}
