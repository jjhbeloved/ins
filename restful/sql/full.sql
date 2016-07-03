CREATE TABLE IF NOT EXISTS device (
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
CREATE TABLE IF NOT EXISTS d2i (
  d2i_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host   TEXT,
  ip     TEXT
);

CREATE TABLE IF NOT EXISTS ip (
  ip TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user (
  user_id   INTEGER PRIMARY KEY AUTOINCREMENT,
  host      TEXT,
  uid       INTEGER,
  name      TEXT    NOT NULL,
  gid       INTEGER NOT NULL,
  password  TEXT,
  last_date DATETIME
);

/**
group 2 group
 */
CREATE TABLE IF NOT EXISTS g2g (
  g2g_id     INTEGER PRIMARY KEY AUTOINCREMENT,
  parent_gid INTEGER, /* 父组 */
  sub_gi     INTEGER  /* 子组 */
);

CREATE TABLE IF NOT EXISTS sgroup (
  sgroup_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host      TEXT,
  gid       INTEGER,
  name      TEXT NOT NULL
);

/**
device 2 app
 */
CREATE TABLE IF NOT EXISTS d2u2a (
  d2a_id INTEGER PRIMARY KEY AUTOINCREMENT,
  host   TEXT,
  uid    INTEGER,
  app_id TEXT
);

CREATE TABLE IF NOT EXISTS app (
  app_id    TEXT,
  app_name  TEXT,
  version   TEXT,
  app_home  TEXT,
  ips       TEXT,
  ports     TEXT,
  last_date DATETIME,
  PRIMARY KEY (app_id)
);




