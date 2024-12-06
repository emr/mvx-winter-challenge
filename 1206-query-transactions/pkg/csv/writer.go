package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"emr/mvx-winter-challenge/1206-query-transactions/pkg/transactions"
)

type Writer struct {
	File      *os.File
	csvWriter *csv.Writer
}

func OpenFile(resultsDir, address string) (*Writer, error) {
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create results directory: %w", err)
	}
	file, err := os.Create(filepath.Join(resultsDir, address+".csv"))
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	csvWriter := csv.NewWriter(file)

	headers := []string{"TxHash", "Nonce", "Timestamp", "Sender", "Receiver", "Function", "Value", "MiniBlockHash", "Status"}
	if err := csvWriter.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}
	return &Writer{
		File:      file,
		csvWriter: csvWriter,
	}, nil
}

func (w *Writer) WriteTransactionsPage(txs []transactions.Transaction) error {
	for _, tx := range txs {
		record := []string{
			tx.TxHash,
			fmt.Sprintf("%d", tx.Nonce),
			time.Unix(tx.Timestamp, 0).UTC().Format(time.RFC3339),
			tx.Sender,
			tx.Receiver,
			tx.Function,
			tx.Value,
			tx.MiniBlockHash,
			tx.Status,
		}
		if err := w.csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}
	w.csvWriter.Flush()
	return nil
}

func (w *Writer) CloseFile() error {
	if w.File != nil {
		return w.File.Close()
	}
	return nil
}
