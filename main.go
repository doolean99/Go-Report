package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

// =================STRUCT untuk Respons GetToken==================//

type ResponseData struct {
	Error bool
	Token string
}

/*================================================================*/
// =================STRUCT untuk Instances=================//

type items struct {
	ID        string `json:"_id"`
	ServerKey string `json:"server_key"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	WaName    string `json:"wa_name"`
	Webhook   string `json:"webhook"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Version   int    `json:"__v"`
}

type InstacesData struct {
	Error   bool
	Message string
	Data    []items `json:"data"`
}

/*================================================================*/
/* function untuk mengubah value variable di file .env*/
func ReadToken(nameFile string) string {
	file, err := os.Open(nameFile + ".txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}

/* function untuk menulis file TOKEN.txt*/
func WriteToken(nameFile string, Values string) {
	file, err := os.Create(nameFile + ".txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(Values)
	fmt.Println("Success Create Token")
}

// function untuk token login //
func GetToken(Server string) string {
	// Mengambil value dari variable .env //
	username := ReadEnv("username")
	password := ReadEnv("password")
	url := ReadEnv(Server)
	var JsonData ResponseData
	client := resty.New()
	fmt.Println(url)
	// data JSON/array map di golang
	data := map[string]interface{}{
		"username": username,
		"password": password,
	}
	// POST data To Url
	response, err := client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(url + "/signin")
	fmt.Println(response)
	if err != nil {
		panic(err)
	}
	resData := response.Body()
	error := json.Unmarshal([]byte(resData), &JsonData)
	if error != nil {
		panic(error)
	}
	return JsonData.Token
}

// function untuk membaca .env //
func ReadEnv(variable string) string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	read := os.Getenv(variable)

	return read
}

// function generate Token //
func GenerateToken() {
	SignIn := GetToken("server50")
	SignIn03 := GetToken("server03")
	WriteToken("tokenRhea", SignIn)
	WriteToken("token03", SignIn03)
}

func checkInstances() {
	// slice array //
	names := []string{}
	// slice array //
	Token := ReadToken("token03")
	Url := ReadEnv("server03")
	client := resty.New()
	var JsonData InstacesData
	response, err := client.R().SetHeader("Authorization", "Bearer "+Token).Get(Url + "/instance/all")
	if err != nil {
		panic(err)
	}
	resData := response.Body()
	error := json.Unmarshal([]byte(resData), &JsonData)

	if error != nil {
		panic(error)
	}
	for _, item := range JsonData.Data {
		if item.Status == "ONLINE" {
			names = append(names, item.Name)
		}
	}
	if len(names) < 400 {
		Report(len(names))
	} else if len(names) >= 400 {
		fmt.Println("INSTANCES AMANN")
	}
}

func Report(length int) {
	Token := ReadToken("tokenRhea")
	Url := ReadEnv("server50")
	client := resty.New()
	group := ReadEnv("group_key")
	key := ReadEnv("rhea_key")
	data := map[string]interface{}{
		"id":      group,
		"message": "JUMLAH INSTANCES " + string(rune(length)),
	}

	response, err := client.R().SetHeader("Authorization", "Bearer "+Token).SetBody(data).Post(Url + "group/ListmemGroup?key=" + key)

	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}

func main() {
	localeTimes, _ := time.LoadLocation("Asia/Jakarta")

	scheduler := cron.NewWithLocation(localeTimes)

	scheduler.AddFunc("0 0 0 * * *", GenerateToken)
	scheduler.AddFunc("0 0 6 * * *", checkInstances)

	go scheduler.Start()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
