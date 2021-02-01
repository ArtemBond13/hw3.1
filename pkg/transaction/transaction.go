package transaction

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	Id      string
	From    string
	To      string
	Amount  int64
	Created int64
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
func (s *Service) AddTrancsation(id, from, to string, amount, created int64) {
	s.mu.Lock()
	trans := &Transaction{id, from, to, amount, created}
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

// Import если не большие файлы
func (s *Service) Import(r io.Reader) error {
	// Сначало надо прочитать файл
	reader := csv.NewReader(r)
	records := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			records = append(records, record)
			break
		}
		records = append(records, record)
	}

	for _, row := range records {
		transaction, err := s.MapRowToTransaction(row)
		if err != nil {
			log.Println(err)
		}
		if _, _ = s.Register(transaction.From, transaction.To, transaction.Amount); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (s *Service) Import2(r io.Reader) error {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
	}
	for _, row := range records {
		transaction, err := s.MapRowToTransaction(row)
		if err != nil {
			return err
		}
		s.AddTrancsation(transaction.Id, transaction.From, transaction.To, transaction.Amount, transaction.Created)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) MapRowToTransaction(rows []string) (Transaction, error) {
	amount, err := strconv.ParseInt(rows[3], 10, 64)
	if err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	created, err := strconv.ParseInt(rows[4], 10, 64)
	if err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	// createdUnix := time.Unix(created, 0)
	return Transaction{
		rows[0],
		rows[1],
		rows[2],
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
