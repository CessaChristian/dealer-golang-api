package utils

import (

    "strings"
)

func NormalizeFuelType(fuel string) string {
    switch strings.ToLower(fuel) {
    case "gas":
        return "Bensin"
    case "diesel":
        return "Diesel"
    case "electric":
        return "Listrik"
    case "hybrid":
        return "Hybrid"
    default:
        return strings.Title(fuel) // fallback
    }
}

func NormalizeTransmission(t string) string {
    switch strings.ToLower(t) {
    case "a":
        return "Automatic"
    case "m":
        return "Manual"
    case "auto":
        return "Automatic"
    case "manual":
        return "Manual"
    default:
        return strings.Title(t)
    }
}

