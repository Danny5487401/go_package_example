package _4_squirrel

import (
	"log"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func TestSquirrel(t *testing.T) {
	// 查询
	users := sq.Select("*").From("users").Join("emails USING (email_id)")

	active := users.Where(sq.Eq{"deleted_at": nil})

	sql, args, err := active.ToSql()

	if err != nil {
		log.Fatal(err)
	}
	if sql != `SELECT * FROM users JOIN emails USING (email_id) WHERE deleted_at IS NULL` {
		t.Errorf("got %s", sql)
	}

	t.Logf("sql %v : ,args:%v ", sql, args)

	// 插入
	now := time.Now()
	sql, args, err = sq.Insert("post").Columns(
		"created_at", "updated_at", "app", "user_id", "tag", "content", "comment_count",
	).Values(now, now, "douyin", 123, "test", "hello,body", 2).ToSql()
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("sql %v : ,args:%v ", sql, args)

}
