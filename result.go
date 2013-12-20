// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package odbc

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

type Result struct {
	rowCount     int64
	lastInsertId int64
	c            *Conn
}

func (r *Result) LastInsertId() (int64, error) {
	if r.lastInsertId == 0 {
		stmt, err := r.c.Prepare("SELECT @@Identity AS lastInsertId")
		if err != nil {
			return 0, err
		}
		var args []driver.Value
		rows, err := stmt.Query(args)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		dest := make([]driver.Value, 1)
		err = rows.Next(dest)
		if err != nil {
			return 0, err
		}
		switch dest[0].(type) {
		case int64:
			r.lastInsertId = dest[0].(int64)
		case int32:
			r.lastInsertId = int64(dest[0].(int32))
		case float64:
			r.lastInsertId = int64(dest[0].(float64))
		case float32:
			r.lastInsertId = int64(dest[0].(float32))
		case string:
			if dest[0].(string) == "NULL" {
				return 0, nil
			}
			r.lastInsertId, err = strconv.ParseInt(dest[0].(string), 10, 64)
			return r.lastInsertId, err
		default:
			return 0, errors.New("Unknow lastInsertId type")
		}
	}
	return r.lastInsertId, nil
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowCount, nil
}
