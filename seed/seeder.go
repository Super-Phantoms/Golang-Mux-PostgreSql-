package seed

import (
	"log"

	"github.com/golangdevm/fullstack/domain"
	"github.com/golangdevm/fullstack/logger"
	"github.com/golangdevm/fullstack/util"
	"github.com/jinzhu/gorm"
)

var users = []domain.User{
	{
		Username:   "Steven",
		Email:      "steven@gmail.com",
		CustomerId: 2000,
		Password:   "password",
		Role:       "user",
	},
	{
		Username:   "tiger",
		Email:      "tiger@gmail.com",
		Password:   "password",
		CustomerId: 3000,
		Role:       "admin",
	},
}

var posts = []domain.Post{
	{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&domain.Post{}, &domain.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&domain.User{}, &domain.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*
		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
	*/

	for i, user := range users {
		hash_password, hashErr := util.Hash(user.Password)
		if hashErr != nil {
			logger.Error("Error while hash password: " + hashErr.Error())

		}
		user.Password = string(hash_password)
		err = db.Debug().Model(&domain.User{}).Create(&user).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = user.ID

		err = db.Debug().Model(&domain.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	} //end of for
}
