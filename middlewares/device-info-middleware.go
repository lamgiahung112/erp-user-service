package middlewares

import (
	"context"
	"net/http"
)

func (md *MiddlewareConfig) GetDeviceInfo(next http.Handler) http.Handler {
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

		userAgent := r.Header.Get("user-agent")
		deviceInfo, err := md.Utils.DeviceInfo.GetDevice(userAgent, iplocation)
		if err != nil {
			md.serveError(err, w, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), DeviceInfoMiddlewareResultKey, deviceInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
