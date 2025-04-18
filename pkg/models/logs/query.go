package logs

import (
	"strings"
	"time"
)

func QueryBefore(end time.Time) Filter {
	return func(meta *Meta) FilterStatus {
		if meta.CreatedAt.Before(end) {
			return OK
		} else {
			return STOP
		}
	}
}

func QueryAfter(start time.Time) Filter {
	return func(meta *Meta) FilterStatus {
		if meta.CreatedAt.After(start) {
			return OK
		} else {
			return SKIP
		}
	}
}

func QueryGroup(group string) Filter {
	if group == "" {
		return func(meta *Meta) FilterStatus {
			return OK
		}
	}
	if strings.HasSuffix(group, "/*") {
		group = strings.TrimSuffix(group, "*")
		return func(meta *Meta) FilterStatus {
			if strings.HasPrefix(meta.Group, group) {
				return OK
			} else {
				return SKIP
			}
		}
	} else {
		return func(meta *Meta) FilterStatus {
			if meta.Group == group {
				return OK
			} else {
				return SKIP
			}
		}
	}
}

func QueryAuthor(author string) Filter {
	return func(meta *Meta) FilterStatus {
		if strings.Contains(meta.Author.String(), author) {
			return OK
		} else {
			return SKIP
		}
	}
}

func QueryIpAddr(ip string) Filter {
	return func(meta *Meta) FilterStatus {
		if strings.Contains(meta.Ip, ip) {
			return OK
		} else {
			return SKIP
		}
	}
}

func QueryMsg(msg string) Filter {
	return func(meta *Meta) FilterStatus {
		if strings.Contains(meta.Msg, msg) {
			return OK
		} else {
			return SKIP
		}
	}
}
