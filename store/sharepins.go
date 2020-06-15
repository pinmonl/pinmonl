package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Sharepins struct {
	*Store
}

type SharepinOpts struct {
	ListOpts
	ShareIDs []string
	PinlIDs  []string
	Status   field.NullValue

	PinlQuery string
	joinPinls bool

	TagIDs []string
}

func NewSharepins(s *Store) *Sharepins {
	return &Sharepins{s}
}

func (s Sharepins) table() string {
	return "sharepins"
}

func (s *Sharepins) List(ctx context.Context, opts *SharepinOpts) (model.SharepinList, error) {
	if opts == nil {
		opts = &SharepinOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table())
	qb = s.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Sharepin, 0)
	for rows.Next() {
		sharepin, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, sharepin)
	}
	return list, nil
}

func (s *Sharepins) Count(ctx context.Context, opts *SharepinOpts) (int64, error) {
	if opts == nil {
		opts = &SharepinOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select("count(*)").From(s.table())
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Sharepins) Find(ctx context.Context, id string) (*model.Sharepin, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	sharepin, err := s.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return sharepin, nil
}

func (s *Sharepins) FindOrCreate(ctx context.Context, data *model.Sharepin) (*model.Sharepin, error) {
	found, err := s.List(ctx, &SharepinOpts{
		ShareIDs: []string{data.ShareID},
		PinlIDs:  []string{data.PinlID},
	})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return found[0], nil
	}

	sharepin := *data
	err = s.Create(ctx, &sharepin)
	if err != nil {
		return nil, err
	}
	return &sharepin, nil
}

func (s *Sharepins) ListWithPinl(ctx context.Context, opts *SharepinOpts) (model.SharepinList, error) {
	if opts == nil {
		opts = &SharepinOpts{}
	}
	opts = opts.JoinPinls()

	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).
		Columns(Pinls{}.columns()...).
		From(s.table())
	qb = s.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*model.Sharepin, 0)
	for rows.Next() {
		var (
			msp model.Sharepin
			mp  model.Pinl
		)
		scanCols := append(s.scanColumns(&msp), Pinls{}.scanColumns(&mp)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		msp.Pinl = &mp
		list = append(list, &msp)
	}
	return list, nil
}

func (s Sharepins) bindOpts(b squirrel.SelectBuilder, opts *SharepinOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.ShareIDs) > 0 {
		b = b.Where(squirrel.Eq{s.table() + ".share_id": opts.ShareIDs})
	}

	if len(opts.PinlIDs) > 0 {
		b = b.Where(squirrel.Eq{s.table() + ".pinl_id": opts.PinlIDs})
	}

	if opts.Status.Valid {
		if sv, ok := opts.Status.Value().(model.Status); ok {
			b = b.Where(s.table()+".status = ?", sv)
		}
	}

	if opts.PinlQuery != "" {
		opts = opts.JoinPinls()
		pinls := Pinls{}.table()
		b = b.Where(squirrel.Or{
			squirrel.Expr(pinls+".title like ?", "%"+opts.PinlQuery+"%"),
			squirrel.Expr(pinls+".description like ?", "%"+opts.PinlQuery+"%"),
			squirrel.Expr(pinls+".url like ?", "%"+opts.PinlQuery+"%"),
		})
	}

	if len(opts.TagIDs) > 0 {
		sq := s.Builder().Select("1").
			From(Taggables{}.table()).
			Where("target_id = "+s.table()+".pinl_id").
			Where("target_name = ?", model.Pinl{}.MorphName()).
			Where(squirrel.Eq{"tag_id": opts.TagIDs}).
			GroupBy("target_id").
			Having("COUNT( DISTINCT tag_id ) >= ?", len(opts.TagIDs)).
			Prefix("EXISTS (").
			Suffix(")")
		b = b.Where(sq)
	}

	if opts.joinPinls {
		b = b.LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.pinl_id", Pinls{}.table(), s.table()))
	}

	return b
}

func (s Sharepins) columns() []string {
	return []string{
		s.table() + ".id",
		s.table() + ".share_id",
		s.table() + ".pinl_id",
		s.table() + ".status",
	}
}

func (s Sharepins) scanColumns(sharepin *model.Sharepin) []interface{} {
	return []interface{}{
		&sharepin.ID,
		&sharepin.ShareID,
		&sharepin.PinlID,
		&sharepin.Status,
	}
}

func (s Sharepins) scan(row database.RowScanner) (*model.Sharepin, error) {
	var sharepin model.Sharepin
	err := row.Scan(s.scanColumns(&sharepin)...)
	if err != nil {
		return nil, err
	}
	return &sharepin, nil
}

func (s *Sharepins) Create(ctx context.Context, sharepin *model.Sharepin) error {
	sharepin2 := *sharepin
	sharepin2.ID = newID()

	qb := s.RunnableBuilder(ctx).
		Insert(s.table()).
		Columns(
			"id",
			"share_id",
			"pinl_id",
			"status").
		Values(
			sharepin2.ID,
			sharepin2.ShareID,
			sharepin2.PinlID,
			sharepin2.Status)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharepin = sharepin2
	return nil
}

func (s *Sharepins) Update(ctx context.Context, sharepin *model.Sharepin) error {
	sharepin2 := *sharepin

	qb := s.RunnableBuilder(ctx).
		Update(s.table()).
		Set("share_id", sharepin2.ShareID).
		Set("pinl_id", sharepin2.PinlID).
		Set("status", sharepin2.Status).
		Where("id = ?", sharepin2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharepin = sharepin2
	return nil
}

func (s *Sharepins) Delete(ctx context.Context, id string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *SharepinOpts) JoinPinls() *SharepinOpts {
	o2 := *o
	o2.joinPinls = true
	return &o2
}
