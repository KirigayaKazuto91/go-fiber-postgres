package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KirigayaKazuto91/go-fiber-postgres/models"
	"github.com/KirigayaKazuto91/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Car struct{
	Merk 	string	`json:"merk"`
	Tipe	string	`json:"tipe"`
	Warna	string	`json:"warna"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateCar(context *fiber.Ctx) error{
	car := Car{}

	err := context.BodyParser(&car)
	if err != nil{
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message" : "Request Failed"})
			return err
	}

	err = r.DB.Create(&car).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message" : "could not create book"})
			return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"car has been added"})
		return nil
}

func (r *Repository) GetCars(context *fiber.Ctx) error {
	carModels := &[]models.Cars{}

	err := r.DB.Find(carModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get cars"})
			return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"car fetched succesfully",
		"data": carModels,
	})
	return nil
}

func (r *Repository) DeleteCar(context *fiber.Ctx) error {
	carModel := models.Cars{}
	id := context.Params("id")
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(carModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"could not delete car",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"car delete successfully",
	})
	return nil
}

func (r *Repository) GetCarById(context *fiber.Ctx) error {
	id := context.Params("id")
	carModel := &models.Cars{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"id cannot be empty",
		})
		return nil
	}
	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(carModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get the car"})
			return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"book id fetched successfully",
		"data":carModel,
	})
	return nil
}

func(r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_cars", r.CreateCar)
	api.Delete("/delete_car/:id", r.DeleteCar)
	api.Get("/get_cars/:id", r.GetCarById)
	api.Get("/cars", r.GetCars)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal(err)
	}

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Couldnt Load the Database")
	}

	err = models.MigrateCars(db)
	if err != nil {
		log.Fatal("could not migrate DB")
	}

	r := Repository{
		DB: db, 
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}