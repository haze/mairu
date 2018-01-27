package str

// S -> string -> nillable string
func S(a string) *string {
	tmp := a
	return &tmp
}
