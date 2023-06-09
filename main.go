package main

import (
	"backend/pkg/app/foods_app"
	"backend/pkg/app/recipe_food_app"
	"backend/pkg/app/recipes_app"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello, I'm Menu Proposer!")

	//foods_appのリクエスト処理
	http.HandleFunc("/backend/foods", foods_app.FetchFoods)
	http.HandleFunc("/backend/insert_food", foods_app.InsertFoods)
	http.HandleFunc("/backend/delete_food", foods_app.DeleteFoods)
	http.HandleFunc("/backend/update_food", foods_app.UpdateFoods)
	http.HandleFunc("/backend/search_foods/", foods_app.SearchFoods)

	//recipes_appのリクエスト
	http.HandleFunc("/backend/recipes", recipes_app.FetchRecipes)
	http.HandleFunc("/backend/insert_recipe", recipes_app.InsertRecipe)
	http.HandleFunc("/backend/update_recipe", recipes_app.UpdateRecipe)
	http.HandleFunc("/backend/delete_recipe", recipes_app.DeleteRecipe)

	//recipe_foodのリクエスト
	http.HandleFunc("/backend/recipe_food", recipe_food_app.FetchRecipesWithFood)
	http.HandleFunc("/backend/recipe_food/update_food_storage/", recipe_food_app.UpdateFoodStorage)
	http.HandleFunc("/backend/recipe_food/update_using_food_quantity", recipe_food_app.UpdateUsingFoodQuantity)

	//goroutineで定時に発火させる
	go recipe_food_app.FetchExpirationFood(nil, nil)

	http.ListenAndServe(":8080", nil)

}
