package main

import (
	"log"
	"menu_proposer/pkg/app/foods_app"
	"menu_proposer/pkg/app/recipe_food_app"
	"menu_proposer/pkg/app/recipes_app"
	"net/http"
)

func main() {
	log.Println("Hello, I'm Menu Proposer!")

	//foods_appのリクエスト処理
	http.HandleFunc("/menu_proposer/foods", foods_app.FetchFoods)
	http.HandleFunc("/menu_proposer/insert_food", foods_app.InsertFoods)
	http.HandleFunc("/menu_proposer/delete_food/", foods_app.DeleteFoods)
	http.HandleFunc("/menu_proposer/update_food/", foods_app.UpdateFoods)

	//recipes_appのリクエスト
	http.HandleFunc("/menu_proposer/recipes", recipes_app.FetchRecipes)
	http.HandleFunc("/menu_proposer/insert_recipe", recipes_app.InsertFood)
	http.HandleFunc("/menu_proposer/update_recipe/", recipes_app.UpdateRecipe)
	http.HandleFunc("/menu_proposer/delete_recipe/", recipes_app.DeleteRecipe)

	//recipe_foodのリクエスト
	http.HandleFunc("/menu_proposer/recipe_food", recipe_food_app.FetchRecipesWithFood)
	http.HandleFunc("/menu_proposer/recipe_food/updata_food_strage/", recipe_food_app.UpdateFoodStrage)
	http.HandleFunc("/menu_proposer/recipe_food/update_using_food_quantity", recipe_food_app.UpdateUsingFoodQuantity)

	go http.HandleFunc("/menu_proposer/foods_expiration", foods_app.FetchExpirationFood)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
