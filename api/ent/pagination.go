// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/dopedao/RYO/api/ent/turn"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vmihailenco/msgpack/v5"
)

// OrderDirection defines the directions in which to order a list of items.
type OrderDirection string

const (
	// OrderDirectionAsc specifies an ascending order.
	OrderDirectionAsc OrderDirection = "ASC"
	// OrderDirectionDesc specifies a descending order.
	OrderDirectionDesc OrderDirection = "DESC"
)

// Validate the order direction value.
func (o OrderDirection) Validate() error {
	if o != OrderDirectionAsc && o != OrderDirectionDesc {
		return fmt.Errorf("%s is not a valid OrderDirection", o)
	}
	return nil
}

// String implements fmt.Stringer interface.
func (o OrderDirection) String() string {
	return string(o)
}

// MarshalGQL implements graphql.Marshaler interface.
func (o OrderDirection) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(o.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (o *OrderDirection) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("order direction %T must be a string", val)
	}
	*o = OrderDirection(str)
	return o.Validate()
}

func (o OrderDirection) reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

func (o OrderDirection) orderFunc(field string) OrderFunc {
	if o == OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

func cursorsToPredicates(direction OrderDirection, after, before *Cursor, field, idField string) []func(s *sql.Selector) {
	var predicates []func(s *sql.Selector)
	if after != nil {
		if after.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeGT
			} else {
				predicate = sql.CompositeLT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					after.Value, after.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.GT
			} else {
				predicate = sql.LT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					after.ID,
				))
			})
		}
	}
	if before != nil {
		if before.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeLT
			} else {
				predicate = sql.CompositeGT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					before.Value, before.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.LT
			} else {
				predicate = sql.GT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					before.ID,
				))
			})
		}
	}
	return predicates
}

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID    int   `msgpack:"i"`
	Value Value `msgpack:"v,omitempty"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	quote := []byte{'"'}
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := msgpack.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(s),
		),
	).Decode(c); err != nil {
		return fmt.Errorf("cannot decode cursor: %w", err)
	}
	return nil
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func getCollectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	field := fc.Field

walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Name == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return getCollectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

// TurnEdge is the edge representation of Turn.
type TurnEdge struct {
	Node   *Turn  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// TurnConnection is the connection containing edges to Turn.
type TurnConnection struct {
	Edges      []*TurnEdge `json:"edges"`
	PageInfo   PageInfo    `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// TurnPaginateOption enables pagination customization.
type TurnPaginateOption func(*turnPager) error

// WithTurnOrder configures pagination ordering.
func WithTurnOrder(order *TurnOrder) TurnPaginateOption {
	if order == nil {
		order = DefaultTurnOrder
	}
	o := *order
	return func(pager *turnPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultTurnOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithTurnFilter configures pagination filter.
func WithTurnFilter(filter func(*TurnQuery) (*TurnQuery, error)) TurnPaginateOption {
	return func(pager *turnPager) error {
		if filter == nil {
			return errors.New("TurnQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type turnPager struct {
	order  *TurnOrder
	filter func(*TurnQuery) (*TurnQuery, error)
}

func newTurnPager(opts []TurnPaginateOption) (*turnPager, error) {
	pager := &turnPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultTurnOrder
	}
	return pager, nil
}

func (p *turnPager) applyFilter(query *TurnQuery) (*TurnQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *turnPager) toCursor(t *Turn) Cursor {
	return p.order.Field.toCursor(t)
}

func (p *turnPager) applyCursors(query *TurnQuery, after, before *Cursor) *TurnQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultTurnOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *turnPager) applyOrder(query *TurnQuery, reverse bool) *TurnQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultTurnOrder.Field {
		query = query.Order(direction.orderFunc(DefaultTurnOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Turn.
func (t *TurnQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...TurnPaginateOption,
) (*TurnConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newTurnPager(opts)
	if err != nil {
		return nil, err
	}

	if t, err = pager.applyFilter(t); err != nil {
		return nil, err
	}

	conn := &TurnConnection{Edges: []*TurnEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := t.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := t.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	t = pager.applyCursors(t, after, before)
	t = pager.applyOrder(t, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		t = t.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		t = t.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := t.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Turn
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Turn {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Turn {
			return nodes[i]
		}
	}

	conn.Edges = make([]*TurnEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &TurnEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

var (
	// TurnOrderFieldUserID orders Turn by user_id.
	TurnOrderFieldUserID = &TurnOrderField{
		field: turn.FieldUserID,
		toCursor: func(t *Turn) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.UserID,
			}
		},
	}
	// TurnOrderFieldLocationID orders Turn by location_id.
	TurnOrderFieldLocationID = &TurnOrderField{
		field: turn.FieldLocationID,
		toCursor: func(t *Turn) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.LocationID,
			}
		},
	}
	// TurnOrderFieldItemID orders Turn by item_id.
	TurnOrderFieldItemID = &TurnOrderField{
		field: turn.FieldItemID,
		toCursor: func(t *Turn) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.ItemID,
			}
		},
	}
	// TurnOrderFieldCreatedAt orders Turn by created_at.
	TurnOrderFieldCreatedAt = &TurnOrderField{
		field: turn.FieldCreatedAt,
		toCursor: func(t *Turn) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.CreatedAt,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f TurnOrderField) String() string {
	var str string
	switch f.field {
	case turn.FieldUserID:
		str = "USER_ID"
	case turn.FieldLocationID:
		str = "LOCATION_ID"
	case turn.FieldItemID:
		str = "ITEM_ID"
	case turn.FieldCreatedAt:
		str = "CREATED_AT"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f TurnOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *TurnOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("TurnOrderField %T must be a string", v)
	}
	switch str {
	case "USER_ID":
		*f = *TurnOrderFieldUserID
	case "LOCATION_ID":
		*f = *TurnOrderFieldLocationID
	case "ITEM_ID":
		*f = *TurnOrderFieldItemID
	case "CREATED_AT":
		*f = *TurnOrderFieldCreatedAt
	default:
		return fmt.Errorf("%s is not a valid TurnOrderField", str)
	}
	return nil
}

// TurnOrderField defines the ordering field of Turn.
type TurnOrderField struct {
	field    string
	toCursor func(*Turn) Cursor
}

// TurnOrder defines the ordering of Turn.
type TurnOrder struct {
	Direction OrderDirection  `json:"direction"`
	Field     *TurnOrderField `json:"field"`
}

// DefaultTurnOrder is the default ordering of Turn.
var DefaultTurnOrder = &TurnOrder{
	Direction: OrderDirectionAsc,
	Field: &TurnOrderField{
		field: turn.FieldID,
		toCursor: func(t *Turn) Cursor {
			return Cursor{ID: t.ID}
		},
	},
}

// ToEdge converts Turn into TurnEdge.
func (t *Turn) ToEdge(order *TurnOrder) *TurnEdge {
	if order == nil {
		order = DefaultTurnOrder
	}
	return &TurnEdge{
		Node:   t,
		Cursor: order.Field.toCursor(t),
	}
}