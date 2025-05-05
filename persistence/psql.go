package persistence

import (
	"context"
	"fmt"
	"kagari/setting"

	pgx "github.com/jackc/pgx/v5"
)

func url() string {
	return fmt.Sprintf("postgresql://postgres.idrvigdwubfgcyjconyy:%s@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres", setting.PsqlPassword())
}
func WithPsqlConnection(ctx context.Context, f func(con *pgx.Conn)) error {
	con, err := pgx.Connect(ctx, url())
	if err != nil {
		return err
	}
	defer con.Close(ctx)
	f(con)
	return nil
}
