package recipe_food_app

import (
	"backend/pkg/app/foods_app"
	"backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Recipe_food struct {
	ID                 int     `json:"id"`
	RecipeId           int     `json:"recipe_id"`
	RecipeName         string  `json:"recipe_name"`
	RecipeDescription  string  `json:"recipe_description"`
	FoodId             int     `json:"food_id"`
	FoodName           string  `json:"food_name"`
	UseAmount          float64 `json:"use_amount"`
	FoodUnit           string  `json:"food_unit"`
	RecipeMakingMethod string  `json:"recipe_making_method"`
}

type RecipeFoodArray struct {
	FoodID    int     `json:"food_id"`
	RecipeID  int     `json:"recipe_id"`
	UseAmount float64 `json:"use_amount"`
}

type FoodsWithExpiration struct {
	ID             int       `json:"id"`
	FoodId         int       `json:"food_id"`
	FoodName       string    `json:"food_name"`
	FoodQuantity   float64   `json:"food_quantity"`
	FoodUnit       string    `json:"food_unit"`
	ExpirationDate time.Time `json:"expiration_date"`
	RecipeId       int       `json:"recipe_id"`
	RecipeName     string    `json:"recipe_name"`
	UseAmount      float64   `json:"use_amount"`
}

type RecipeID struct {
	RecipeID int `json:"recipe_id"`
}

// 使用する食材の名前と量の情報を保持しているレシピの一覧表示。
// curl http://localhost:8080/backend/recipe_food
func FetchRecipesWithFood(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT rf.id, r.id, r.name, r.description, f.id, f.name, rf.use_amount, f.unit, r.making_method FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id")
	if err != nil {
		log.Fatal(err.Error())
	}

	rfArgs := make([]Recipe_food, 0)
	for rows.Next() {
		var recipe_food Recipe_food
		err = rows.Scan(&recipe_food.ID, &recipe_food.RecipeId, &recipe_food.RecipeName, &recipe_food.RecipeDescription, &recipe_food.FoodId, &recipe_food.FoodName, &recipe_food.UseAmount, &recipe_food.FoodUnit, &recipe_food.RecipeMakingMethod)
		if err != nil {
			log.Fatal(err.Error())
		}
		rfArgs = append(rfArgs, recipe_food)
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Write([]byte("Show the recipe with foods.\n"))
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rfArgs)
}

// curl -X GET  http://localhost:8080/backend/recipes/{id}
func FetchRecipeDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/backend/recipes/")
	db := db.Connect()
	defer db.Close()

	// var recipeID RecipeID
	// err := json.NewDecoder(r.Body).Decode(&recipeID)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	stmt, err := db.Prepare("SELECT rf.id, r.id, r.name, r.description, f.id, f.name, rf.use_amount, f.unit, r.making_method FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id WHERE r.id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	var recipe_foods []Recipe_food
	for rows.Next() {
		var recipe_food Recipe_food
		err = rows.Scan(&recipe_food.ID, &recipe_food.RecipeId, &recipe_food.RecipeName, &recipe_food.RecipeDescription, &recipe_food.FoodId, &recipe_food.FoodName, &recipe_food.UseAmount, &recipe_food.FoodUnit, &recipe_food.RecipeMakingMethod)
		if err != nil {
			log.Fatal(err.Error())
		}
		recipe_foods = append(recipe_foods, recipe_food)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}
	fmt.Print("\nYou use API")
	fmt.Print(recipe_foods)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(recipe_foods)
}

// 料理の作成後、使用した食材分foodsから引く機能
// curl -X PUT http://localhost:8080/backend/recipe_food/updata_quantity/1
func UpdateFoodStorage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/backend/recipe_food/update_food_storage/")

	db := db.Connect()
	defer db.Close()

	update, err := db.Prepare("UPDATE foods AS f INNER JOIN recipe_food AS rf ON f.id = rf.food_id SET f.quantity = f.quantity - rf.use_amount WHERE rf.recipe_id = ? AND f.quantity >= rf.use_amount")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = update.Exec(id)
	if err != nil {
		log.Fatal(err.Error())
		log.Println("Ingredients is out of stock!")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}

	fmt.Printf("You've finished cooking! I've finished updating list you use ingredient. Nice Cooking!")
	w.WriteHeader(http.StatusOK)
}

