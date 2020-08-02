package transaction

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	Id string `json:"id"`
	From string `json:"from"`
	To string `json:"to"`
	Amount int64 `json:"amount"`
	Created int64 `json:"created"`
}

type Service struct {
	mu sync.Mutex
	Transactions []*Transaction
}

func NewService() *Service{
	return &Service{}
}

func (s *Service) Register(from, to string, amount int64) (string, error) {
	t := &Transaction{
		Id: "xxxxxx",
		From: from,
		To: to,
		Amount: amount,
		Created: time.Now().Unix(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Transactions = append(s.Transactions, t)

	return t.Id, nil
}

func (s *Service) ExportCSV(writer io.Writer) error {
	s.mu.Lock()
	if len(s.Transactions) == 0 {
		s.mu.Unlock()
		return nil
	}
	records := make([][]string, len(s.Transactions))
	for _, t := range s.Transactions {
		record := []string{
			t.Id,
			t.From,
			t.To,
			strconv.FormatInt(t.Amount, 10),
			strconv.FormatInt(t.Created, 10),
		}
		records = append(records, record)
	}
	s.mu.Unlock()

	w := csv.NewWriter(writer)
	return w.WriteAll(records)
}
func (s *Service) ExportJSON(writer io.Writer) error {
	s.mu.Lock()
	if len(s.Transactions) == 0 {
		s.mu.Unlock()
		return nil
	}
	encoded, err := json.Marshal(s.Transactions)
	if err != nil {
		log.Println(err)
		return err
	}
	s.mu.Unlock()
	_, err = io.WriteString(writer, string(encoded))
	return err
}

func MapRowToTransaction(record []string) *Transaction {
	if len(record) != 5 {
		return nil
	}
	amount, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		log.Println(err)
	}
	created, err := strconv.ParseInt(record[4], 10, 64)
	if err != nil {
		log.Println(err)
	}
	return &Transaction{
		Id:      record[0],
		From:    record[1],
		To:      record[2],
		Amount:  amount,
		Created: created,
	}
}

func (s *Service) ImportCSV(reader io.Reader) error {
	r := csv.NewReader(reader)
	records := make([][]string, 0)
	for {
		record, err := r.Read()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return err
			}
			records = append(records, record)
			break
		}
		records = append(records, record)
	}
	s.mu.Lock()
	for _, v := range records {
		s.Transactions = append(s.Transactions, MapRowToTransaction(v))
	}
	s.mu.Unlock()
	return nil
}

func (s *Service) ImportJSON(reader io.Reader) error {
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return err
	}
	s.mu.Lock()
	err = json.Unmarshal(d, &s.Transactions)
	s.mu.Unlock()
	return err
}