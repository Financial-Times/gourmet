package gmlog

import "net/http"

func WithField(key string, val interface{}) Field {
	return func(e *Event) {
		e.put(key, val)
	}
}

func WithRequest(r *http.Request) Field {
	return func(e *Event) {
		route := r.URL.Path
		query := r.URL.RawQuery
		e.put("proto", r.Proto)
		e.put("method", r.Method)
		e.put("route", route)

		if len(query) > 0 {
			e.put("query", query)
		}
	}
}

func WithStatusCode(code int) Field {
	return WithField("status_code", code)
}

func WithTransactionsID(tid string) Field {
	return WithField("transaction_id", tid)
}
