package report_renderers

import "unsafe"

func toInterfaceSlice(strs []string) []interface{} {
	out := make([]interface{}, len(strs))
	for i, s := range strs {
		out[i] = s
	}
	return out
}

func escapeMarkdown(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	buf := make([]byte, 0, n+n/8)
	for i := 0; i < n; i++ {
		ch := s[i]
		switch ch {
		case '!', '*', '_', '-', '~', '`', '[', ']', '(', ')', '>', '|', '.', '+':
			// skip if reflection already exists
			if i != 0 && s[i-1] == '\\' {
				continue
			}
			buf = append(buf, '\\', ch)
		default:
			buf = append(buf, ch)
		}
	}
	return *(*string)(unsafe.Pointer(&buf))
}

func safeCutTextAlongRune(s string, lim int) string {
	if len(s) <= lim {
		return s
	}
	i := lim
	for i > 0 && (s[i]&0xC0) == 0x80 {
		i--
	}
	return s[:i]
}
