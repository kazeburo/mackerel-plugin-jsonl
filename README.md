
# mackerel-plugin-jsonl

JSON Lines 形式のログからメトリックを抽出し、Mackerelに投稿するためのカスタムプラグインです。

## インストール方法

mkrを使うと楽です

```
mkr plugin --upgrade kazeburo/mackerel-plugin-jsonl
```

## 基本的な使い方

```
./mackerel-plugin-jsonl \
	--key-name count \
	--json-key foo.bar \
	--aggregator count \
	--log-file /path/to/your.log \
	--prefix custom.jsonl
```

複数のキーや集計方法を指定する場合は、`--key-name`, `--json-key`, `--aggregator` を同じ数だけ指定してください。

## mackerel-agent.conf 設定例

```ini
[plugin.metrics.jsonl]
command = "/path/to/mackerel-plugin-jsonl --key-name count --json-key foo.bar --aggregator count --log-file /var/log/app.log --prefix custom.jsonl"
```


## aggregatorの種類と説明

`--aggregator` には以下の種類を指定できます。

- `count` : 対象となるキーが存在する行をカウントします
	- 例: `json.total.count        80000.000000    1760281381`
- `group_by` : 指定したキーの値ごとに件数を集計します。
	- 例: `json.status.2xx 24769.733333    1760281381`
- `group_by_with_percentage` : group_byの各値ごとの割合（%）も出力します。
	- 例: `json.status_percentage.5xx      23.139667       1760281381`
- `percentile` : 指定した数値キーのパーセンタイル（mean, p90, p95, p99）を出力します。
	- 例: `json.latency.p95        0.030000        1760281381`

複数のaggregatorを同時に指定することで、複数のメトリクスを一度に出力できます。

---

## json-keyの変換関数

`--json-key` にはパイプ（`|`）区切りで変換関数を指定できます。

例: `foo.bar|tolower|replace('a','b')|trimspace`

利用可能な関数:

- `tolower` : 小文字化
- `toupper` : 大文字化
- `trimspace` : 前後の空白を除去
- `replace('pattern','repl')` : 正規表現で置換

例:
```
--json-key "user.name|trimspace|tolower"
--json-key "message|replace('error','warn')"
```

## JSONによるアクセスログのメトリクス作成例

```
% /path/to/mackerel-plugin-jsonl --prefix json --log-file json.log -k total.count -j time -a count -k status -j 'status|replace("^(.).+$","${1}xx")' -a group_by_with_percentage -k latency -j reqtime -a percentile
json.total.count        80000.000000    1760281381
json.status.4xx 24481.600000    1760281381
json.status.2xx 24769.733333    1760281381
json.status.3xx 12236.933333    1760281381
json.status.5xx 18511.733333    1760281381
json.status_percentage.5xx      23.139667       1760281381
json.status_percentage.4xx      30.602000       1760281381
json.status_percentage.2xx      30.962167       1760281381
json.status_percentage.3xx      15.296167       1760281381
json.latency.mean       0.030000        1760281381
json.latency.p90        0.030000        1760281381
json.latency.p95        0.030000        1760281381
json.latency.p99        0.030000        1760281381
```


## ライセンス

MIT License
