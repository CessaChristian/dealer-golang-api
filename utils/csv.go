package utils

import "strconv"

func SafeCSV(s string) string {
    if s == "" {
        return "-"
    }
    return s
}

func SafeCSVInt(v int) string {
    return strconv.Itoa(v)
}

func SafeCSVFloat(v float64) string {
    return strconv.FormatFloat(v, 'f', 2, 64)
}
