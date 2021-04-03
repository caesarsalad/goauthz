package authorization

import (
	"log"

	"github.com/caesarsalad/goauthz/database"
	"gorm.io/gorm"
)

func InitStaticTypesDB() {
	var methods []database.HTTPMethod
	var err error
	for k, v := range HttpMethodIDMap {
		method := database.HTTPMethod{Model: gorm.Model{ID: v}, Method: k}
		methods = append(methods, method)
	}
	err = database.DB.Create(&methods).Error
	if err != nil {
		log.Println("error while inserting HTTPMethods to DB ", err)
	}

	var meta_locations []database.MetaLocation
	for k, v := range MetaLocationIDMap {
		meta_location := database.MetaLocation{Model: gorm.Model{ID: v}, MetaLocation: k}
		meta_locations = append(meta_locations, meta_location)
	}
	err = database.DB.Create(&methods).Error
	if err != nil {
		log.Println("error while inserting MetaLocations to DB ", err)
	}
}
