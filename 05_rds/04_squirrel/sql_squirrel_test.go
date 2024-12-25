package _4_squirrel

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
)

func TestSquirrel(t *testing.T) {
	users := sq.Select("*").From("users").Join("emails USING (email_id)")

	active := users.Where(sq.Eq{"deleted_at": nil})

	sql, args, err := active.ToSql()

	if err == nil {
		t.Logf("sql %v : ,args:%v ", sql, args)
	}

}
