package main

import (
	"fmt"
	"strings"
	"time"
	"math"
)

func formatDuration(d time.Duration) string {
	d = time.Duration(math.Ceil(d.Seconds())) * time.Second
	seconds := int(d.Seconds())

	if seconds < 60 {
		return fmt.Sprintf("%d секунд", seconds)
	}

	minutes := seconds / 60
	seconds = seconds % 60

	if minutes < 60 {
		return fmt.Sprintf("%d минут %d секунд", minutes, seconds)
	}

	hours := minutes / 60
	minutes = minutes % 60

	if hours < 24 {
		return fmt.Sprintf("%d часа %d минут %d секунд", hours, minutes, seconds)
	}

	days := hours / 24
	hours = hours % 24

	return fmt.Sprintf("%d дней %d часа %d минут %d секунд", days, hours, minutes, seconds)
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}