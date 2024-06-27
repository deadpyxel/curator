package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/deadpyxel/curator/internal/database"
)

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description string        `xml:"description"`
	Language    string        `xml:"language"`
	Item        []RSSFeedItem `xml:"item"`
}

type RSSFeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}

// startFeedScrapping initiates the scraping operation with the specified parameters.
// It calls scrapeFeed to fetch feeds from the database with the given concurrency and time interval between requests.
func startFeedScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	logger.Info("Starting scrape operation", "concurrency", concurrency, "interval", timeBetweenRequest.String())

	ticker := time.NewTicker(timeBetweenRequest)
	for range ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			logger.Error("Error fetching feeds", "error", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

// scrapeFeed fetches and processes the feed data.
// It fetches the feed data, marks the feed as fetched in the database and logs the new posts found.
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		logger.Error("Error marking feed as fetched", "feedID", feed.ID, "feedName", feed.Name)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		logger.Error("Error fetching feed data", "feedID", feed.ID, "error", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		logger.Info(fmt.Sprintf("New post found on %s with title %v", feed.Name, item.Title), "feedID", feed.ID)
	}
	logger.Info("Feed scrapping complete", "feedID", feed.ID, "numPosts", len(rssFeed.Channel.Item))
}
