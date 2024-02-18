package database

import (
	"backend/src/models"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	DB, err = gorm.Open(mysql.Open("user:user_password@tcp(db:3306)/mycooookbook?parseTime=true"), &gorm.Config{})
	// DB, err = gorm.Open(mysql.Open(os.Getenv("DB_USER")+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_ADDRESS")+")/"+os.Getenv("DB_NAME")+"?parseTime=true"), &gorm.Config{})
	// DB, err = gorm.Open(mysql.Open("root:@tcp(host.docker.internal:3306)/mycooookbooklocal?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic("Could not connect with tha database!")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Food{}, models.FoodUnit{}, models.FoodType{}, models.Recipe{}, models.RecipeFoodRelation{})
}

func AutoCreateRecprd() {
	units := []*models.FoodUnit{
		{Unit: "個"}, {Unit: "本"}, {Unit: "匹"}, {Unit: "g"}, {Unit: "ml"},
	}
	DB.Create(&units)

	types := []*models.FoodType{
		{Type: "穀物"}, {Type: "肉"}, {Type: "魚"}, {Type: "野菜"}, {Type: "乳製品"}, {Type: "果物"}, {Type: "卵"},
	}
	DB.Create(&types)
}

// for _, user := range users {
//   user.ID // 1,2,3
// }
