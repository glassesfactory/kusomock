package kusomock

import (
  "encoding/json"
  "io"
  "net/http"

  toml "github.com/pelletier/go-toml"
  "github.com/zenazn/goji"
  "github.com/zenazn/goji/graceful"
  "github.com/zenazn/goji/web"
  mgo "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

// Init 初期化
func Init(filename string) {
  if filename == "" {
    filename = "./config.toml"
  }
  // Config ファイルを読む
  config, err := toml.LoadFile(filename)
  if err != nil {
    panic(err)
  }

  // db の設定を取る
  dbConfig := config.Get("database").(*toml.TomlTree)
  dbName := dbConfig.Get("db_name").(string)
  dbHost := dbConfig.Get("hostname").(string)

  mgoSession, err := mgo.Dial(dbHost)
  if err != nil {
    panic("mongo sikutta")
  }
  defer mgoSession.Close()

  mgoSession.SetMode(mgo.Monotonic, true)

  db = mgoSession.DB(dbName)

  // 静的ファイル
  static := web.New()
  publicPath := config.Get("general.public_path").(string)
  // safe join 的なのなかったっけ
  publicURL := "/" + config.Get("general.public_url").(string) + "/"
  static.Get(publicURL + "*", http.StripPrefix(publicURL, http.FileServer(http.Dir(publicPath))))

  http.Handle(publicURL, static)

  goji.Get("/api/:collection", SomeIndex)
  goji.Get("/api/:collection/:id", SomeShow)
  goji.Post("/api/:collection/", SomePost)
  goji.Put("/api/:collection/:id", SomePut)
  goji.Delete("/api/:collection/:id", SomeDelete)

  graceful.PostHook(func() {
  })
  goji.Serve()
}

// Start 開始
func Start() {
  goji.Serve()
}

// CreateResponse Create response
func CreateResponse(body string, w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json")
  io.WriteString(w, body)
}

// SomeIndex some index
func SomeIndex(c web.C, w http.ResponseWriter, r *http.Request) {
  colName := c.URLParams["collection"]
  col := db.C(colName)

  page := c.URL.Query().Get("page").(int)
  if page == 0 {
    page = 1
  }

  limit := c.URL.Query().GET("limit").(int)
  if limit == 0 {
    limit = 20
  }

  var results []struct{}
  // どうすっかな
  col.Find(nil).Limit(limit).All(&results)

  body, _ := json.Marshal(results)
  CreateResponse(string(body), w)
}

// SomeShow some collection detail
func SomeShow(c web.C, w http.ResponseWriter, r *http.Request) {

  colName := c.URLParams["collection"]
  id := c.URLParams["id"]
  col := db.C(colName)

  var result struct{}
  col.Find(bson.M{"id": id}).One(result)

  body, _ := json.Marshal(result)
  CreateResponse(string(body), w)
}

// SomePost Some collection insert
func SomePost(c web.C, w http.ResponseWriter, r *http.Request) {
  colName := c.URLParams["collection"]
  col := db.C(colName)

  // decode json
  decoder := json.NewDecoder(r.Body)
  var req struct{}
  err := decoder.Decode(&req)

  if err != nil {
    panic(err)
  }

  err = col.Insert(req)
  if err != nil {
    panic(err)
  }

  body, _ := json.Marshal(req)
  CreateResponse(string(body), w)
}

// SomePut update some collection detail
func SomePut(c web.C, w http.ResponseWriter, r *http.Request) {
  colName := c.URLParams["collection"]
  id := c.URLParams["id"]
  col := db.C(colName)

  var data struct{}
  err := col.Find(bson.M{"id": id}).One(data)
  if err != nil {
    panic(err)
  }

  decoder := json.NewDecoder(r.Body)
  var req struct{}
  err = decoder.Decode(&req)
  if err != nil {
    panic(err)
  }

  err = col.Update(bson.M{"id": id}, req)
  if err != nil {
    panic(err)
  }

  body, _ := json.Marshal(req)
  CreateResponse(string(body), w)
}

// SomeDelete Delete Some collection data
func SomeDelete(c web.C, w http.ResponseWriter, r *http.Request) {
  colName := c.URLParams["collection"]
  id := c.URLParams["id"]
  col := db.C(colName)

  err := col.Remove(bson.M{"id": id})
  if err != nil {
    panic(err)
  }

  res := struct {
    Id string `json:"id"`
  }{}
  res.Id = id
  body, _ := json.Marshal(res)
  CreateResponse(string(body), w)
}
