package recipes_app

import (
	db "backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// レシピの構造体json形式のデータ変換等も行う
type Recipe struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         any    `json:"image"`
	Making_method string `json:"making_method"`
}

// レシピの一覧表示。frontにて表示件数を絞る必要が出てくるかも　→ フロントにて考慮のはず
func FetchRecipes(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		log.Fatal(err.Error())
	}

	recipeArgs := make([]Recipe, 0)
	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Image, &recipe.Making_method)
		if err != nil {
			log.Fatal(err.Error())
		}
		recipeArgs = append(recipeArgs, recipe)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	json.NewEncoder(w).Encode(recipeArgs)
}

// レシピの検索に使うが、他の検索アルゴリズムに置き換わる可能性大
func FetchRecipeByKey(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()
}

// 新しいレシピの項目追加
// curl -X POST -H "Content-Type: application/json" -d '{"id": 2, "name": "カレー", "description": "日 本家庭の一般的料理。", "image": null, "making_method": "hoge"}' http://localhost:8080/backend/insert_recipe
func InsertFood(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
	defer db.Close()

	var recipe Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Fatal(err.Error())
	}

	crat, err := db.Prepare("INSERT INTO recipes (id, name, description, image, making_method) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = crat.Exec(recipe.ID, recipe.Name, recipe.Description, recipe.Image, recipe.Making_method)
	if err != nil {
		log.Fatal(err.Error())
	}

	bytes, err := json.Marshal(recipe)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.Write([]byte("Insert a New Recipe\n"))
	w.Write(bytes)

}

// レシピの変更をこのコードにて処理する。
// curl -X PUT -H "Content-Type: application/json" -d '{"id": n, "name": "hoge", "description": "hoge", "image": null, "making_method": "hoge"}' http://localhost:8080/backend/update_recipe/(id)
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/backend/update_recipe/")

	db := db.Connect()
	defer db.Close()

	var recipe Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Fatal(err.Error())
	}

	update, err := db.Prepare("UPDATE recipes SET name = ?, description = ?, image = ?, making_method = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = update.Exec(recipe.Name, recipe.Description, recipe.Image, recipe.Making_method, id)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	fmt.Fprintf(w, "%s has been updated.", recipe.Name)
}

// レシピのデータベースのフィールドそのものを削除する。再びその食品を使うには再度Insertrecipeを叩かなければならなくなる。
// curl -X DELETE localhost:8080/backend/delete_recipe/(id)
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/backend/delete_recipe/")

	db := db.Connect()
	defer db.Close()

	delt, err := db.Prepare("DELETE FROM recipes WHERE name = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = delt.Exec(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	fmt.Printf("Recipe which you select has been deleted")
}
