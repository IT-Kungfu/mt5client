package mt5client

type Config struct {
	MT5Host           string
	MT5Port           int
	MT5Login          string
	MT5Password       string
	MT5APIVersion     string
	MT5APIAgent       string
	MT5PingTimeout    int
	MT5RequestTimeout int
}
