package utils

/*

e.g.

var things FlagSliceString

flag.Var(&things, "thing", "a single thing")

*/

type FlagSliceString []string

func (f *FlagSliceString) String() string {
	return "" // TODO: is this important?
}

func (f *FlagSliceString) Set(value string) error {
	*f = append(*f, value)

	return nil
}
