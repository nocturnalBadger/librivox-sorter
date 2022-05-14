package main

import (
	"fmt"
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

    fmt.Println("I am alive!")
	http.ListenAndServe(":3333", r)
}

func GetRSSFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Println("i am a teapot")

	fmt.Println(idInt)
	feed, err := GetLibrivoxFeed(idInt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println(feed.Root())
	SortFeedItems(feed)
}

func GetLibrivoxFeed(feedID int) (*etree.Document, error) {
	url := fmt.Sprintf("%s%d", librivoxFeedURL, feedID)
	fmt.Println(url)
	fmt.Println("i am a teapot")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("i am a teapot")

	fmt.Println(resp)
	fmt.Println(string(body))

	doc := etree.NewDocument()
	err = doc.ReadFromBytes(body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func SortFeedItems(feed *etree.Document) {
	for _, t := range feed.FindElements("//item") {
		fmt.Println("Title:", t.Text())
	}
}
