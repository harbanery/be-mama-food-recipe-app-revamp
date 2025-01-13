package helper

type EnvLoad struct {
	APP_NAME string

	WEB_PREFORK string
	WEB_PORT    string

	CORS_ALLOW_ORIGINS  string
	CORS_ALLOW_METHODS  string
	CORS_ALLOW_HEADERS  string
	CORS_EXPOSE_HEADERS string

	DB_HOST          string
	DB_PORT          string
	DB_USER          string
	DB_PASS          string
	DB_NAME          string
	DB_SSLMODE       string
	DB_POOL_IDLE     string
	DB_POOL_MAX      string
	DB_POOL_LIFETIME string

	JWT_SECRET_KEY string
}
