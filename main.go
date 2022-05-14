package main

import (
	"fmt"
	"time"
	"github.com/beevik/etree"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strconv"
)

const librivoxFeedURL = "https://librivox.org/rss/"

func main() {

	r := chi.NewRouter()
	r.Get("/rss/{id}", GetRSSFeed)

	http.ListenAndServe(":3333", r)
}

func GetRSSFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Println(idInt)
	feed, err := GetLibrivoxFeed(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println(feed.Root())
	SortFeedItems(feed)

	feed.WriteTo(w)
}

func GetLibrivoxFeed(feedID int) (*etree.Document, error) {
	url := fmt.Sprintf("%s%d", librivoxFeedURL, feedID)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	err = doc.ReadFromBytes(body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func SortFeedItems(feed *etree.Document) {
	startTime := time.Now()
	for i, t := range feed.FindElements("//item") {
		pubDate := t.CreateElement("pubDate")
		time := startTime.Add(time.Minute * time.Duration(i))
		pubDate.SetText(time.Format("2006-01-02 15:04:05.000000"))
	}
}
