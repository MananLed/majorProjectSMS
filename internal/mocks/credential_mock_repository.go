package mocks 

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// )

// type MockDB struct {
// 	Queries map[string]any // query string → return value
// 	ExecErr map[string]error
// }

// // Implements QueryRow
// func (m *MockDB) QueryRow(query string, args ...any) *sql.Row {
// 	key := fmt.Sprintf(query, args...)
// 	val, ok := m.Queries[key]
// 	if !ok {
// 		// return a Row that will scan to false
// 		return &sql.Row{}
// 	}
	
// 	// we can’t easily fake *sql.Row, but for testing Delete we only check `Scan(&exists)`
// 	// so we can wrap using sql.Row manually OR instead return a struct implementing scanning
// 	return &mockRow{val: val}
// }

// // Implements Exec
// func (m *MockDB) Exec(query string, args ...any) (sql.Result, error) {
// 	key := fmt.Sprintf(query, args...)
// 	if err, ok := m.ExecErr[key]; ok {
// 		return nil, err
// 	}
// 	return mockResult{}, nil
// }

// // ---- helpers ----
// type mockRow struct {
// 	val any
// }

// func (r *mockRow) Scan(dest ...any) error {
// 	if len(dest) != 1 {
// 		return errors.New("mockRow supports only single column scan")
// 	}
// 	switch d := dest[0].(type) {
// 	case *bool:
// 		*d = r.val.(bool)
// 	default:
// 		return errors.New("unsupported scan type")
// 	}
// 	return nil
// }

// type mockResult struct{}

// func (mockResult) LastInsertId() (int64, error) { return 0, nil }
// func (mockResult) RowsAffected() (int64, error) { return 1, nil }