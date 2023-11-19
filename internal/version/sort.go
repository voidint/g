package version

type Collection []*Version

func (c Collection) Len() int {
	return len(c)
}

func (c Collection) Less(i, j int) bool {
	return c[i].sv.LessThan(c[j].sv)
}

func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
