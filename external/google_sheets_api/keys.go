package google_sheets_api

const (
	MajorDimensionRows    string = "ROWS"    // use in sheets.ValueRange. array is interpreted as strings (each subarray is a string)
	MajorDimensionColumns string = "COLUMNS" // use in sheets.ValueRange. the array is interpreted as columns

	ValueInputOptionRaw         string = "RAW"          // writes data as is (text, numbers)
	ValueInputOptionUserEntered string = "USER_ENTERED" // interprets data as user input (formulas, numbers, text)

	InsertDataOptionInsertRows string = "INSERT_ROWS" // data is inserted as new rows (default)
	InsertDataOptionOverwrite  string = "OVERWRITE"   // existing data can be replaced

	MergeTypeAll     string = "MERGE_ALL"     // merges all cells in a range into one
	MergeTypeColumns string = "MERGE_COLUMNS" // merges cells by column (each column is merged separately)
	MergeTypeRows    string = "MERGE_ROWS"    // merges cells by row (each row is merged separately).

	UpdateCellFieldUserEnteredValue string = "userEnteredValue"
	UserEnteredValueString          string = "stringValue"
	UserEnteredValueNumber          string = "numberValue"
	UserEnteredValueBool            string = "boolValue"
	UserEnteredValueFormula         string = "formulaValue"
	UserEnteredValueError           string = "errorValue"

	UpdateCellFieldUserEnteredFormat      string = "userEnteredFormat"
	UserEnteredFormatBackColor            string = "backgroundColor"
	UserEnteredFormatTextFormat           string = "textFormat"
	UserEnteredFormatTextFormatBold       string = "bold"
	UserEnteredFormatTextFormatItalic     string = "italic"
	UserEnteredFormatTextFormatUnderline  string = "underline"
	UserEnteredFormatTextFormatFontFamily string = "fontFamily"
	UserEnteredFormatTextFormatFontSize   string = "fontSize"
	UserEnteredFormatHorizontalAlignment  string = "horizontalAlignment"
	UserEnteredFormatVerticalAlignment    string = "verticalAlignment"
	UserEnteredFormatWrapStrategy         string = "wrapStrategy"
	UserEnteredFormatBorders              string = "borders"
	UserEnteredFormatPadding              string = "padding"

	UpdateCellFieldNote           string = "note"
	UpdateCellFieldHyperlink      string = "hyperlink"
	UpdateCellFieldDataValidation string = "dataValidation"
	UpdateCellFieldPivotTable     string = "pivotTable"
	UpdateCellFieldEffectiveValue string = "effectiveValue"

	BorderStyleNone        string = "NONE"
	BorderStyleDotted      string = "DOTTED"
	BorderStyleDashed      string = "DASHED"
	BorderStyleSolid       string = "SOLID"
	BorderStyleSolidMedium string = "SOLID_MEDIUM"
	BorderStyleSolidThick  string = "SOLID_THICK"
	BorderStyleSolidDouble string = "DOUBLE"

	TextHorizontalAlignmentLeft   string = "LEFT"
	TextHorizontalAlignmentCenter string = "CENTER"
	TextHorizontalAlignmentRight  string = "RIGHT"

	TextVerticalAlignmentTop    string = "TOP"
	TextVerticalAlignmentMiddle string = "MIDDLE"
	TextVerticalAlignmentBottom string = "BOTTOM"

	TextWrapStrategyOverflowCell string = "OVERFLOW_CELL" // Text extends beyond cell
	TextWrapStrategyWrap         string = "WRAP"          // Wrap text on a new line
	TextWrapStrategyClip         string = "CLIP"          // Trim text that doesn't fit in a cell

	TextDirectionLeftRight string = "LEFT_TO_RIGHT" // Text from left to right
	TextDirectionRightLeft string = "RIGHT_TO_LEFT" // Text from right to left
)
