package render

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
)

func Respond(w http.ResponseWriter, r *http.Request, v interface{}) {
	// Present the object.
	if presenter, ok := r.Context().Value(presenterCtxKey).(Presenter); ok {
		r, v = presenter.Present(r, v)
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Chan:
		switch getResponseContentType(r) {
		case ContentTypeEventStream:
			channelEventStream(w, r, v)
			return
		default:
			v = channelIntoSlice(w, r, v)
		}
	}

	// Format data based on Content-Type.
	switch getResponseContentType(r) {
	case ContentTypeJSON:
		JSON(w, r, v)
	case ContentTypeXML:
		XML(w, r, v)
	default:
		JSON(w, r, v)
	}
}

func PlainText(w http.ResponseWriter, r *http.Request, v string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write([]byte(v))
}

func Data(w http.ResponseWriter, r *http.Request, v []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(v)
}

func HTML(w http.ResponseWriter, r *http.Request, v string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write([]byte(v))
}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func XML(w http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := xml.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	// Try to find <?xml header in first 100 bytes (just in case there're some XML comments).
	findHeaderUntil := len(b)
	if findHeaderUntil > 100 {
		findHeaderUntil = 100
	}
	if bytes.Index(b[:findHeaderUntil], []byte("<?xml")) == -1 {
		// No header found. Print it out first.
		w.Write([]byte(xml.Header))
	}

	w.Write(b)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func channelEventStream(w http.ResponseWriter, r *http.Request, v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Chan {
		panic(fmt.Sprintf("render.EventStream() expects channel, not %v", reflect.TypeOf(v).Kind()))
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(200)

	ctx := r.Context()
	for {
		switch chosen, recv, ok := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(v)},
		}); chosen {
		case 0: // equivalent to: case <-ctx.Done()
			w.Write([]byte("event: error\ndata: {\"error\":\"Server Timeout\"}\n\n"))
			return

		default: // equivalent to: case v, ok := <-stream
			if !ok {
				w.Write([]byte("event: EOF\n\n"))
				return
			}
			v := recv.Interface()

			// Present each channel item.
			if presenter, ok := r.Context().Value(presenterCtxKey).(Presenter); ok {
				r, v = presenter.Present(r, v)
			}

			// TODO: Can't use json Encoder - it panics on bufio.Flush(). Why?!
			// enc := json.NewEncoder(w)
			// enc.Encode(v.Interface())
			bytes, err := json.Marshal(v)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("event: error\ndata: {\"error\":\"%v\"}\n\n", err)))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				continue
			}
			w.Write([]byte(fmt.Sprintf("event: data\ndata: %s\n\n", bytes)))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

// channelIntoSlice buffers channel data into a slice.
func channelIntoSlice(w http.ResponseWriter, r *http.Request, from interface{}) interface{} {
	ctx := r.Context()

	var to []interface{}
	for {
		switch chosen, recv, ok := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(from)},
		}); chosen {
		case 0: // equivalent to: case <-ctx.Done()
			http.Error(w, "Server Timeout", 504)
			return nil

		default: // equivalent to: case v, ok := <-stream
			if !ok {
				return to
			}
			v := recv.Interface()

			// Present each channel item.
			if presenter, ok := r.Context().Value(presenterCtxKey).(Presenter); ok {
				r, v = presenter.Present(r, v)
			}

			to = append(to, v)
		}
	}
	panic("unreachable")
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi render context value " + k.name
}
