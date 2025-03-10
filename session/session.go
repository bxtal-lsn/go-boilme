package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
}

func (b *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// how long should sessions last?
	minutes, err := strconv.Atoi(b.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	// should cookies persist?
	if strings.ToLower(b.CookiePersist) == "true" {
		persist = true
	}

	// must cookies be secure?
	if strings.ToLower(b.CookieSecure) == "true" {
		secure = true
	}

	// create session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = b.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = b.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// which session store?
	switch strings.ToLower(b.SessionType) {
	case "redis":
		session.Store = redisstore.New(b.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(b.DBPool)
	case "postgres", "postgresql":
		session.Store = postgresstore.New(b.DBPool)
	default:
		// cookie
	}

	return session
}
