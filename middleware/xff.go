package cMiddleware

import (
	"github.com/labstack/echo/v4"
	"log"
	"net"
	"strings"
)

func RealIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		clientIP := getRealIP(c)
		log.Printf("Client IP %s", clientIP)

		c.Set("client_ip", clientIP)
		return next(c)
	}
}

// Learning purposes.
// Echo utiliza la funcion ExtractIPFromRealIPHeader() y ExtractIPFromXFFHeader() para
// extraer la ip real de nuestro cliente.

// Cuando utilizamos una proxy por encima de nuestra web, esta proxy nos reenviara las solicitudes
// de nuestros clientes, lo que quiere decir que la IP que vendra, sera otra. Pero los headers
// de la request incluyen la IP original, para encontrarla debemos parsear la primera IP que
// aparece en el header que nuestro proxy provider usa.
func getRealIP(c echo.Context) string {
	// Headers que esperamos:
	headers := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range headers {
		ip := c.Request().Header.Get(header)
		if ip != "" {
			// Spliteamos las multiples IPS y nos quedamos con l aprimera
			ips := strings.Split(ip, ",")
			clientIP := strings.TrimSpace(ips[0])

			// Utilizamos ParseIP para asegurarnos de que la IP tiene correcto formato
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}

		}
	}

	// Si no encontramos Proxy Headers, utilizaremos la remote adress, que es
	// la IP que nos ha enviado al request.
	ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	if net.ParseIP(ip) != nil {
		return ip
	}

	return "Unknown"

}
