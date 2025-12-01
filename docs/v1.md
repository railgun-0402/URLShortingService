# v1

とりあえずシンプルに作成して内容を掴む

## IDを一意にする方式
### 候補(とりあえず思いついたやつ書くだけ)
- ハッシュ
- Base62

### 大規模トラフィック
v1では考慮しないので、ランダム・衝突時のリトライのみの
カスAPIにする。

### 文字種
- 文字種：0-9a-zA-Z（= 62文字）
- 長さ：6〜8文字くらい（とりあえず 8 にしておけば十分）
- 生成方法：crypto/rand でランダムバイト → 62文字にマッピング
- 衝突対策：
  - 生成したIDが既に使われていたら再生成（DB or ストアを見て確認）

---

## 衝突を避ける方法
v1では・・・

- generateID(8) で ID を生成
- ストア（DB or メモリ）の id キーにデータがないか確認
- あれば「衝突なので再生成」
- なければ保存して採用

実運用なら DB の UNIQUE 制約 or 

PutItem(ConditionExpression) みたいなので守るのが定石な気がするが・・

まずは in-memory 実装 → 後から DB に差し替え

---

## 301 or 302 の選択

### 301（Moved Permanently）

- SEOを意識した本格サービスはこっちを使うことが多いらしい
- クライアントや検索エンジンがキャッシュしやすい（URLマッピングが変えづらくなる）

### 302（Found / 一時的リダイレクト）

- 「将来マッピング変えるかも」「クリック計測したい」など、
サーバー側の自由度を残したいときに使われがち

- 実際のURL短縮サービスも 302 or 307 を使ってるケース多め

→302を採用

別URL切り替えなどの余地を残しておく。

--- 

## 2️⃣ v1 のAPI仕様イメージ

### POST /shorten

- 入力：

```json
{
"url": "https://example.com/very/long/url/..."
}
```

- 出力：

```json
{
"id": "aBcD12eF",
"short_url": "https://short.local/aBcD12eF"
}
```


### GET /:id

- GET /aBcD12eF にアクセス

  → 元URLに 302 Redirect

# Let's try!

---

## POST /shorten

---

### request

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{
        "url": "https://example.com/foo/bar?param=123"
      }'
```

### result
```bash
{
  "id": "aBcD12eF",
  "short_url": "http://localhost:8080/aBcD12eF"
}
```

## GET /:id

---

### request

```bash
curl -i http://localhost:8080/aBcD12eF
```

### result
```bash
HTTP/1.1 302 Found
Location: https://example.com/foo/bar?param=123
```

## Not Exist ID

---

### request

```bash
curl -i http://localhost:8080/notfound123
```

### result
```bash
HTTP/1.1 404 Not Found
{"message":"short url not found"}
```
## empty URL

---

### request

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```

### result
```bash
HTTP/1.1 400 Bad Request
{"message":"url is required"}
```

## irregular URL

---

### request
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "hogehoge"}'
```

### result
```bash
HTTP/1.1 400 Bad Request
{"message":"invalid url format"}
```
