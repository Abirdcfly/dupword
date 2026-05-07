package raw_string_dup

func example() string {
	const s = `hello the the world` // want `Duplicate words \(the\) found`
	return s
}
