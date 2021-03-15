package db

import "database/sql"

func GetNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func GetNullInt32(i uint32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(i),
		Valid: i != 0,
	}
}
