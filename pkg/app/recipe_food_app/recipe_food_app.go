package recipe_food_app

import (
	"encoding/json"
	"fmt"
	"log"
	"menu_proposer/pkg/db"
	"net/http"
	"strings"
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
// curl http://localhost:8080/menu_proposer/recipe_food
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

	w.Write([]byte("Show the recipe with ingridients.\n"))
	w.Write([]byte(v))
}

// 料理の作成後、使用した食材分foodsから引く機能
// curl -X PUT http://localhost:8080/menu_proposer/recipe_food/updata_quantity/1
func UpdateFoodStrage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/menu_proposer/recipe_food/updata_food_strage/")

	db := db.Connect()
	defer db.Close()

	updt, err := db.Prepare("UPDATE foods AS f INNER JOIN recipe_food AS rf ON f.id = rf.food_id SET f.quantity = f.quantity - rf.use_amount WHERE rf.recipe_id = ? AND f.quantity >= rf.use_amount")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = updt.Exec(id)
	if err != nil {
		log.Fatal(err.Error())
		log.Println("Ingridients is out of stock!")
	}

	fmt.Printf("You've finished cooking! I've finished updating list you use ingridient. Nice Cooking!")

}

// レシピにて使用する食材の量を変更する処理
// curl -X PUT -H "Content-Type: application/json" -d '{"recipe_id": 2, "food_id": 3, "use_amount": 3}' http://localhost:8080/menu_proposer/recipe_food/update_using_food_quantity
func UpdateUsingFoodQuantity(w http.ResponseWriter, r *http.Request) {
	// recipe_id := strings.TrimPrefix(r.URL.Path, "/menu_proposer/recipe_food/update_using_food_quantity/")

	db := db.Connect()
	defer db.Close()

	var recipe_food Recipe_food
	err := json.NewDecoder(r.Body).Decode(&recipe_food)
	if err != nil {
		log.Fatal(err.Error())
	}

	updt, err := db.Prepare("UPDATE recipe_food SET use_amount = ? WHERE recipe_id = ? AND food_id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = updt.Exec(recipe_food.UseAmount, recipe_food.RecipeId, recipe_food.FoodId)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Hey, you alter the amount of ingridients in recipe. OK, I accept")
}
