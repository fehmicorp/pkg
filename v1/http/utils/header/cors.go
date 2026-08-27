package header

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(corsConfig *CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if corsConfig == nil {
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		method := c.Request.Method

		// 1. Normalize Dynamic Origins Configuration
		allowedOrigins := parseDynamicConfig(corsConfig.Origins)

		// If explicit false/empty config and there is an origin, treat it strictly as same-origin
		if len(allowedOrigins) == 1 && allowedOrigins[0] == "false" {
			allowedOrigins = []string{}
		}

		// 2. Normalize Dynamic Headers Configuration
		allowedHeaders := parseDynamicConfig(corsConfig.Headers)
		headersAllDenied := len(allowedHeaders) == 1 && allowedHeaders[0] == "false"

		// 3. Extract requested headers (relevant for OPTIONS preflight checks)
		reqHeadersStr := c.GetHeader("Access-Control-Request-Headers")
		var requestHeaders []string
		if reqHeadersStr != "" {
			for _, h := range strings.Split(reqHeadersStr, ",") {
				requestHeaders = append(requestHeaders, strings.TrimSpace(h))
			}
		}

		// Validate the HTTP method against allowed CORS methods early
		if !containsIgnoreCase(corsConfig.Methods, "*") && !containsIgnoreCase(corsConfig.Methods, method) && method != "OPTIONS" {
			c.AbortWithStatusJSON(405, gin.H{"error": "CORS request forbidden: method unauthorized"})
			return
		}

		// 4. Run Validation Functions for Cross-Origin requests
		if origin != "" {
			// Validate Origin
			if !CORSValidate(origin, allowedOrigins) {
				c.AbortWithStatusJSON(403, gin.H{"error": "CORS request forbidden: origin unauthorized"})
				return
			}

			// Validate Headers
			if headersAllDenied && len(requestHeaders) > 0 {
				c.AbortWithStatusJSON(403, gin.H{"error": "CORS request forbidden: all custom headers are denied"})
				return
			}
			if !CORSHeadersValidate(requestHeaders, allowedHeaders) {
				c.AbortWithStatusJSON(403, gin.H{"error": "CORS request forbidden: headers unauthorized"})
				return
			}

			// 5. Append required Headers to the response
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)

			if len(corsConfig.Methods) > 0 {
				c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.Methods, ", "))
			}

			if len(allowedHeaders) > 0 && !headersAllDenied {
				if containsIgnoreCase(allowedHeaders, "*") && reqHeadersStr != "" {
					c.Writer.Header().Set("Access-Control-Allow-Headers", reqHeadersStr)
				} else {
					c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
				}
			}
		}

		// 6. Instantly short-circuit Preflight Options Requests
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func parseDynamicConfig(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return []string{}
	}

	var boolVal bool
	if err := json.Unmarshal(raw, &boolVal); err == nil {
		if !boolVal {
			return []string{"false"}
		}
		return []string{"*"}
	}

	var strVal string
	if err := json.Unmarshal(raw, &strVal); err == nil {
		return []string{strVal}
	}

	var sliceVal []string
	if err := json.Unmarshal(raw, &sliceVal); err == nil {
		return sliceVal
	}

	return []string{}
}

func CORSHeadersValidate(requestHeaders []string, allowedHeaders []string) bool {
	if len(allowedHeaders) == 0 {
		return true
	}
	if containsIgnoreCase(allowedHeaders, "*") {
		return true
	}
	for _, header := range requestHeaders {
		if !containsIgnoreCase(allowedHeaders, header) {
			return false
		}
	}
	return true
}

func CORSValidate(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return origin == ""
	}
	return containsIgnoreCase(allowedOrigins, "*") || containsIgnoreCase(allowedOrigins, origin)
}

func containsIgnoreCase(list []string, value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	for _, item := range list {
		if strings.TrimSpace(strings.ToLower(item)) == value {
			return true
		}
	}
	return false
}
