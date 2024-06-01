package middlewares

import (
	"context"
	"log"
	"net/http"
)

func (md *MiddlewareConfig) GetIpLocation(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		iplocation, err := md.Utils.IpLocation.GetLocationDatafromIP(IPAddress)

		if err != nil {
			md.serveError(err, w, http.StatusInternalServerError)
			return
		}

		log.Println(iplocation)
		ctx := context.WithValue(r.Context(), IpLocationMiddlewareResultKey, iplocation)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
