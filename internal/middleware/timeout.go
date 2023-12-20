package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/xfang04/go-web/internal/log"
)

type TimeoutMiddleWare struct {
	Next http.Handler
}

func (t *TimeoutMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.Next == nil {
		t.Next = http.DefaultServeMux
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel() // Call the cancel function to avoid a context leak
	r = r.WithContext(ctx)
	ch := make(chan struct{})
	go func() {
		t.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
		log.Info.Println("[Request handle time out!]")
	}
	ctx.Done()
}
