## remo-local-api-handler

Nature Remo のローカル API を経由して, シグナルを飛ばす為の API です。[Cloud API](http://swagger.nature.global/) と違い, ローカル環境のみで実行できます。

Local API ではシグナルを記憶しておく方法がないため, この API を経由することで JSON ファイルに記録させておくことが出来ます。

## config/config.json

最初の設定として, `config/config.json` ファイルを作成する必要があります。

```json5
{
  // 接続詞たい Nature Remo のポート
  "remo_ip": "192.168.1.2",
  // サーバーを動かすポート
  "port": 8000
}
```

## Signal を取得

シグナルは, Nature Remo に向けて送信したあと, `/signal/` に GET すると JSON で取得できます。

取得した JSON を以下のフォーマットで `config/signals.json` に保存してください。

```json5
{
  "signals": [
    {
      // 任意の名前をつける.
      "name": "light_off",
      // 取得した JSON は "content" の中.
      "content": {"freq":37,"data":[3439,1792,388,482,388,482,388,1356,387,1353,391,480,390,1356,385,483,387,486,377,491,387,1358,387,479,389,481,389,1359,386,483,387,1352,391,481,389,1355,388,482,391,479,387,1353,391,481,388,482,388,483,388,481,389,1355,387,1355,388,484,386,1359,384,483,388,1354,388,484,387,482,388,486,384,1353,389,482,393,477,388,482,388,1357,386,481,387,486,387,65535,0,9543,3442,1791,388,484,387,483,388,1358,384,1354,382,493,379,1362,388,486,384,485,385,489,381,1359,386,482,388,482,392,1354,385,487,383,1361,382,489,381,1357,388,483,387,484,387,1356,386,484,387,484,387,486,384,484,386,1354,389,1357,386,485,388,1357,384,486,385,1354,389,486,384,485,385,486,386,1354,392,478,389,484,387,481,389,1357,386,483,387,484,388],"format":"us"}
    }
  ]
}
```

## Signal を送信

* POST `/signal/`

```json5
{
  // config/signals.json 内の名前を指定.
  "name": "light_off",
  // *または*, 直接シグナルを指定.
  "content": {"freq":37,"data":[3439,1792,388,482,388,482,388,1356,387,1353,391,480,390,1356,385,483,387,486,377,491,387,1358,387,479,389,481,389,1359,386,483,387,1352,391,481,389,1355,388,482,391,479,387,1353,391,481,388,482,388,483,388,481,389,1355,387,1355,388,484,386,1359,384,483,388,1354,388,484,387,482,388,486,384,1353,389,482,393,477,388,482,388,1357,386,481,387,486,387,65535,0,9543,3442,1791,388,484,387,483,388,1358,384,1354,382,493,379,1362,388,486,384,485,385,489,381,1359,386,482,388,482,392,1354,385,487,383,1361,382,489,381,1357,388,483,387,484,387,1356,386,484,387,484,387,486,384,484,386,1354,389,1357,386,485,388,1357,384,486,385,1354,389,486,384,485,385,486,386,1354,392,478,389,484,387,481,389,1357,386,483,387,484,388],"format":"us"}
}
```
## TODO

* [ ] Nature Remo 自動検出
* [ ] 認識したシグナルを自動保存
* [ ] フロント実装 (優先度低)
