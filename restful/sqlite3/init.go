package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"asiainfo.com/ins/restful/models"
	"time"
)

var create_sql = `
CREATE TABLE if not exists device (
  host    TEXT,
  memory  INTEGER,
  cpus    INTEGER,
  storage TEXT,
  os      TEXT,
  PRIMARY KEY (host)
);

/**
device to ip
 */
CREATE TABLE if not exists d2i (
  d2i_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host   TEXT,
  ip     TEXT
);

CREATE TABLE if not exists ip (
  ip TEXT PRIMARY KEY
);

/**
device 2 user
 */
CREATE TABLE if not exists d2u (
  d2u_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host   TEXT,
  uid    INTEGER
);

CREATE TABLE if not exists user (
  uid       INTEGER,
  name      TEXT    NOT NULL,
  gid       INTEGER NOT NULL,
  password  TEXT,
  last_date DATETIME,
  PRIMARY KEY (uid)
);

/**
group 2 group
 */
CREATE TABLE if not exists g2g (
  g2g_id     INTEGER PRIMARY KEY AUTOINCREMENT,
  parent_gid INTEGER, /* 父组 */
  sub_gi     INTEGER  /* 子组 */
);

CREATE TABLE if not exists sgroup (
  gid  INTEGER,
  name TEXT NOT NULL,
  PRIMARY KEY (gid)
);

/**
device 2 app
 */
CREATE TABLE if not exists d2u2a (
  d2a_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host   TEXT,
  uid	 INTEGER,
  app_id TEXT
);

CREATE TABLE if not exists app (
  app_id    TEXT,
  app_name  TEXT,
  version   TEXT,
  app_home  TEXT,
  ips       TEXT,
  ports     TEXT,
  last_date DATETIME,
  PRIMARY KEY (app_id)
);
`

func Open(name string) (*sql.DB, error) {
	return sql.Open("sqlite3", name)
}

func Close(db *sql.DB) error {
	return db.Close()
}

func CreateBasicTables(db *sql.DB) error {
	_, err := db.Exec(create_sql)
	return err
}

func Insert(db *sql.DB, sql string) (*sql.Stmt, error) {
	return db.Prepare(sql)
}

func main() {
	group := &models.Group{
		Name: "root",
		GID: 0,
	}
	user := &models.User{
		Name: "root",
		UID: 0,
		Group:group,
		Password:"123",
		Last_date: time.Now(),
	}
	device := &models.Device{
		Host:"xiaoxiao",
		IP:[]string{"192.168.2.105"},
		Memory: 8,
		CPUs: 2,
		Storage: "512G",
		OSType:"OS X Yosemite",
	}
	fmt.Println(len(device.Users))
	fmt.Println(len(user.APPs))

	db, err := Open("/tmp/c.c")
	if err != nil {
		fmt.Println(err)
	}
	CreateBasicTables(db)
	stmt, _ := Insert(db, models.INSERT_DEVICE)
	res, err := device.Insert_Device(stmt)
	if err!= nil {
		fmt.Println(err)
	}
	fmt.Println(res.LastInsertId())
	if err := Close(db); err != nil {
		fmt.Println(err)
	}
}