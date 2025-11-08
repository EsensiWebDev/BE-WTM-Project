package constant

import "time"

var AsiaJakarta = mustLoadLocation("Asia/Jakarta")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic("invalid timezone: " + name)
	}
	return loc
}
