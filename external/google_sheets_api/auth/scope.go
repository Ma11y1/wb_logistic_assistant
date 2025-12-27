package auth

import "google.golang.org/api/sheets/v4"

type Scope string

const (
	SheetsAllScope      Scope = sheets.SpreadsheetsScope
	SheetsReadonlyScope Scope = sheets.SpreadsheetsReadonlyScope
)
