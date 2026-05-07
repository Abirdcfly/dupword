package raw_string_multiline

func example() string {
	const q = ` // want `Duplicate words \(the\) found`
SELECT the the column
FROM table
WHERE id = 1;
`
	return q
}
