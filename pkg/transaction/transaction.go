package transaction

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	XMLName string `xml:"transaction"`
	Id      string `json:"id" xml:"id"`
	From    string `json:"from" xml:"from"`
	To      string `json:"to" xml:"to"`
	Amount  int64 `json:"amount" xml:"amount"`
	Created int64 `json:"created" xml:"created"`
}

type Service struct {
	mu           sync.Mutex
	transactions []*Transaction
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Register(from, to string, amount int64) (string, error) {
	t := &Transaction{
		Id:      "xxxx", //FIXME: use uuid later
		From:    from,
		To:      to,
		Amount:  amount,
		Created: time.Now().Unix(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions = append(s.transactions, t)

	return t.Id, nil
}

// AddTrancsation добавляет транзакцию в историю
func (s *Service) AddTrancsation(xmlName, id, from, to string, amount, created int64) {
	s.mu.Lock()
	trans := &Transaction{xmlName,id, from, to,amount, created}
	s.transactions = append(s.transactions, trans)
	s.mu.Unlock()

}

func (s *Service) Export(writer io.Writer) error {
	s.mu.Lock()
	if len(s.transactions) == 0 {
		s.mu.Unlock()
		return nil
	}

	records := make([][]string, len(s.transactions))
	for _, transaction := range s.transactions {
		record := []string{
			transaction.Id,
			transaction.From,
			transaction.To,
			strconv.FormatInt(transaction.Amount, 10), // преобразование числа в строку
			strconv.FormatInt(transaction.Created, 10),
		}
		records = append(records, record)
	}
	s.mu.Unlock()
	w := csv.NewWriter(writer)
	return w.WriteAll(records) // не используем defer,потому что тогда lock будет висеть доокончания записи
}

func (s *Service) ExportJSON(filename string) error {
	s.mu.Lock()
	if len(s.transactions) == 0 {
		s.mu.Unlock()
		return nil
	}

	encoded, err := json.Marshal(s.transactions)
	if err != nil {
		log.Println(err)
	}
	s.mu.Unlock()
	if err = ioutil.WriteFile(filename, encoded, 0666); err != nil {
		log.Println(err)
	}

	return nil
}

func (s *Service) ExportXML(filename string) error {
	s.mu.Lock()
	if len(s.transactions) == 0 {
		s.mu.Unlock()
		return nil
	}

	encoded, err := json.Marshal(s.transactions)
	if err != nil {
		log.Println(err)
	}
	encoded = append([]byte(xml.Header), encoded...)
	s.mu.Unlock()
	if err = ioutil.WriteFile(filename, encoded, 0666); err != nil {
		log.Println(err)
	}

	return nil
}

func (s * Service) Import(r io.Reader) error {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
	}
	for _,row:=range records{
		transaction, err := s.MapRowToTransaction(row)
		if err != nil{
			return  err
		}
		s.AddTrancsation(transaction.XMLName, transaction.Id, transaction.From, transaction.To, transaction.Amount, transaction.Created)
	}
	if err != nil{
		return err
	}
	return nil
}

func (s * Service) ImportJSON(filename io.Reader) error {
	file, err := ioutil.ReadAll(filename)
	if err != nil {
		fmt.Printf("Cannot read file %s\n", filename)
		log.Println(err)
		return err
	}
	var decoded []Transaction

	// Важно: передаём указатель, чтобы функция смогла записать данные
	if err = json.Unmarshal(file, &decoded); err != nil {
		fmt.Printf("Cannot unmarshaling file %s\n", filename)
		log.Println(err)
		return err
	}
	log.Printf("%#v\n", decoded)

	for _, transaction := range decoded{
		s.AddTrancsation(transaction.XMLName, transaction.Id, transaction.From, transaction.To, transaction.Amount, transaction.Created)
	}
	return nil
}

func (s * Service) ImportXML(filename io.Reader) error {
	file, err := ioutil.ReadAll(filename)
	if err != nil {
		fmt.Printf("Cannot read file %s\n", filename)
		log.Println(err)
		return err
	}
	var decoded []Transaction

	// Важно: передаём указатель, чтобы функция смогла записать данные
	if err = xml.Unmarshal(file, &decoded); err != nil {
		fmt.Printf("Cannot unmarshaling file %s\n", filename)
		log.Println(err)
		return err
	}
	log.Printf("%#v\n", decoded)

	for _, transaction := range decoded{
		s.AddTrancsation(transaction.XMLName, transaction.Id, transaction.From, transaction.To, transaction.Amount, transaction.Created)
	}
	return nil
}

func (s *Service) MapRowToTransaction(rows []string) (Transaction, error) {
	amount, err := strconv.ParseInt(rows[4], 10, 64)
	if err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	created, err := strconv.ParseInt(rows[5], 10, 64)
	if err != nil {
		log.Println(err)
		return Transaction{}, err
	}

	return Transaction{
		rows[0],
		rows[1],
		rows[2],
		rows[3],
		amount,
		created,
	}, nil
}

func Compress(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	out, err := os.Create(filename + ".gz")
	if err != nil {
		log.Println(err)
	}

	gzout := gzip.NewWriter(out)
	_, err = io.Copy(gzout, file)
	gzout.Close()
	return err
}