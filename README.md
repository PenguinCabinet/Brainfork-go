# 🧠 Brainfork-go
[n4mlzさん](https://github.com/n4mlz)が[セキュリティキャンプ2025](https://www.ipa.go.jp/jinzai/security-camp/2025/camp/zenkoku/index.html)で提案した[Brainfork](https://github.com/n4mlz/Brainfork)の言語処理系です。   
Go言語で書かれています。

## 🏗ビルド

```
go build
```
## 🔨使い方
```
brainfork-go input-brainfork-sourcecode.bf
```

## 📃テスト
```
brainfork-go test/Correct-Cases/Parallel.bf
```

[Correct-Cases](test/Correct-Cases) ,[Incorrect-Cases](test/Incorrect-Cases)のテストケースを使ってテストが行われます。

## 追加された構文
以下は、[Brainfork](https://github.com/n4mlz/Brainfork)のソースコードを読みながら、私が推測したものです。一部、[n4mlzさん](https://github.com/n4mlz)にご教示いただきました。     
これら以外のものは[Brainfuxk](https://ja.wikipedia.org/wiki/Brainfuck)と同等です。     
?は私が理解できていない部分です。
|文字|動作|
|---|---|
|{|`}`までのプログラムを、スレッドにして開始する|
|\||そのスレッドにおける並行処理の区切り。`+++.`と`---.`を平行処理したい場合、`+++.\|---.`になる|
|}|スレッド開始の終端|
|(|現在のptrが指し示すメモリデータをロックする|
|)|対応する`)`が実行された時のptrが指し示すメモリデータをアンロックする|
|~|100ms待つ|
|^|現在のptrが指し示すメモリデータで待機し、処理を一時停止する|
|v|現在のptrが指し示すメモリデータで待機していたすべてのスレッドを再開する|
|;|これ以降、改行までコメントアウトする|

## 🎫LICENSE

[MIT](./LICENSE)

## ✍Author

[PenguinCabinet](https://github.com/PenguinCabinet)