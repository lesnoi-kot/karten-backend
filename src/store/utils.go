package store

import "database/sql"

func NoRowsAffected(result sql.Result) bool {
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return true
	}

	return false
}
