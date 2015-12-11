package sitemap

type ByTime Items

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	return a[i].LastMod.UnixNano() <
		a[j].LastMod.UnixNano()
}
