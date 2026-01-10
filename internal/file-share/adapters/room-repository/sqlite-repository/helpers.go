package sqliterepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func readMaxSQLVars(ctx context.Context, db *sql.DB, fallback int) int {
	rows, err := db.QueryContext(ctx, `PRAGMA compile_options;`)
	if err != nil {
		return fallback
	}
	defer rows.Close()

	const prefix = "MAX_VARIABLE_NUMBER="
	for rows.Next() {
		var opt string
		if err := rows.Scan(&opt); err != nil {
			continue
		}
		if val, ok := strings.CutPrefix(opt, prefix); ok {
			var n int
			_, _ = fmt.Sscanf(val, "%d", &n)
			if n > 0 {
				return n
			}
		}
	}
	return fallback
}

func makePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(2*n - 1)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('?')
	}
	return b.String()
}

func chunkStrings(in []string, size int) [][]string {
	if size <= 0 || len(in) == 0 {
		return nil
	}
	out := make([][]string, 0, (len(in)+size-1)/size)
	for i := 0; i < len(in); i += size {
		end := i + size
		if end > len(in) {
			end = len(in)
		}
		out = append(out, in[i:end])
	}
	return out
}

func argsFromStrings(in []string) []any {
	out := make([]any, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}
