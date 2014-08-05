package sitemap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	footer   = `	</urlset>`
	template = `
	 <url>
	   <loc>%s</loc>
	   <lastmod>%s</lastmod>
	   <changefreq>%s</changefreq>
	   <priority>%.1f</priority>
	 </url> 	`

	indexHeader = `<?xml version="1.0" encoding="UTF-8"?>
      <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexFooter = `
</sitemapindex>
	`
	indexTemplate = `
    <sitemap>
       <loc>%s%s</loc>
       <lastmod>%s</lastmod>
    </sitemap>
	`
)

type UrlSet struct {
	Ursl []*Item `xml:"url"`
}

type Item struct {
	Loc        string    `xml:"loc"`
	LastMod    time.Time `xml:"lastmod"`
	Changefreq string    `xml:"changefreq"`
	Priority   float32   `xml:"priority"`
}

func (item *Item) String() string {
	//2012-08-30T01:23:57+08:00
	//Mon Jan 2 15:04:05 -0700 MST 2006
	return fmt.Sprintf(template, item.Loc, item.LastMod.Format("2006-01-02T15:04:05+08:00"), item.Changefreq, item.Priority)
}

func SiteMap(f string, items []*Item) error {
	var buffer bytes.Buffer
	buffer.WriteString(header)
	for _, item := range items {
		if item == nil {
			continue
		}
		_, err := buffer.WriteString(item.String())
		if err != nil {
			return err
		}
	}
	fo, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fo.Close()
	buffer.WriteString("\n" + footer)

	_, err = fo.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return err
}

func SiteMapIndex(folder, indexFile, baseurl string) error {
	var buffer bytes.Buffer
	buffer.WriteString(indexHeader)
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml") {
			fmt.Println(f.Name())
			s := fmt.Sprintf(indexTemplate, baseurl, f.Name(), time.Now().Format("2006-01-02T15:04:05+08:00"))
			//fmt.Println(s)
			buffer.WriteString(s)
		}
	}
	buffer.WriteString(indexFooter)
	err = ioutil.WriteFile(indexFile, buffer.Bytes(), 0755)
	return err
}

func Add(f string, item *Item) error {
	if ExistLoc(f, item.Loc) {
		return errors.New("Location already exist")
	}
	fi, err := os.Stat(f)
	if fi == nil && err != nil {
		return SiteMap(f, []*Item{item})
	}
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	lines := string(content)
	xml := strings.Replace(lines, "\n"+footer, item.String()+"\n"+footer, 1)
	err = ioutil.WriteFile(f, []byte(xml), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Update(fp string, item *Item) error {
	if !ExistLoc(fp, item.Loc) {
		return errors.New("Location not exist")
	}
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	v := UrlSet{}
	err = xml.Unmarshal(f, &v)
	if err != nil {
		return err
	}
	for _, v := range v.Ursl {
		if item.Loc == v.Loc {
			v.LastMod = item.LastMod
		}
	}
	SiteMap(fp, v.Ursl)
	return nil
}

func Delete(fp string, item *Item) error {
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	v := UrlSet{}
	err = xml.Unmarshal(f, &v)
	if err != nil {
		return err
	}
	for i, k := range v.Ursl {
		if item.Loc == k.Loc {
			v.Ursl[i] = nil
		}
	}
	SiteMap(fp, v.Ursl)
	return nil
}

func ExistLoc(fp string, l string) bool {
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return false
	}
	loc := "<loc>" + l + "</loc>"
	if strings.Index(string(f), loc) != -1 {
		return true
	}
	return false
}
