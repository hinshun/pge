package pge

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var (
	spacePattern = regexp.MustCompile(`\s+`)

	removeCommentsRegexp = regexp.MustCompile(`((?m)--.*$)`)
)

type QueryOption func(*Query)

func WithInsertColumns(cols int) QueryOption {
	return func(q *Query) {
		q.insertCols = cols
	}
}

func WithSuffix(suffix string) QueryOption {
	return func(q *Query) {
		q.suffix = suffix
	}
}

func WithCursorColumn(column string, paramNumber int, ascendingOrDescending Order) QueryOption {
	return func(q *Query) {
		q.cursorCol = column
		q.nextParamNumber = paramNumber
		q.ascendingOrDescending = ascendingOrDescending
	}
}

func WithCursorFromString(cursorFromString func(string) (interface{}, error)) QueryOption {
	return func(q *Query) {
		q.cursorFromString = cursorFromString
	}
}

type Query struct {
	Name                  string

	query                 string
	insertCols            int
	prefix                string
	suffix                string
	cursorCol             string
	nextParamNumber       int
	cursorFromString      func(string) (interface{}, error)
	ascendingOrDescending Order
	paginator             *Paginator
	tmpl                  *template.Template
	templateParams        map[string]interface{}
}

func NewQuery(name, query string, opts ...QueryOption) Query {
	q := Query{
		Name:  name,
		query: query,
	}
	q.query = removeCommentsRegexp.ReplaceAllString(query, "")
	for _, opt := range opts {
		opt(&q)
	}
	q.tmpl = template.Must(template.New("query").Parse(q.query))
	return q
}

func (q Query) String() string {
	if q.insertCols > 0 {
		q = q.WithValues(1)
	}
	var buf bytes.Buffer
	err := q.tmpl.Execute(&buf, q.templateParams)
	if err != nil {
		return trimString(q.prefix + q.query + q.suffix)
	}
	return trimString(q.prefix + buf.String() + q.suffix)
}

func (q Query) CursorColumn() string {
	return q.cursorCol
}

func (q Query) NextParamNumber() int {
	return q.nextParamNumber
}

func (q Query) AscendingOrDescending() Order {
	return q.ascendingOrDescending
}

func trimString(qStr string) string {
	return spacePattern.ReplaceAllString(strings.TrimSpace(qStr), " ")
}

func (q Query) WithPrefix(prefix string) Query {
	q.prefix = prefix + "\n" + q.prefix
	return q
}

func (q Query) WithSuffix(suffix string) Query {
	q.suffix = q.suffix + "\n" + suffix
	return q
}

func (q Query) WithTemplateParameter(param string, paramValue interface{}) Query {
	updatedParams := make(map[string]interface{})
	for k, v := range q.templateParams {
		updatedParams[k] = v
	}
	updatedParams[param] = paramValue
	q.templateParams = updatedParams
	return q
}

func (q Query) WithValues(rows int) Query {
	if rows == 0 {
		return q
	}

	cols := q.insertCols

	var b strings.Builder

	// For each row, we write `(),` = 3*rows
	// For each col, we write `$,` = 2*rows*cols
	query := spacePattern.ReplaceAllString(strings.TrimSpace(q.query), " ")
	prefix := spacePattern.ReplaceAllString(strings.TrimSpace(q.prefix), " ")
	suffix := spacePattern.ReplaceAllString(strings.TrimSpace(q.suffix), " ")
	length := len(prefix) + len(query) + 3*rows + 2*rows*cols + numDigits(rows*cols) + len(suffix)
	b.Grow(length)

	b.WriteString(prefix)
	b.WriteString(query)
	b.WriteString(" VALUES ")

	for i := 0; i < rows; i++ {
		b.WriteString("(")
		for j := 0; j < cols; j++ {
			b.WriteString("$")
			b.WriteString(strconv.Itoa(i*cols + j + 1))
			if j != cols-1 {
				b.WriteString(",")
			}
		}
		b.WriteString(")")
		if i != rows-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(" ")
	b.WriteString(suffix)

	return NewQuery(q.Name, b.String())
}

// CursorFromString converts a cursor from a string to the format we want to
// compare in the DB query.
func (q Query) CursorFromString(cursor string) (interface{}, error) {
	if q.cursorFromString == nil {
		return cursor, nil
	}
	return q.cursorFromString(cursor)
}

type Order string

const (
	Ascending  Order = "ASC"
	Descending Order = "DESC"
)

func (o Order) Reverse() Order {
	if o == Ascending {
		return Descending
	}
	return Ascending
}

func numDigits(n int) int {
	sum := 0
	for i := 1; i <= n; i *= 10 {
		sum += n - i + 1
	}
	return sum
}
