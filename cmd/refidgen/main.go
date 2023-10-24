package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"strings"
)

const tplText = `
// Code generated by refidgen. DO NOT EDIT.
// generated from: {{.Origin}}

package {{.Pkg}}

import (
	"fmt"

	"github.com/dropwhile/refid"
)

const tagVal{{.Name}} = {{.Value}}

type {{.Name}} struct {
	refid.RefID
}

func (r *{{.Name}}) checkResult(err error) error {
	if err != nil {
		return err
	}
	if !r.RefID.HasTag(tagVal{{.Name}}) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}

func (r *{{.Name}}) Scan(src interface{}) error {
	err := r.RefID.Scan(src)
	return r.checkResult(err)
}

func (r *{{.Name}}) UnmarshalJSON(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func (r *{{.Name}}) UnmarshalBinary(b []byte) error {
	err := r.RefID.UnmarshalBinary(b)
	return r.checkResult(err)
}

func New{{.Name}}() ({{.Name}}, error) {
	v, err := refid.NewTagged(tagVal{{.Name}})
	return {{.Name}}{v}, err
}

func Parse{{.Name}}(s string) ({{.Name}}, error) {
	v, err := refid.ParseWithRequire(s, refid.HasTag(tagVal{{.Name}}))
	return {{.Name}}{v}, err
}

func Parse{{.Name}}WithRequire(s string, reqs ...refid.Requirement) ({{.Name}}, error) {
	reqs = append(reqs, refid.HasTag(tagVal{{.Name}}))
	v, err := refid.ParseWithRequire(s, reqs...)
	return {{.Name}}{v}, err
}

func {{.Name}}FromBytes(input []byte) ({{.Name}}, error) {
	var r {{.Name}}
	err := r.UnmarshalBinary(input)
	return r, err
}

type Null{{.Name}} struct {
	refid.NullRefID
}

func (u *Null{{.Name}}) checkResult(err error) error {
	if err != nil {
		return err
	}
	n := u.NullRefID
	if n.Valid && !n.RefID.HasTag(tagVal{{.Name}}) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}


func (u *Null{{.Name}}) Scan(src interface{}) error {
	err := u.NullRefID.Scan(src)
	return u.checkResult(err)
}

func (u *Null{{.Name}}) UnmarshalJSON(b []byte) error {
	err := u.NullRefID.UnmarshalJSON(b)
	return u.checkResult(err)
}

type {{.Name}}Matcher struct{}

func (a {{.Name}}Matcher) Match(v interface{}) bool {
	var r refid.RefID
	var err error
	switch x := v.(type) {
	case {{.Name}}:
		r = x.RefID
	case *{{.Name}}:
		r = x.RefID
	case string:
		r, err = refid.Parse(x)
	case []byte:
		r, err = refid.FromBytes(x)
	default:
		return false
	}
	if err != nil {
		return false
	}
	return r.HasTag(tagVal{{.Name}})
}
`

func main() {
	var (
		output  string
		prefix  string
		suffix  string
		typeVal int
	)
	flag.StringVar(&prefix, "t", "", "type prefix")
	flag.StringVar(&suffix, "s", ".gen.go", "file prefix")
	flag.IntVar(&typeVal, "v", 0, "type value")
	flag.Parse()

	fname := strings.TrimSuffix(os.Getenv("GOFILE"), ".go")
	pkg := os.Getenv("GOPACKAGE")

	if prefix == "" {
		log.Fatal("Param prefix is required")
	}

	if suffix == "" {
		log.Fatal("Param suffix is required")
	}

	if typeVal <= 0 {
		log.Fatal("Param value is required")
	}

	output = fmt.Sprintf("%s_refid%s", fname, suffix)
	fmt.Printf("generating %s\n", path.Base(output))

	t, err := template.New("fileTemplate").Parse(strings.TrimLeft(tplText, "\n"))
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	writer := bufio.NewWriter(w)
	defer writer.Flush()

	data := map[string]any{
		"Name":   fmt.Sprintf("%sRefID", prefix),
		"Value":  typeVal,
		"Origin": os.Getenv("GOFILE"),
		"Pkg":    pkg,
	}
	err = t.Execute(writer, data)
	if err != nil {
		log.Fatal(err)
	}
}