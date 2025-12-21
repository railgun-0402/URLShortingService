# パーティション

## パーティションを試す目的
- DBを分割することで大量リクエスト・データ等に対しスケーラビリティ向上
- 今回はRDB(Postgres)を例にするが、NoSQLやインメモリDBでも利用される
- ボトルネックの分散を実際に試す

**→とにかくデータが均等に分散されることが重要！**

## スケーラビリティを向上すると？
- DBサーバーには物理的制約あり
  - CPU処理能力
  - メモリ容量
  - ディスクI/O
  - NW帯域幅
- 上記により、リクエストやデータ量が増加するとパフォーマンス低下につながる


## テーブル設計（PostgreSQL）

- short_urlsのテーブル設計を思い出してみる↓

```sql
CREATE TABLE short_urls (
    id           VARCHAR(16) PRIMARY KEY,   -- 短縮ID (Base62)
    original_url TEXT        NOT NULL,      -- 元URL
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
    expired_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_short_urls_created_at ON short_urls (created_at);
```
- 例えばURLを「期間で集計」する場合は・・・
- 2025/12、2026/1など該当期間のパーティションのみを読むことになる

### メリット

- 該当期間のパーティションのみを読むので、読み取りI/Oが減りやすい
- 保持期限(expired_at)の運用が楽になる
  - データに対しDELETEを実施する等必要なくなる
  - 論理削除で放っておくでもいいが、、データ量が多いPJは読み取り負荷に影響あり
- パーティションごとにインデックスが分かれる
  - 1個あたりのインデックスが小さくなる等

#### イメージ

```
click_events（親：論理テーブル）
  ├─ click_events_2025_12（子：2025/12の分）
  ├─ click_events_2026_01（子：2026/01の分）
  └─ ...
```

### デメリット

- 運用が増える(今回はURLの期間ごとにできるパーティション)
- 親テーブルに対するALTERが影響する場合があり、デグレの可能性
- パーティションを増やし過ぎると上記影響が顕著になる
- パーティションキーのデータ分割が難しい
  - 日付で分散することで、データ分布が不均等になる可能性も・・・

## 手順

```bash
docker exec -it urlshort_pg psql -U urlshort -d urlshort
```

### 親テーブル作成

- ここにデータは入らない

```sql
CREATE TABLE IF NOT EXISTS click_events (
  occurred_at TIMESTAMPTZ NOT NULL,
  id BIGSERIAL NOT NULL,
  short_url_id BIGINT NOT NULL,
  referrer TEXT NULL,
  user_agent TEXT NULL,
  PRIMARY KEY (occurred_at, id)
) PARTITION BY RANGE (occurred_at);

```

### 月次パーティションを作る（例：2025/12, 2026/01）

```sql
CREATE TABLE IF NOT EXISTS click_events_2025_12
  PARTITION OF interviewcat.click_events
  FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE IF NOT EXISTS click_events_2026_01
  PARTITION OF interviewcat.click_events
  FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
```

- インデックス

```sql
-- 親に作る（環境によっては各パーティションにも作られる/必要なら各子にも作成）
CREATE INDEX IF NOT EXISTS idx_click_events_url_time
  ON click_events (short_url_id, occurred_at);
```

### partition pruning

- 実際にどの子テーブルを読んだか
- どれくらい時間がかかったか
を“見える化”してみよう

- 期間条件を指定すると、パーティション効果あり

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM click_events
WHERE short_url_id = 42
  AND occurred_at >= '2025-12-01'
  AND occurred_at <  '2026-01-01';
```

- 期間条件無視

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM click_events
WHERE short_url_id = 42;
```

Append の下で 参照されてるパーティションが少ない（or 1個だけ）になってたら pruning 成功！

### 保持期間（retention）：古い月を捨てる

例：2025/12 を捨てるなら、パーティションを落とすだけ。

```sql
DROP TABLE click_events_2025_12;
```

- DELETE FROM click_events WHERE occurred_at < ... は重い
- パーティション単位削除は速いしVACUUM地獄を避けやすい

### Goからもデータ投入で体感してみよう！

- まずは 50万〜200万件 くらい click_events を入れる
- COPY でもいいが、最初はバルクINSERTでも可

#### データ投入の方針（擬似）：

- occurred_at を 2025-12 と 2026-01 に散らして入れる
- 計測クエリを2パターン流す（期間あり/なし）

## リファクタ案

現在は50万件のデータを500 * 1000回INSERTに試しにしているが・・・
batchSizeを2000にするとINSERTは250回で済む

```go
// とりあえず1回のINSERTに詰める行数を設定
const batchSize = 500
```

### batchSizeが多：メリット

- 1度実行する際のINSERT回数減(NW/COMMIT/実行回数の負担が減る)
- 単純に速くなる

### batchSizeが多：デメリット

- 1度のクエリが巨大になるので実行が重くなる
- 1度の実行によるメモリ使用量が増える
- postgresでmax_stack_depthやwork_memに引っかかる可能性
- RDSに接続する際NWで不安定になる可能性

### ベストな値を探る

当然計測が必要になる。

1. まず固定条件を決める
- n（例：200,000）
- 同じRDS/同じ時間帯

2. batchSizeだけ変えて計測

- 500 / 1000 / 2000 / 5000 あたり

3. rows/sec と 失敗率 と CPU/メモリ を見る

ここは詳しくないが、AIによると↓なので要検証

```bash
さらに速さ追求するなら、Postgresは結局 COPY が王者です。バッチINSERTは「手軽さ優先」
```
