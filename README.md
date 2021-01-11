# 入力されたjsonと指定したファイルのjsonのkeyが全て一致しているかチェックするツール
サーバーサイドの人が自分の作ったエンドポイントのjsonが正しいかをチェックするのに使ってください

# 使用例
## 例1
jsonをそのままコマンドに渡す
`go run main.go '{"key1": "hoge", "key2": "huga"}' -file ./correctJson/correct.json`

## 例2
curlの結果をコマンドに渡す
`curl -X POST 'http://localhost:3030/path1' | go run main.go -file ./correctJson/correct.json`