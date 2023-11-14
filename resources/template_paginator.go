package resources

import (
	"fmt"
)

type Page struct {
	display   string
	IsCurrent bool
}

func (p *Page) String() string {
	return p.display
}

type Paginator struct {
	Pages     [][]int
	ShowNear  int
	ShowStart int
	ShowEnd   int
}

func NewPaginator(showNear, showStart, showEnd int) *Paginator {
	return &Paginator{
		Pages:     make([][]int, 0),
		ShowNear:  showNear,
		ShowStart: showStart,
		ShowEnd:   showEnd,
	}
}

func (p *Paginator) AddPage(start, end int) *Paginator {
	p.Pages = append(p.Pages, []int{start, end})
	return p
}

func (p *Paginator) AddPages(size, step int) *Paginator {
	for i := 0; i < size; i++ {
		if i%step == 0 {
			p.AddPage(i+1, i+step)
		}
	}
	return p
}

func (p *Paginator) Paginate(current int) []*Page {
	out := make([]*Page, 0)
	max := len(p.Pages)
	prevWasDot := false
	for i := 0; i < max; i++ {
		pg := &Page{display: "..."}
		if i == current-1 {
			pg.IsCurrent = true
		}
		if i < p.ShowStart || i > (max-1)-p.ShowEnd || (i > ((current-2)-p.ShowNear) && i < current+p.ShowNear) {
			pg.display = fmt.Sprintf("%d", i+1)
			out = append(out, pg)
			prevWasDot = false
			continue
		}

		if prevWasDot {
			continue
		}

		out = append(out, pg)
		prevWasDot = true
	}
	return out
}

type PaginationResult struct {
	Pages   []*Page
	Start   int
	Stop    int
	Size    int
	HasPrev bool
	HasNext bool
}

type PgInput struct {
	// baseurl to work around some funky issues with browser pushstate
	BaseUrl string
	Max     int
	Step    int
	Current int
}

func NewPgInput(max, step, current int, baseUrl string) *PgInput {
	return &PgInput{baseUrl, max, step, current}
}

func CalculateMaxPageNum(size, step int) int {
	maxPage := size / step
	if size%step != 0 {
		maxPage++
	}
	return maxPage
}
