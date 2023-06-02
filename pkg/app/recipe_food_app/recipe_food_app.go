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
	RecipeId           int     `json:"recipe_id"`
	RecipeName         string  `json:"recipe_name"`
	RecipeDescription  string  `json:"recipe_description"`
	FoodId             int     `json:"food_id"`
	FoodName           string  `json:"food_name"`
	UseAmount          float64 `json:"use_amount"`
	FoodUnit           string  `json:"food_unit"`
	RecipeMakingMethod string  `json:"recipe_making_method"`
}

// 使用する食材の名前と量の情報を保持しているレシピの一覧表示。
// curl http://localhost:8080/backend/recipe_food
func FetchRecipesWithFood(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT r.name, r.description, f.name, rf.use_amount, f.unit, r.making_method FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id")
	if err != nil {
		log.Fatal(err.Error())
	}

	rfArgs := make([]Recipe_food, 0)
	for rows.Next() {
		var recipe_food Recipe_food
		err = rows.Scan(&recipe_food.RecipeName, &recipe_food.RecipeDescription, &recipe_food.FoodName, &recipe_food.UseAmount, &recipe_food.FoodUnit, &recipe_food.RecipeMakingMethod)
		if err != nil {
			log.Fatal(err.Error())
		}
		rfArgs = append(rfArgs, recipe_food)
	}

	v, err := json.Marshal(rfArgs)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Write([]byte("Show the recipe with foods.\n"))
	w.Write([]byte(v))
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

	fmt.Printf("You've finished cooking! I've finished updating list you use ingredient. Nice Cooking!")

}

// レシピにて使用する食材の量を変更する処理
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

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Hey, you alter the amount of ingredients in recipe. OK, I accept")
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

			recipeRows, err := db.Query("SELECT r.name, f.name, rf.use_amount FROM recipe_food rf JOIN foods f ON rf.food_id = f.id JOIN recipes r ON rf.recipe_id = r.id WHERE f.expiration_date >= DATE(NOW()) AND f.expiration_date <= DATE_ADD(DATE(NOW()), INTERVAL 5 DAY)")
			if err != nil {
				log.Fatal(err.Error())
			}

			recipeWithExpirationFoodsArgs := make([]Recipe_food, 0)
			for recipeRows.Next() {
				var recipe_food Recipe_food
				err = recipeRows.Scan(&recipe_food.RecipeName, &recipe_food.FoodName, &recipe_food.UseAmount)
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
