package foods_app

import (
	db "backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		foodArgs = append(foodArgs, food)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Print(foodArgs)

	json.NewEncoder(w).Encode(foodArgs)
}

// 食品の検索に使うが、他の検索アルゴリズムに置き換わる可能性大
// func SearchFoods(w http.ResponseWriter, r *http.Request) {

// 	elasticURL := "http://localhost:9200"
// 	elasticUser := "root"
// 	elasticPass := "hoge"

// 	esClient, err := elastic.NewClient(
// 		elastic.SetURL(elasticURL),
// 		elastic.SetBasicAuth(elasticUser, elasticPass),
// 	)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	db := db.Connect()
// 	defer db.Close()

// 	rows, err := db.Query("SELECT * FROM foods")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	for rows.Next() {
// 		var food Food
// 		err = rows.Scan(&food.ID, &food.Name, &food.Quantity, &food.Unit, &food.ExpirationDate, &food.Type)
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}

// 		//Elasticsearchへデータを保存
// 		_, err = esClient.Index().
// 			Index("foods").
// 			BodyJson(food).
// 			Do(context.Background())
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}
// 	}

// 	query := r.URL.Query().Get("name")
// 	// 間違ってるかも

// 	searchResult, err := esClient.Search().
// 		Index("foods").
// 		Query(elastic.NewMatchQuery("name", query)).
// 		From(0).Size(10).
// 		Do(context.Background())
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	//elasticServer上からヒットした検索結果を[]Foodのスライスに格納押したのちにfoodsへとappendしている。
// 	foods := make([]Food, 0)
// 	for _, hit := range searchResult.Hits.Hits {
// 		var food Food
// 		err := json.Unmarshal(hit.Source, &food)
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}
// 		foods = append(foods, food)
// 	}
// 	v, err := json.Marshal(foods)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Methods", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")

// 	w.Write(v)
// }

// 新しい食品の項目追加
// ↓実行コマンド
// curl -X POST -d '{"name": "豚肉", "quantity": 250, "unit": "g", "expiration_date": "2023-06-10T00:00:00Z", "type": "精肉"}' http://localhost:8080/backend/insert_food
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)

}

// 食品の数量、個数の変化をこのコードにて処理する。0の量もこのデータにて扱う

// curl -X PUT -H "Content-Type: application/json" -d '{"name": "キャベツ", "quantity": 0.3, "unit": " 個", "expiration_date": "2023-05-02T00:00:00Z", "type": "野菜"}' http://localhost:8080/backend/update_food
func UpdateFoods(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Fatal(err.Error())
	}

	update, err := db.Prepare("UPDATE foods SET id = ?, name = ?, quantity = ?, unit = ?, expiration_date = ?, type = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = update.Exec(food.ID, food.Name, food.Quantity, food.Unit, food.ExpirationDate, food.Type, food.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)
}

// 食品のデータベースのフィールドそのものを削除する。再びその食品を使うには再度InsertFoodsを叩かなければならなくなる。
// curl -X DELETE -H "Content-Type: application.json" -d '{"food_id": 1}' http://localhost:8080/backend/delete_food
func DeleteFoods(w http.ResponseWriter, r *http.Request) {
	// id := strings.TrimPrefix(r.URL.Path, "/backend/delete_food/")

	db := db.Connect()
	defer db.Close()

	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Fatal(err.Error())
	}

	delete, err := db.Prepare("DELETE FROM foods WHERE id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = delete.Exec(food.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)
}
