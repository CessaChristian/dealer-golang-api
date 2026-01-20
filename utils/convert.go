package utils

import "strconv"

func ToInt(s string) int {
    v, _ := strconv.Atoi(s)
    return v
}