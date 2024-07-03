// Package ipfilter - Пакет для функциональности фильтрации запросов по IP адресу
package ipfilter

import (
	"net"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"go.uber.org/zap"
)

// NewIPFilterMiddleware - создать middleware фильтра запросов по ip
func NewIPFilterMiddleware(trustedIPNet *net.IPNet) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Применяем фильтр только если доверенная подсеть задана
			if trustedIPNet != nil {

				// проверяем IP адрес из заголовка X-Real-IP, который должен всегда устанавливаться на прокси (если используется)
				realIPString := r.Header.Get("X-Real-IP")
				var clientIP net.IP
				if realIPString != "" {
					ip := net.ParseIP(realIPString)
					if ip == nil {
						logger.Log.Debug("error parsing ip address while filtering", zap.String("x-real-ip", realIPString))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					clientIP = ip
				} else {
					// Если IP адрес не задан в заголовке X-Real-IP - тогда берем его из RemoteAddr request'а
					addr := r.RemoteAddr
					// метод возвращает адрес в формате host:port
					// нужна только подстрока host
					ipStr, _, err := net.SplitHostPort(addr)
					if err != nil {
						logger.Log.Debug("error splitting host addr", zap.String("addr", addr))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					// парсим ip
					ip := net.ParseIP(ipStr)
					if ip == nil {
						logger.Log.Debug("error parsing ip address while filtering", zap.String("request-ip", ipStr))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					clientIP = ip
				}

				// Доступ запрещен, если адрес не принадлежит доверенной подсети
				if !trustedIPNet.Contains(clientIP) {
					logger.Log.Debug("request filtered: ip address is not in trusted subnet", zap.String("ip", clientIP.String()), zap.String("trusted-subnet", trustedIPNet.String()))
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			// передаём управление хендлеру
			h.ServeHTTP(w, r)
		})
	}
}
