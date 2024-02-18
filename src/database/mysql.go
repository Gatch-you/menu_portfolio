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
