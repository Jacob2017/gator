package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Request creation - %v", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request - %v", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading response - %v", err)
	}

	var feed RSSFeed
	if err = xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("Unmarshaling XML - %v", err)
	}

	pFeed := &feed

	pFeed.Channel.Title = html.UnescapeString(pFeed.Channel.Title)
	pFeed.Channel.Description = html.UnescapeString(pFeed.Channel.Description)

	for i := range pFeed.Channel.Item {
		pFeed.Channel.Item[i].Description = html.UnescapeString(pFeed.Channel.Item[i].Description)
		pFeed.Channel.Item[i].Title = html.UnescapeString(pFeed.Channel.Item[i].Title)
	}

	return pFeed, nil
}
