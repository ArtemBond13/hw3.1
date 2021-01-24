package transaction

import (
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
			strconv.FormatInt(transaction.Amount, 10),
			strconv.FormatInt(transaction.Created, 10),
		}
		records = append(records, record)
	}
	s.mu.Unlock()
	w := csv.NewWriter(writer)
	return w.WriteAll(records)

}

// Чтение файла
//func (s *Service) Import(reader io.Reader) error {
//	file, err := os.Open("export.scv")
//	if err != nil{
//		log.Println(err)
//	}
//
//	defer func(c io.Closer) {
//		if cerr := c.Close(); cerr != nil{
//			log.Println(err)
//			if err == nil {
//				err = cerr
//			}
//		}
//	}(file)
//
//	// Срез для хранения содержимот
//	content := make([]byte, 0)
//
//	// буфер для чтения
//	buf := make([]byte, 4096) // 4096 - количство байт
//
//	for {
//		n, err := file.Read(buf)
//		if err != nil {
//			// io.EOF - ошибка, сигнализирующая о том, что дочитали данные до конца (файл закончился)
//			if err != io.EOF {
//				log.Println(err)
//			}
//			// "перекладываем" данные из буфера в слайс со всем содержимым
//			content = append(content, buf[:n]...)
//			break
//		}
//		content = append(content, buf[:n]...)
//	}
//
//}

func (s *Service) Import(r io.Reader) error {
	//// Открыть файл
	//file, err := os.Open("import.scv")
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//defer func(c io.Closer) {
	//	if cerr := c.Close(); cerr != nil {
	//		log.Println(err)
	//		if err == nil {
	//			err = cerr
	//		}
	//	}
	//}(file)

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
		s.Register(transaction.From, transaction.To, transaction.Amount)

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
		Id:      rows[1],
		From:    rows[2],
		To:      rows[3],
		Amount:  amount,
		Created: time.Unix(created, 0),
	}, nil
}