// レシピにて使用する食材の追加処理を行う関数. 将来的には一括で食材を登録することができるようにするために実装
// curl -X POST -d '[{"recipe_id": 1, "food_id": 3, "use_amount": 2},{"recipe_id": 5, "food_id": 21, "use_amount": 100}]' http://localhost:8080/backend/recipe_food/insert_use_food_array
func InsertUseFoodArray(w http.ResponseWriter, r *http.Request) {
	// recipeFoodArrayへとリクエストボディに受け取ったjsonデータを配列として受け取る
	var recipeFoodArray []RecipeFoodArray
	err := json.NewDecoder(r.Body).Decode(&recipeFoodArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	// dbへの接続処理
	db := db.Connect()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO recipe_food (food_id, recipe_id, use_amount) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, RecipeFoodArray := range recipeFoodArray {
		_, err := insert.Exec(RecipeFoodArray.FoodID, RecipeFoodArray.RecipeID, RecipeFoodArray.UseAmount)
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)
}

// レシピにて使用する食材の登録処理 実装済み
func InsertUseFood(w http.ResponseWriter, r *http.Request) {
	// dbへの接続処理
	db := db.Connect()
	defer db.Close()

	// recipeFoodArrayへとリクエストボディに受け取ったjsonデータを配列として受け取る
	var recipeFoodArray RecipeFoodArray
	err := json.NewDecoder(r.Body).Decode(&recipeFoodArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	inst, err := db.Prepare("INSERT INTO recipe_food (recipe_id, food_id, use_amount) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = inst.Exec(recipeFoodArray.RecipeID, recipeFoodArray.FoodID, recipeFoodArray.UseAmount)
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
	json.NewEncoder(w).Encode(recipeFoodArray)
}

// レシピにて使用する食材の量を変更する処理 実装済み
// curl -X PUT -H "Content-Type: application/json" -d '{"recipe_id": 2, "food_id": 3, "use_amount": 3}' http://localhost:8080/backend/recipe_food/update_using_food_quantity
func UpdateUsingFoodQuantity(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	var recipe_food Recipe_food
	err := json.NewDecoder(r.Body).Decode(&recipe_food)
	if err != nil {
		log.Fatal(err.Error())
	}

	update, err := db.Prepare("UPDATE recipe_food SET use_amount = ? WHERE recipe_id = ? AND food_id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = update.Exec(recipe_food.UseAmount, recipe_food.RecipeId, recipe_food.FoodId)
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

	fmt.Println("Hey, you alter the amount of ingredients in recipe. OK, I accept")
}

// レシピにて使用する食材の削除 実装済み
// curl -X DELETE -d '{"recipe_id": 1, "food_id": 14}' http://localhost:8080/backend/delete_using_food
func DeleteUsingFood(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	var recipeFoodArray RecipeFoodArray
	err := json.NewDecoder(r.Body).Decode(&recipeFoodArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	delete, err := db.Prepare("DELETE FROM recipe_food WHERE recipe_id = ? AND food_id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = delete.Exec(recipeFoodArray.RecipeID, recipeFoodArray.FoodID)
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
	json.NewEncoder(w).Encode(recipeFoodArray)
}

// 定時になったら、賞味期限が指定した日時以内の食品の一覧を表示し、
// その食品を使って作ることができるレシピと、その食材の使用量を出力する関数
func FetchExpirationFood(w http.ResponseWriter, r *http.Request) []Recipe_food {
	db := db.Connect()
	defer db.Close()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		now := time.Now().In(loc)

		fmt.Println(now)

		if now == time.Date(now.Year(), now.Month(), now.Day(), 15, now.Minute(), now.Second(), now.Nanosecond(), loc) {

			foodRows, err := db.Query("SELECT name, quantity, unit, expiration_date FROM foods WHERE expiration_date >= DATE(NOW()) AND expiration_date <= DATE_ADD(DATE(NOW()), INTERVAL 5 DAY)")
			if err != nil {
				log.Fatal(err.Error())
			}

			expirationFoodArgs := make([]foods_app.Food, 0)
			for foodRows.Next() {
				var food foods_app.Food
				err = foodRows.Scan(&food.Name, &food.Quantity, &food.Unit, &food.ExpirationDate)
				if err != nil {
					log.Fatal(err.Error())
				}
				expirationFoodArgs = append(expirationFoodArgs, food)
			}

			recipeRows, err := db.Query("SELECT rf.id, r.name, f.name, rf.use_amount FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id WHERE f.expiration_date >= DATE(NOW()) AND f.expiration_date <= DATE_ADD(DATE(NOW()), INTERVAL 5 DAY)")
			if err != nil {
				log.Fatal(err.Error())
			}

			recipeWithExpirationFoodsArgs := make([]Recipe_food, 0)
			for recipeRows.Next() {
				var recipe_food Recipe_food
				err = recipeRows.Scan(&recipe_food.ID, &recipe_food.RecipeName, &recipe_food.FoodName, &recipe_food.UseAmount)
				if err != nil {
					log.Fatal(err.Error())
				}
				recipeWithExpirationFoodsArgs = append(recipeWithExpirationFoodsArgs, recipe_food)
			}

			fmt.Println("Hello, Foods!")
			fmt.Println(expirationFoodArgs)
			fmt.Println(recipeWithExpirationFoodsArgs)

			// jsonへと変換
			// v, err := json.Marshal(expirationFoodArgs)
			// if err != nil {
			// 	log.Fatal(err.Error())
			// }
			// fmt.Println(v)

			// goroutineで並列処理を実装するとmainから渡された関数の引数nilとwがぶつかってエラーが出たので、今はw.Writeはつかわない。
			// w.Write([]byte("Show the Foods which expiration date having been closed in 3 days\n"))
			// w.Write([]byte(v))
		}
		time.Sleep(time.Hour * 5)
	}
}

func ShowFoodsWithExpiration(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT rf.id, f.id, f.name, f.quantity, f.unit, f.expiration_date, r.id, r.name, rf.use_amount, f.unit FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id WHERE f.expiration_date >= DATE(NOW()) AND f.expiration_date <= DATE_ADD(DATE(NOW()), INTERVAL 5 DAY)")
	if err != nil {
		log.Fatal(err.Error())
	}

	foodArgs := make([]FoodsWithExpiration, 0)
	for rows.Next() {
		var food_expiration FoodsWithExpiration
		err = rows.Scan(&food_expiration.ID, &food_expiration.FoodId, &food_expiration.FoodName, &food_expiration.FoodQuantity, &food_expiration.FoodUnit, &food_expiration.ExpirationDate, &food_expiration.RecipeId, &food_expiration.RecipeName, &food_expiration.UseAmount, &food_expiration.FoodUnit)
		if err != nil {
			fmt.Print("missing query")
			log.Fatal(err.Error())
		}
		foodArgs = append(foodArgs, food_expiration)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(foodArgs)
}
