package pge

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

// Paginator provides pagination parameters for a select query.
type Paginator struct {
	// as in "get the first <n> results..."
	First int64
	// as in "get the last <n> results..."
	Last int64

	Cursors
}

// QueryParams returns query parameters which encode the pagination settings and
// offset.
func (p Paginator) QueryParams() url.Values {
	values := make(url.Values)

	if p.First != 0 {
		values.Set("first", strconv.FormatInt(p.First, 10))
	} else if p.Last != 0 {
		values.Set("last", strconv.FormatInt(p.Last, 10))
	}

	if p.AfterCursor != "" {
		values.Set("after", p.AfterCursor)
	} else if p.BeforeCursor != "" {
		values.Set("before", p.BeforeCursor)
	}

	return values
}

// Cursors can either be used as part of pagination input to a SELECT query
// (as part of a Paginator), or new cursors provided as a result of the SELECT.
type Cursors struct {
	// ...after <cursor> (optional)
	AfterCursor string
	// ...or before <cursor> (optional)
	BeforeCursor string
}

var removeTablePrefix = regexp.MustCompile(`^(.*)\.`)

func (q Query) WithPaginator(p *Paginator) (Query, error) {
	cursorCol := q.CursorColumn()
	if cursorCol == "" {
		return q, errors.New("cannot paginate query because it does not declare a cursor column")
	}

	reverseOrder := false
	ascOrDesc := q.AscendingOrDescending()
	// Get an extra row so we can tell whether to include a pagination cursor
	limit := p.First + 1
	gtOrLt := ">"
	if p.Last != 0 {
		limit = p.Last + 1
		ascOrDesc = ascOrDesc.Reverse()
		reverseOrder = true
	}
	if ascOrDesc == Descending {
		gtOrLt = "<"
	}

	paginationWhere := "true"
	if p.AfterCursor != "" {
		paginationWhere = fmt.Sprintf(" %s %s $%d", cursorCol, gtOrLt, q.NextParamNumber())
	} else if p.BeforeCursor != "" {
		paginationWhere = fmt.Sprintf(" %s %s $%d", cursorCol, gtOrLt, q.NextParamNumber())
	}
	q = q.WithTemplateParameter("PaginationWhere", paginationWhere)

	paginationOrderBy := fmt.Sprintf("%s %s LIMIT %d", cursorCol, ascOrDesc, limit)
	q = q.WithTemplateParameter("PaginationOrderBy", paginationOrderBy)

	if reverseOrder {
		// Wrap the query so the DB gives us the page in the expected sort order
		q = q.WithPrefix("SELECT * FROM (").
			WithSuffix(") unsorted ORDER BY unsorted." +
				removeTablePrefix.ReplaceAllLiteralString(cursorCol, "") + " " +
				string(q.AscendingOrDescending()))
	}
	q.paginator = p
	return q, nil
}

func extractCursors(p *Paginator, queryResults interface{}) (cursors Cursors, err error) {
	val := reflect.ValueOf(queryResults)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return Cursors{}, errors.New("query results are not a pointer to a slice")
	}

	l := val.Elem().Len()
	if l == 0 {
		return cursors, nil
	}

	if p.Last != 0 {
		if p.BeforeCursor != "" {
			cursors.AfterCursor, err = cursorFromValue(val.Elem().Index(l - 1))
			if err != nil {
				return
			}
		}
		if int64(val.Elem().Len()) > p.Last {
			cursors.BeforeCursor, err = cursorFromValue(val.Elem().Index(1))
			if err != nil {
				return
			}
			val.Elem().Set(val.Elem().Slice(1, l))
		}
	} else {
		if p.AfterCursor != "" {
			cursors.BeforeCursor, err = cursorFromValue(val.Elem().Index(0))
			if err != nil {
				return
			}
		}
		if int64(val.Elem().Len()) > p.First {
			cursors.AfterCursor, err = cursorFromValue(val.Elem().Index(l - 2))
			if err != nil {
				return
			}
			val.Elem().SetLen(l - 1)
		}
	}

	return cursors, nil
}

func cursorFromValue(v reflect.Value) (string, error) {
	cursorMethod := v.MethodByName("ToCursor")
	if !cursorMethod.IsValid() {
		return "", errors.New("missing ToCursor method")
	}
	returnVals := cursorMethod.Call(nil)
	if len(returnVals) != 1 || returnVals[0].Kind() != reflect.String {
		return "", errors.New("ToCursor did not return a string")
	}
	return returnVals[0].String(), nil
}
