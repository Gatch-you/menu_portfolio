package foods_app

import (
	"encoding/json"
	"fmt"
	"log"
	db "menu_proposer/pkg/db"
	"net/http"
	"strings"
	"time"
)

// 食品の構造体json形式のデータ変換等も行う
type Food struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	ExpirationDate time.Time `json:"expiration_date"`
	Type           string    `json:"type"`
}

// 食品の一覧表示。frontにて表示件数を絞る必要が出てくるかも　→ フロントにて考慮のはず
func FetchFoods(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM foods")
	if err != nil {
		log.Fatal(err.Error())
	}

	foodArgs := make([]Food, 0)
	for rows.Next() {
		var food Food
		err = rows.Scan(&food.ID, &food.Name, &food.Quantity, &food.Unit, &food.ExpirationDate, &food.Type)
		if err != nil {
			log.Fatal(err.Error())
		}
		// z := food.ExpirationDate
		foodArgs = append(foodArgs, food)
		// fmt.Println(z)
	}

	v, err := json.Marshal(foodArgs)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Write([]byte("Show the Foods\n"))
	w.Write([]byte(v))
}

// 食品の検索に使うが、他の検索アルゴリズムに置き換わる可能性大
func FetchFoodByKey(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()
}

// 新しい食品の項目追加
// ↓実行コマンド
// curl -X POST -H "Content-Type: application/json" -d
// '{"id": value, "name": "value", "quantity": value, "unit": "value", "expiration_date": "XXXX-YY-ZZT00:00:00Z", "type": "value"}'
// http://localhost:8080/api/user/create

// curl -X POST -H "Content-Type: application/json" -d '{"id": 2, "name": "キャベツ", "quantity": 0.5, "unit": "個", "expiration_date": "2023-04-21T00:00:00Z", "type": "野菜"}' http://localhost:8080/menu_proposer/insert_food
func InsertFoods(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Fatal(err.Error())
	}

	expirationDate, err := time.Parse("2006-01-02", food.ExpirationDate.Format("2006-01-02"))
	if err != nil {
		log.Fatal(food.ExpirationDate)
	}
	food.ExpirationDate = expirationDate

	inst, err := db.Prepare("INSERT INTO foods (id, name, quantity, unit, expiration_date, type) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = inst.Exec(food.ID, food.Name, food.Quantity, food.Unit, food.ExpirationDate, food.Type)
	if err != nil {
		log.Fatal(err.Error())
	}

	bytes, err := json.Marshal(food)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Write([]byte("Insert New Food\n"))
	w.Write(bytes)

}

// 食品の数量、個数の変化をこのコードにて処理する。0の量もこのデータにて扱う

// curl -X PUT -H "Content-Type: application/json" -d
// '{"name": "キャベツ", "quantity": 0.3, "unit": " 個", "expiration_date": "2023-04-21T00:00:00Z", "type": "野菜"}'

// http://localhost:8080/menu_proposer/update_food/(id)
func UpdateFoods(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/menu_proposer/update_food/")

	db := db.Connect()
	defer db.Close()

	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Fatal(err.Error())
	}

	updt, err := db.Prepare("UPDATE foods SET name = ?, quantity = ?, unit = ?, expiration_date = ?, type = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = updt.Exec(food.Name, food.Quantity, food.Unit, food.ExpirationDate, food.Type, id)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Fprintf(w, "%s has been updated and quantity is altered.", food.Name)
}

// 食品のデータベースのフィールドそのものを削除する。再びその食品を使うには再度InsertFoodsを叩かなければならなくなる。
// curl -X DELETE localhost:8080/menu_proposer/delete_food/(id)
func DeleteFoods(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/menu_proposer/delete_food/")

	db := db.Connect()
	defer db.Close()

	delt, err := db.Prepare("DELETE FROM foods WHERE id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = delt.Exec(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Food which you select has been deleted")

}
