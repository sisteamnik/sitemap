package sitemap

import (
	"fmt"
	"time"
)

type Item struct {
	Loc        string    `xml:"loc"`
	LastMod    time.Time `xml:"lastmod"`
	ChangeFreq string    `xml:"changefreq,omitempty"`
	Priority   float32   `xml:"priority,omitempty"`

	isIndex bool `xml:"-"`
}

func NewItem(loc string, lastMod time.Time, changeFreq string, priority float32) Item {
	return Item{
		Loc:        loc,
		LastMod:    lastMod,
		ChangeFreq: changeFreq,
		Priority:   priority,
	}
}

func NewIndexItem(loc string, lastMod time.Time) Item {
	return Item{
		Loc:     loc,
		LastMod: lastMod,

		isIndex: true,
	}
}

func (item Item) String() string {
	//2012-08-30T01:23:57+08:00
	//Mon Jan 2 15:04:05 -0700 MST 2006
	if !item.isIndex {
		return fmt.Sprintf(template, item.Loc, item.LastMod.Format(TimeFormat), item.ChangeFreq, item.Priority)
	} else {
		return fmt.Sprintf(indexTemplate, "", item.Loc, item.LastMod.Format(TimeFormat))
	}
}
