package main

import (
	"backend/internal/app/foods_app"
	"backend/internal/app/recipe_food_app"
	"backend/internal/app/recipes_app"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello, I'm Menu Proposer!")

	// CORSミドルウェア関数
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// レスポンスヘッダーに対するCORS設定を追加
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	//foods_appのリクエスト処理
	http.HandleFunc("/backend/foods", foods_app.FetchFoods)
	http.HandleFunc("/backend/insert_food", foods_app.InsertFoods)
	http.HandleFunc("/backend/delete_food", foods_app.DeleteFoods)
	http.HandleFunc("/backend/update_food", foods_app.UpdateFoods)

	//recipes_appのリクエスト
	http.HandleFunc("/backend/recipes", recipes_app.FetchRecipes)
	http.HandleFunc("/backend/insert_recipe", recipes_app.InsertRecipe)
	http.HandleFunc("/backend/update_recipe", recipes_app.UpdateRecipe)
	http.HandleFunc("/backend/delete_recipe", recipes_app.DeleteRecipe)

	//recipe_foodのリクエスト
	http.HandleFunc("/backend/recipe_food", recipe_food_app.FetchRecipesWithFood)
	http.HandleFunc("/backend/recipe_food/update_food_storage/", recipe_food_app.UpdateFoodStorage)
	http.HandleFunc("/backend/recipe_food/update_using_food_quantity", recipe_food_app.UpdateUsingFoodQuantity)
	http.HandleFunc("/backend/recipe_food/insert_use_food_array", recipe_food_app.InsertUseFoodArray)
	http.HandleFunc("/backend/recipes/", recipe_food_app.FetchRecipeDetail)
	http.HandleFunc("/backend/delete_using_food", recipe_food_app.DeleteUsingFood)
	http.HandleFunc("/backend/recipe_food/insert_use_food", recipe_food_app.InsertUseFood)
	http.HandleFunc("/backend/recipe_food/foods_expiration", recipe_food_app.ShowFoodsWithExpiration)
	//goroutineで定時に発火させる
	go recipe_food_app.FetchExpirationFood(nil, nil)

	// CORSミドルウェアを適用
	handler := corsMiddleware(http.DefaultServeMux)
	http.ListenAndServe(":8080", handler)
}
