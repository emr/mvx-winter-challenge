package csv

import (
	"os"
	"path/filepath"
	"testing"

	"emr/mvx-winter-challenge/1206-query-transactions/pkg/transactions"
)

func TestWriter(t *testing.T) {
	testDir := t.TempDir()
	address := "erd1test"

	writer, err := OpenFile(testDir, address)
	if err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	txs := []transactions.Transaction{
		{
			TxHash:        "hash1",
			Nonce:         1,
			Timestamp:     1609459200,
			Sender:        "sender1",
			Receiver:      "receiver1",
			Function:      "transfer",
			Value:         "1000",
			MiniBlockHash: "block1",
			Status:        "success",
		},
		{
			TxHash:        "hash2",
			Nonce:         2,
			Timestamp:     1609459200,
			Sender:        "sender2",
			Receiver:      "receiver2",
			Function:      "transfer",
			Value:         "2000",
			MiniBlockHash: "block2",
			Status:        "success",
		},
	}

	if err := writer.WriteTransactionsPage(txs); err != nil {
		t.Fatalf("WriteTransactionsPage() error = %v", err)
	}

	// Read and verify file contents
	content, err := os.ReadFile(filepath.Join(testDir, address+".csv"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	expected := "TxHash,Nonce,Timestamp,Sender,Receiver,Function,Value,MiniBlockHash,Status\n" +
		"hash1,1,2021-01-01T00:00:00Z,sender1,receiver1,transfer,1000,block1,success\n" +
		"hash2,2,2021-01-01T00:00:00Z,sender2,receiver2,transfer,2000,block2,success\n"

	if string(content) != expected {
		t.Errorf("File content = %v, want %v", string(content), expected)
	}

	if err := writer.CloseFile(); err != nil {
		t.Fatalf("CloseFile() error = %v", err)
	}
}
