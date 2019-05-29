package catsql

import (
  "net"

  "github.com/go-pg/pg"
  "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"

  "github.com/Liquid-Labs/env/go"
)

var db *pg.DB

// RawConnect sets up the application connection without initializing
// application specific prepared statements. It's typically used by
// non-application tools. Most application users will want to use 'Connect'
func RawConnect() *pg.DB {
  options := pg.Options{
    User:     env.MustGet("CLOUDSQL_USER"),
    Password: env.MustGet("CLOUDSQL_PASSWORD"), // NOTE: password may NOT be empty
    Database: env.MustGet("CLOUDSQL_DB"),
  }
  if env.IsTest() {
    options.Dialer = func(network, addr string) (net.Conn, error) {
      return proxy.Dial(env.MustGet(`CLOUDSQL_CONNECTION_NAME`))
    }
  } else {
    options.Addr = env.MustGet("CLOUDSQL_CONNECTION_NAME")
  }

  db = pg.Connect(&options)
  return db
}

// Connect connects the application to the database and prepares 
func Connect() *pg.DB {
  RawConnect()
  Warm()
  return db
}
