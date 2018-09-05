package apitool

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// connectToDB do the connection request to the database depending on provided parameters
func ConnectToDB(username string, dbName string, password string, host string) (DB *gorm.DB, err error) {
	log.Println("CONNECTING TO [" + dbName + "] DB...")
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", "host="+host+" user="+username+" dbname="+dbName+" sslmode=disable password="+password)
		if err != nil {
			log.Println("Still trying...")
		} else {
			DB.SingularTable(true)
			log.Println("Database status: [Connected]")
			return DB, nil
		}
		time.Sleep(5 * time.Second)
	}
	return
}
