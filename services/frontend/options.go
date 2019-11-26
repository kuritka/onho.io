package frontend

type Options struct {
	ClientSecret string
	ClientID     string
	CookieStoreKey string
	Port int
	Idp Idp
	QueueConnectionString string
}