package db

import (
	"log"

	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded. Skipping seed.")
		return nil
	}

	log.Println("Seeding initial data (Users, Tags, Articles, Comments)...")

	// 1. Seed Users
	user1 := models.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashedpassword123",
		Bio:      "I work at State Farm and like coding in Go.",
		Image:    "https://fastly.picsum.photos/id/575/200/200.jpg?hmac=u8uMtAWK-6Ug08Vo4nf84xQLlwJqyrXpfzsU9a3YpCY",
	}
	user2 := models.User{
		Username: "janedoe",
		Email:    "jane@example.com",
		Password: "hashedpassword123",
		Bio:      "Tech enthusiast and Echo framework fan.",
		Image:    "https://fastly.picsum.photos/id/42/200/200.jpg?hmac=jc_eDuYgXmIOC_4gl2wEY0jgxC2rMPJbDF6QJdynR7Q",
	}

	if err := db.Create(&user1).Error; err != nil {
		return err
	}
	if err := db.Create(&user2).Error; err != nil {
		return err
	}

	// 2. Seed Tags
	tagGo := models.Tag{Name: "golang"}
	tagEcho := models.Tag{Name: "echo"}
	tagPostgres := models.Tag{Name: "postgres"}
	tagRealWorld := models.Tag{Name: "realworld"}

	if err := db.Create(&tagGo).Error; err != nil {
		return err
	}
	if err := db.Create(&tagEcho).Error; err != nil {
		return err
	}
	if err := db.Create(&tagPostgres).Error; err != nil {
		return err
	}
	if err := db.Create(&tagRealWorld).Error; err != nil {
		return err
	}

	// 3. Seed Articles
	article1 := models.Article{
		Slug:           "how-to-train-your-dragon",
		Title:          "How to train your dragon",
		Description:    "Ever wonder how?",
		Body:           "It takes a lot of patience and dragon nip.",
		FavoritesCount: 2,
		AuthorID:       user1.ID,
		Tags:           []models.Tag{tagGo, tagEcho},
	}
	article2 := models.Article{
		Slug:           "learning-go-echo-framework",
		Title:          "Learning Go Echo Framework",
		Description:    "A beginner's guide to building REST APIs with Go and Echo.",
		Body:           "Echo is a high performance, extensible, minimalist Go web framework.",
		FavoritesCount: 5,
		AuthorID:       user2.ID,
		Tags:           []models.Tag{tagGo, tagPostgres, tagRealWorld},
	}

	if err := db.Create(&article1).Error; err != nil {
		return err
	}
	if err := db.Create(&article2).Error; err != nil {
		return err
	}

	// 4. Seed Comments
	comment1 := models.Comment{
		Body:      "His name was King! Great article.",
		ArticleID: article1.ID,
		AuthorID:  user2.ID,
	}
	comment2 := models.Comment{
		Body:      "Echo framework is super fast and clean!",
		ArticleID: article2.ID,
		AuthorID:  user1.ID,
	}

	if err := db.Create(&comment1).Error; err != nil {
		return err
	}
	if err := db.Create(&comment2).Error; err != nil {
		return err
	}

	log.Println("Database seeding completed successfully.")
	return nil
}
