package Handlers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type PersonDetails struct {
	Name              string          `json:"name"`
	AmountSentFromYou ReceiverDetails `json:"amount_sent_from_you"`
	AmountSentToYou   SenderDetails   `json:"amount_sent_to_you"`
	AmountLeft        int64           `json:"amount_left"`
}

type CreateAccount struct {
	Name                  string `json:"name"`
	MobileNumber          string `json:"mobile_number"`
	DateOfCreatingAccount string `json:"date"`
	AmountDeposited       int64  `json:"amount_deposited"`
	SetPin                string `json:"set_pin"` //has to set by the user just to make transaction very safe.
}

type ReceiverDetails struct {
	MobileNumber  string `json:"mobile_number"`
	Amount        int64  `json:"amount"`
	Date          string `json:"date"`
	TransactionID string `json:"transaction_id"`
	Pin           string `json:"pin"` //The user need to enter the pin every time.
}

type SenderDetails struct { //Money you get form another person
	Amount        int64  `json:"amount"` //In cryptocurrencies the transactions are very anonymous, Therefore you won't be able to see the details of the person whoever sent money to you.
	Date          string `json:"date"`
	TransactionID string `json:"transaction_id"`
}

type Block struct {
	PrevHash  string
	Pos       int
	Details   PersonDetails
	Timestamp string
	Hash      string
}

type BlockChain []Block

var Chain BlockChain
var createUser CreateAccount

func CreateUser(w http.ResponseWriter, r *http.Request) { //Creates the user for very first time
	err := json.NewDecoder(r.Body).Decode(&createUser)
	if err != nil {
		log.Printf("Couldn't able to Decode the request %v", err)
		w.Write([]byte("Could not able to create an User"))
	}
	currentTime := time.Now()
	createUser.DateOfCreatingAccount = currentTime.Format("2006-01-02 15:04:05.000000000")
	if err != nil {
		log.Printf("Couldn't able to marshal %v", err)
		w.Write([]byte("User not created"))
	} else {
		w.Write([]byte("User successfully created"))
	}
}

// creates the unique randomIDs used for TransactionIDs
func createRandomId() string {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(newUUID)
}

func SendMoney(w http.ResponseWriter, r *http.Request) { //send money to the receiver
	var receiver ReceiverDetails
	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Printf("Couldn't able to Decode the request %v", err)
		w.Write([]byte("Money did not sent."))
	}
	//creating a random unique ID
	receiver.TransactionID = createRandomId()
	currentTime := time.Now()
	receiver.Date = currentTime.Format("2006-01-02 15:04:05.000000000")
	id := md5.New() //creating new transaction ID
	io.WriteString(id, receiver.Date)
	receiver.TransactionID = fmt.Sprintf("%x", id.Sum(nil))
	//calculating the balance left
	if createUser.AmountDeposited < receiver.Amount {
		log.Printf("Insufficient Funds. Please add the funds to process transaction.")
		w.Write([]byte("Insufficient Funds"))
	} else if receiver.Pin != createUser.SetPin {
		panic("incorrect pin")
		w.Write([]byte("incorrect pin"))
	} else {
		createUser.AmountDeposited = createUser.AmountDeposited - receiver.Amount
	}
	var dupSender SenderDetails //empty struct
	personAssets(receiver, dupSender)
	if err != nil {
		log.Printf("Couldn't able to marshal Receiver Details %v", err)
		w.Write([]byte("Couldn't able to send Money."))
	} else {
		w.Write([]byte("Transaction completed"))
	}
}

func ReceiveMoney(w http.ResponseWriter, r *http.Request) { //receive money from the sender
	var sender SenderDetails
	err := json.NewDecoder(r.Body).Decode(&sender)
	if err != nil {
		log.Printf("Couldn't able to Decode the request %v", err)
		w.Write([]byte("Money did not receive."))
	}
	sender.TransactionID = createRandomId()
	currentTime := time.Now()
	sender.Date = currentTime.Format("2006-01-02 15:04:05.000000000")
	createUser.AmountDeposited = sender.Amount + createUser.AmountDeposited
	var dupReceiver ReceiverDetails // empty struct
	personAssets(dupReceiver, sender)
	if err != nil {
		log.Printf("Couldn't able to marshal sender Details %v", err)
		w.Write([]byte("Couldn't able to receive Money."))
	} else {
		w.Write([]byte("Amount successfully received"))
	}

}

func personAssets(receiver ReceiverDetails, sender SenderDetails) {
	var assets PersonDetails
	assets.Name = createUser.Name
	assets.AmountSentFromYou = receiver
	assets.AmountSentToYou = sender
	assets.AmountLeft = createUser.AmountDeposited
	createBlock(assets)
}

func generateHash(block Block) string {
	bytes, _ := json.Marshal(block.Details)
	data := string(block.Pos) + block.Timestamp + string(bytes) + block.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	block.Hash = hex.EncodeToString(hash.Sum(nil))
	return block.Hash
}

// returns the Previous Hash in the blockchain
func prevHash() string {
	if len(Chain) == 0 {
		return strconv.Itoa(0)
	} else {
		previousHash := Chain[len(Chain)-1].Hash
		return previousHash
	}
}

func createBlock(assets PersonDetails) {
	var block Block
	block.PrevHash = prevHash()
	block.Details = assets
	block.Timestamp = time.Now().Format("2006-01-02 15:04:05.000000000")
	block.Hash = generateHash(block)
	chain(block)
}

func chain(block Block) {
	block.Pos = len(Chain) + 1
	Chain = append(Chain, block)
}

func PrintBlockchain(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(Chain)
	w.Write(res)
}

func Authentication(w http.ResponseWriter, r *http.Request) {

}
