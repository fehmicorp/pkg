package os

import "os"

func GetEnv(key, defValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defValue
}

func GetEnvV2(key string, defValue ...string) (bool, string) {
	if val := os.Getenv(key); val != "" {
		return true, val
	}
	if len(defValue) > 0 {
		return true, defValue[0]
	}
	return false, ""
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func ClearEnv() {
	os.Clearenv()
}

func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}
