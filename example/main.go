package main

import (
	"encoding/xml"
	"fmt"
	"github.com/sisteamnik/sitemap"
	"io/ioutil"
	"time"
)

type UrlSet struct {
	Ursl []Url `xml:"url"`
}

type Url struct {
	Loc        string    `xml:"loc"`
	Lastmod    time.Time `xml:"lastmod"`
	Changefreq string    `xml:"changefreq"`
	Priority   float32   `xml:"priority"`
}

var sl = "sitemap.xml"

func main() {
	err := sitemap.Add(sl, &sitemap.Item{
		Loc:        "http://golang.org",
		LastMod:    time.Now(),
		Changefreq: "weekly",
		Priority:   0.5,
	})
	if err != nil {
		fmt.Println(err)
	}
	PrintSitemap()
	time.Sleep(time.Second * 3)
	err = sitemap.Update(sl, &sitemap.Item{
		Loc:        "http://golang.org",
		LastMod:    time.Now(),
		Changefreq: "weekly",
		Priority:   0.5,
	})
	if err != nil {
		fmt.Println(err)
	}
	PrintSitemap()

	err = sitemap.Delete(sl, &sitemap.Item{
		Loc:        "http://golang.org",
		LastMod:    time.Now(),
		Changefreq: "weekly",
		Priority:   0.5,
	})
	if err != nil {
		fmt.Println(err)
	}
	PrintSitemap()

}

func PrintSitemap() {
	f, err := ioutil.ReadFile(sl)
	if err != nil {
		fmt.Println(err)
	}
	v := UrlSet{}
	err = xml.Unmarshal(f, &v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v.Ursl)
}
