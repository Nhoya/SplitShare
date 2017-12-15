package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SSSaaS/sssa-golang"
	"github.com/TheZ3ro/go-pgp/pgp"
	"github.com/go-mail/mail"
	"github.com/spf13/viper"
)

type User struct {
	keypath string
	mail    string
	share   string
}

type MailConf struct {
	server   string
	port     int
	username string
	password string
	subject  string
	sender   string
}

func splitSecret(secret []byte) []string {
	var minShares int
	var totalShares int

	fmt.Print("Insert Number of Slices: ")
	fmt.Scan(&totalShares)
	fmt.Print("Insert Minumim number of Slices to decrypt the secret: ")
	fmt.Scan(&minShares)
	if (minShares >= totalShares) || (minShares <= 0) {
		fmt.Println("The minimum nuber of Slices must be lower than the total number of Slices")
		os.Exit(1)
	}

	fmt.Println("[+] Generating Slices")
	shares, err := sssa.Create(minShares, totalShares, string(secret))
	if err != nil {
		fmt.Println("Unable to create Shamir's Shares")
		os.Exit(1)
	}
	return shares
}

func registerUsers(users []User, share string) []User {
	newUser := User{}
	newUser.share = share
	fmt.Print("Enter email: ")
	fmt.Scanln(&newUser.mail)
	fmt.Print("Enter keypath: ")
	fmt.Scanln(&newUser.keypath)
	users = append(users, newUser)
	return users
}

func encryptPassword(person User, conf MailConf) {
	fmt.Println("[+] Encrypting message for " + person.mail)
	key, err := ioutil.ReadFile(person.keypath)
	if err != nil {
		fmt.Println("Unable to read keyfile " + person.keypath)
		os.Exit(1)
	}

	pubEntity, err := pgp.GetEntity(key, []byte{})
	if err != nil {
		fmt.Println("Unable to generate public Entity")
		os.Exit(1)
	}

	encrypted, err := pgp.Encrypt(pubEntity, []byte(person.share))
	if err != nil {
		fmt.Println("Unable to encrypt secret")
		os.Exit(1)
	}

	fmt.Println("[+] Sending mail to: " + person.mail)
	sendViaMail(person, string(encrypted), conf)

}

func sendViaMail(person User, message string, conf MailConf) {
	m := mail.NewMessage()
	m.SetHeader("From", conf.sender)
	m.SetHeader("To", person.mail)
	m.SetHeader("Subject", conf.subject)
	m.SetBody("text/html", message)

	d := mail.NewDialer(conf.server, conf.port, conf.username, conf.password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func initConfig(conf MailConf) MailConf {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf.server = viper.GetString("SMTP.server")
	conf.port = viper.GetInt("SMTP.port")
	conf.username = viper.GetString("SMTP.username")
	conf.password = viper.GetString("SMTP.password")

	conf.subject = viper.GetString("MAIL.subject")
	conf.sender = viper.GetString("MAIL.sender")

	return conf

}

func main() {
	users := []User{}
	conf := MailConf{}

	if len(os.Args) != 2 {
		fmt.Println("You must specify the secret file")
		os.Exit(1)
	}
	secretFile := os.Args[1]
	secret, err := ioutil.ReadFile(secretFile)
	if err != nil {
		panic("Unable to read secret file")
	}
	shares := splitSecret(secret)
	fmt.Println("[+] Adding Users")
	for _, share := range shares {
		users = registerUsers(users, share)
	}

	conf = initConfig(conf)

	for _, person := range users {
		encryptPassword(person, conf)
	}
}
