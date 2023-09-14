package refid

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

var (
	refTagTest     = byte(1)
	testValWoutTag = "000baxr70ja4ggc0jbgw5dzx7vb52"
	testValWithTag = "040baxr70ja4ggc0jbgw5dzx7vb52"
)

func UintToBinary(n uint64) string {
	return strconv.FormatUint(n, 2)
}

func IntToBinary(n int64) string {
	return strconv.FormatInt(n, 2)
}

func TestGetTime(t *testing.T) {
	t.Parallel()

	// divide times by 10, so we are close enough
	t0 := time.Now().UTC().Unix() / 10
	refId := MustNew()
	vz := refId.Time().UTC().Unix() / 10
	assert.Equal(t, t0, vz)
}

func TestBase64RoundTrip(t *testing.T) {
	t.Parallel()

	refId := MustParse(testValWithTag)
	b64 := refId.ToBase64String()
	refId2, err := FromBase64String(b64)
	assert.NilError(t, err)
	assert.Equal(t, refId.String(), refId2.String())
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()
	u := MustNew()
	r := MustParse(u.String())
	assert.Check(t, !u.HasTag(refTagTest))
	assert.Check(t, !r.HasTag(refTagTest))
	assert.Equal(t, u.String(), r.String())

	u = MustNewTagged(refTagTest)
	r = MustParse(u.String())
	assert.Check(t, u.HasTag(refTagTest))
	assert.Check(t, r.HasTag(refTagTest))
	assert.Equal(t, u.String(), r.String())
}

func TestSetTag(t *testing.T) {
	t.Parallel()

	refId := MustParse(testValWoutTag)
	assert.Check(t, !refId.HasTag(refTagTest))
	assert.Equal(t, refId.String(), testValWoutTag)
	assert.Equal(t, (&refId).String(), testValWoutTag)

	refId.SetTag(refTagTest)
	assert.Check(t, refId.HasTag(refTagTest))
	assert.Equal(t, refId.String(), testValWithTag)
	assert.Equal(t, (&refId).String(), testValWithTag)

	refId.ClearTag()
	assert.Check(t, !refId.HasTag(refTagTest))
	assert.Equal(t, refId.String(), testValWoutTag)
	assert.Equal(t, (&refId).String(), testValWoutTag)
}

func TestAmbiguous(t *testing.T) {
	t.Parallel()

	rd0 := MustParse(testValWoutTag)
	rd1 := MustParse(testValWoutTag)
	rd2 := MustParse(testValWoutTag)
	assert.Assert(t,
		rd0.String() == rd1.String() && rd1.String() == rd2.String(),
	)
}

func TestTemplateStringer(t *testing.T) {
	t.Parallel()
	s := MustParse(testValWoutTag)
	assert.Equal(t, fmt.Sprintf("%s", s), testValWoutTag)
	tpl := template.Must(template.New("name").Parse(`{{.}}`))
	var b bytes.Buffer
	err := tpl.Execute(&b, s)
	assert.NilError(t, err)
	assert.Equal(t, b.String(), testValWoutTag)
}
