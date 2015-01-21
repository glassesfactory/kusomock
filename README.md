kusomock
====

思いつきで適当に30分ぐらいで書いた雑な RESTful API モックツクール

思いつきで書いただけなのでおいおい全部の機能がちゃんと動くようにします。

* mongodb 
* goji `github.com/zenazn/goji`
* mgo `gopkg.in/mgo.v2`
* toml `github.com/pelletier/go-toml`

URL中に任意のコレクション名をつけてリクエストを飛ばすことで
あれば使うしなければ作るしな
mongodb 先生に仕事をしてもらえる。

`GET /api/:collection_name`

```json
{
	"somecollection": [...]
}
```

という感じ。
とりあえずダミーでいいから CRUD な API ほしい時に使ってもらえれば。

使い方
----

```golang

import "github.com/glassesfactory/kusomock"

func main() {
	kusomock.Start()
}
```

要望とかあれば issue に突っ込んでください。