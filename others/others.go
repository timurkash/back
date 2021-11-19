package others

import (
	"strings"
	"time"
)

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func DoEvery(d time.Duration, f func(time.Time)) {
	go func() {
		for x := range time.Tick(d) {
			f(x)
		}
	}()
}

func GetWeek(offset int) (*time.Time, *time.Time) {
	now := time.Now()
	weekBegin := time.Time{}
	switch now.Weekday() {
	case time.Monday:
		weekBegin = now
	case time.Tuesday:
		weekBegin = now.Add(-24 * time.Hour)
	case time.Wednesday:
		weekBegin = now.Add(-48 * time.Hour)
	case time.Thursday:
		weekBegin = now.Add(-72 * time.Hour)
	case time.Friday:
		weekBegin = now.Add(-96 * time.Hour)
	case time.Saturday:
		weekBegin = now.Add(-120 * time.Hour)
	case time.Sunday:
		weekBegin = now.Add(-144 * time.Hour)
	}
	weekBegin = weekBegin.Truncate(24 * time.Hour).Add(time.Duration(offset*168) * time.Hour)
	weekEnd := weekBegin.Add(168 * time.Hour)
	return &weekBegin, &weekEnd
}

func Shortize(str string, limit int) string {
	if len(str) < limit {
		return str
	}
	str = FirstN(str, limit-3) + "..."
	return str
}

func FirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func GetNameExt(filename string) (string, string) {
	parts := strings.Split(filename, ".")
	l := len(parts)
	if l == 1 {
		return filename, ""
	}
	return strings.Join(parts[:l-1], "."), parts[l-1]
}

func Min(value int, values ...int) int {
	for _, v := range values {
		if v < value {
			value = v
		}
	}
	return value
}
