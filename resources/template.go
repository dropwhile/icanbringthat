package resources

import (
	"embed"
	"fmt"
	"os"

	//"html/template"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/google/safehtml/template"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates
var templateEmbedFS embed.FS

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
	Start   int
	Stop    int
	Size    int
	HasPrev bool
	HasNext bool
	Pages   []*Page
}

type PgInput struct {
	Max     int
	Step    int
	Current int
	// baseurl to work around some funky issues with browser pushstate
	BaseUrl string
}

func NewPgInput(max, step, current int, baseUrl string) *PgInput {
	return &PgInput{max, step, current, baseUrl}
}

func CalculateMaxPageNum(size, step int) int {
	maxPage := size / step
	if size%step != 0 {
		maxPage++
	}
	return maxPage
}

var templateFuncMap = template.FuncMap{
	"titlecase": cases.Title(language.English).String,
	"lowercase": func(s fmt.Stringer) string {
		return cases.Lower(language.English).String(s.String())
	},
	"formatTS": func(t time.Time) string {
		return t.UTC().Format("2006-01-02T15:04Z07:00")
	},
	"formatTSLocal": func(t time.Time, zone string) string {
		loc, _ := time.LoadLocation(zone)
		return t.In(loc).Format("2006-01-02T15:04")
	},
	"formatDateTime": func(t time.Time) string {
		return t.Format("2006-01-02 15:04 MST")
	},
	"paginate": func(pg *PgInput) *PaginationResult {
		size, step, current := pg.Max, pg.Step, pg.Current
		maxPage := CalculateMaxPageNum(size, step)
		if current < 1 {
			current = 1
		}
		if current > maxPage {
			current = maxPage
		}
		start := (step * (current - 1)) + 1
		stop := start + step - 1
		if current == maxPage && size%step != 0 {
			stop = ((current - 1) * step) + (size % step)
		}
		s := &PaginationResult{
			start, stop, size, current > 1, current < maxPage,
			NewPaginator(3, 3, 3).AddPages(size, step).Paginate(current),
		}
		return s
	},
	"mod": func(i, j int) int {
		return i % j
	},
	"add": func(i, j int) int {
		return i + j
	},
	"mult": func(i, j int) int {
		return i * j
	},
	"min": func(i, j int) int {
		return min(i, j)
	},
	"subt": func(i, j int) int {
		return i - j
	},
	"set": func(ac reflect.Value, k reflect.Value, v reflect.Value) error {
		switch ac.Kind() {
		case reflect.Map:
			if k.Type() == ac.Type().Key() {
				ac.SetMapIndex(k, v)
				return nil
			}
		}
		return fmt.Errorf("calling set with unsupported type %q (%T) -> %q (%T)", ac.Kind(), ac, k.Kind(), k)
	},
	// isset is a helper func from hugo
	"isset": func(ac reflect.Value, kv reflect.Value) (bool, error) {
		switch ac.Kind() {
		case reflect.Array, reflect.Slice:
			k := 0
			switch kv.Kind() {
			case reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64:
				k = int(kv.Int())
			case reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64:
				k = int(kv.Uint())
			case reflect.String:
				v, err := strconv.ParseInt(kv.String(), 0, 0)
				if err != nil {
					return false, fmt.Errorf("unable to cast %#v of type %T to int64", kv, kv)
				}
				k = int(v)
			default:
				return false, fmt.Errorf("unable to cast %#v of type %T to int", kv, kv)
			}
			if ac.Len() > k {
				return true, nil
			}
		case reflect.Map:
			if kv.Type() == ac.Type().Key() {
				return ac.MapIndex(kv).IsValid(), nil
			}
		default:
			log.Info().
				Msgf("calling IsSet with unsupported type %q (%T) will always return false", ac.Kind(), ac)
		}

		return false, nil
	},
	"eqorempty": func(arg0 reflect.Value, arg1 reflect.Value) (bool, error) {
		k1 := arg0.Kind()
		k2 := arg1.Kind()
		if k1 != k2 {
			return false, fmt.Errorf("non-comparable types %s: %v, %s: %v", arg0, arg0.Type(), arg1.Type(), arg1)
		}

		truth := false
		switch arg0.Kind() {
		case reflect.String:
			truth = arg0.String() == "" || arg0.String() == arg1.String()
		case reflect.Invalid:
			truth = true
		}
		return truth, nil
	},
}

type TemplateMap map[string]*template.Template

func (tm *TemplateMap) Get(name string) (*template.Template, error) {
	if v, ok := (*tm)[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("template not found for name %s", name)
}

func MustParseTemplates(templatesDir string) TemplateMap {
	templates := make(TemplateMap, 0)

	var templateFS template.TrustedFS
	var nameFS fs.FS
	if templatesDir == "embed" {
		var err error
		tfs := template.TrustedFSFromEmbed(templateEmbedFS)
		templateFS, err = tfs.Sub(template.TrustedSourceFromConstant("templates"))
		if err != nil {
			panic(err)
		}
		nameFS, err = fs.Sub(templateEmbedFS, "templates")
		if err != nil {
			panic(err)
		}
	} else {
		templateFS = template.TrustedFSFromTrustedSource(template.TrustedSourceFromEnvVar("TPL_DIR"))
		nameFS = os.DirFS(templatesDir)
	}

	nonViewTemplates, err := template.New("").Funcs(templateFuncMap).ParseFS(
		templateFS,
		"layout/*.gohtml",
		"partial/*.gohtml",
	)
	if err != nil {
		panic(err)
	}

	viewSub, err := fs.Sub(nameFS, "view")
	if err != nil {
		panic(err)
	}

	err = fs.WalkDir(viewSub, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".gohtml" {
			name := filepath.Base(p)
			c, err := nonViewTemplates.Clone()
			if err != nil {
				panic(err)
			}
			t, err := c.New(name).Funcs(templateFuncMap).ParseFS(
				templateFS, fmt.Sprintf("view/%s", name),
			)
			if err != nil {
				panic(err)
			}
			templates[name] = t
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return templates
}
