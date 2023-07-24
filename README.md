# MyCookBook

MyCookBookはアプリケーションはユーザーの料理アイデアを記録し、レシピの提案、食材の効率的な管理を可能とするアプリケーションです。  
MyCookBookは、あなたの料理体験を向上させ、創造性を引き出すためのパートナーとなることでしょう。是非MyCookBookをご活用ください！  
このリポジトリはサーバーサイドのものとなっています。  
アプリケーションリンク: [MyCookBook](http://mycooookbook.com:3000)  
フロントエンドのリポジトリ: [frontend](https://github.com/Gatch-you/menu_proposer_frontend)

**使用技術**  
![Static Badge](https://img.shields.io/badge/golang-1.20.1-blue) ![Static Badge](https://img.shields.io/badge/MySQL-14.14-green)
 ![Static Badge](https://img.shields.io/badge/AWS-EC2,RDS-yellow)

## 機能の紹介
このアプリケーションでは、以下の機能を実装しています。

### 食材情報の登録・変更・削除
ユーザーが入手した食材に対し名前、数量、賞味期限などの登録を行い、その食材情報の変更、削除等が可能です。

### レシピの登録・変更・削除
ユーザーの作りたい料理(自身の自慢のレシピや作ってみたい料理等)の登録を行い、そのレシピ情報の変更、削除が可能です。  
また、登録した食材情報と結びつけて使用数の設定や作成時の在庫管理を行うことができます。

### レシピの提案
ユーザーの食材登録時に入力した賞味期限に対して5日以内の食材の情報を取得し、その食材にて作成できるレシピを提案することが可能です。

