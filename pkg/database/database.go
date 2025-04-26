package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"trading-dashboard/pkg/models"

	"github.com/dgraph-io/badger/v4"
)

// DB is the database instance
var DB *badger.DB

// Initialize sets up the database connection
func Initialize() error {
	// Get user-specific app data directory
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("ERROR: Could not get user config directory: %v", err)
		// Fall back to relative path if user directory can't be determined
		appDataDir = "."
	}
	log.Printf("DEBUG: Using app data directory: %s", appDataDir)

	// Create a specific subdirectory for our application
	dataDir := filepath.Join(appDataDir, "TradingDashboard", "data")
	log.Printf("DEBUG: Full data directory path: %s", dataDir)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("ERROR: Failed to create data directory: %v", err)
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Check if directory exists and is writable
	testFile := filepath.Join(dataDir, "test.tmp")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		log.Printf("ERROR: Data directory not writable: %v", err)
		return fmt.Errorf("data directory not writable: %w", err)
	}
	// Clean up test file
	os.Remove(testFile)

	log.Printf("DEBUG: Initializing BadgerDB at: %s", dataDir)

	// Configure BadgerDB options
	opts := badger.DefaultOptions(dataDir).WithLoggingLevel(badger.INFO)

	// Open the database
	log.Printf("DEBUG: Opening BadgerDB...")
	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("ERROR: Failed to open BadgerDB: %v", err)
		return fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	DB = db
	log.Println("DEBUG: BadgerDB initialized successfully")

	// Test database connection with a simple set/get
	testKey := "test_init_key"
	testValue := "test_value"
	log.Printf("DEBUG: Testing database with key: %s", testKey)

	// Try to write to DB
	err = Set(testKey, testValue)
	if err != nil {
		log.Printf("ERROR: Failed to write test value: %v", err)
		return fmt.Errorf("failed to write test value: %w", err)
	}

	// Try to read from DB
	var retrievedValue string
	err = Get(testKey, &retrievedValue)
	if err != nil {
		log.Printf("ERROR: Failed to read test value: %v", err)
		return fmt.Errorf("failed to read test value: %w", err)
	}

	if retrievedValue != testValue {
		log.Printf("ERROR: Test value mismatch. Expected: %s, Got: %s", testValue, retrievedValue)
		return fmt.Errorf("test value mismatch")
	}

	// Clean up test key
	err = Delete(testKey)
	if err != nil {
		log.Printf("WARNING: Failed to delete test key: %v", err)
	}

	log.Println("SUCCESS: BadgerDB test successful - database is operational")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		log.Println("DEBUG: Closing BadgerDB connection...")
		err := DB.Close()
		if err != nil {
			log.Printf("ERROR: Failed to close BadgerDB: %v", err)
		} else {
			log.Println("DEBUG: BadgerDB connection closed successfully")
		}
	}
}

// Helper functions for key-value operations

// Set stores a key-value pair in BadgerDB
func Set(key string, value interface{}) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("ERROR: Failed to marshal value for key %s: %v", key, err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	log.Printf("DEBUG: Setting key: %s (value size: %d bytes)", key, len(bytes))
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), bytes)
	})
}

// Get retrieves a value by key from BadgerDB
func Get(key string, result interface{}) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	log.Printf("DEBUG: Getting key: %s", key)
	var valCopy []byte
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			log.Printf("ERROR: Failed to get item for key %s: %v", key, err)
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		log.Printf("ERROR: DB view operation failed for key %s: %v", key, err)
		return err
	}

	err = json.Unmarshal(valCopy, result)
	if err != nil {
		log.Printf("ERROR: Failed to unmarshal value for key %s: %v", key, err)
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	log.Printf("DEBUG: Successfully retrieved key: %s (value size: %d bytes)", key, len(valCopy))
	return nil
}

// Delete removes a key-value pair from BadgerDB
func Delete(key string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	log.Printf("DEBUG: Deleting key: %s", key)
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// GetByPrefix retrieves all items with a specific prefix
func GetByPrefix(prefix string, results interface{}) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	log.Printf("DEBUG: Getting all keys with prefix: %s", prefix)
	var items [][]byte
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefixBytes := []byte(prefix)
		count := 0
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			count++
			item := it.Item()
			key := item.Key()
			log.Printf("DEBUG: Found key: %s", string(key))

			err := item.Value(func(val []byte) error {
				// Make a copy to use outside the transaction
				valCopy := append([]byte{}, val...)
				items = append(items, valCopy)
				return nil
			})

			if err != nil {
				log.Printf("ERROR: Failed to process value for key %s: %v", string(key), err)
				return err
			}
		}
		log.Printf("DEBUG: Found %d items with prefix: %s", count, prefix)
		return nil
	})

	if err != nil {
		log.Printf("ERROR: Failed to get items with prefix %s: %v", prefix, err)
		return err
	}

	// Handle empty results case - initialize with empty slice
	if len(items) == 0 {
		log.Printf("DEBUG: No items found with prefix: %s, returning empty slice", prefix)

		// Handle different result types for empty results
		switch resultsPtr := results.(type) {
		case *[]*models.RiskAssessment:
			*resultsPtr = []*models.RiskAssessment{}
		case *[]*models.StockRating:
			*resultsPtr = []*models.StockRating{}
		case *[]*models.Trade:
			*resultsPtr = []*models.Trade{}
		default:
			log.Printf("ERROR: Unknown result type for GetByPrefix: %T", results)
			return fmt.Errorf("unknown result type for prefix %s", prefix)
		}
		return nil
	}

	// Unmarshal each item directly into the appropriate type
	switch resultsPtr := results.(type) {
	case *[]*models.RiskAssessment:
		for _, item := range items {
			var assessment models.RiskAssessment
			if err := json.Unmarshal(item, &assessment); err != nil {
				log.Printf("ERROR: Failed to unmarshal RiskAssessment: %v", err)
				return err
			}
			*resultsPtr = append(*resultsPtr, &assessment)
		}
	case *[]*models.StockRating:
		for _, item := range items {
			var rating models.StockRating
			if err := json.Unmarshal(item, &rating); err != nil {
				log.Printf("ERROR: Failed to unmarshal StockRating: %v", err)
				return err
			}
			*resultsPtr = append(*resultsPtr, &rating)
		}
	case *[]*models.Trade:
		for _, item := range items {
			var trade models.Trade
			if err := json.Unmarshal(item, &trade); err != nil {
				log.Printf("ERROR: Failed to unmarshal Trade: %v", err)
				return err
			}
			*resultsPtr = append(*resultsPtr, &trade)
		}
	default:
		log.Printf("ERROR: Unknown result type for GetByPrefix: %T", results)
		return fmt.Errorf("unknown result type for prefix %s", prefix)
	}

	log.Printf("DEBUG: Successfully unmarshaled %d items with prefix: %s", len(items), prefix)
	return nil
}

// GenerateKey generates a key with a prefix and timestamp
func GenerateKey(prefix string) string {
	key := fmt.Sprintf("%s_%s", prefix, time.Now().Format(time.RFC3339))
	log.Printf("DEBUG: Generated key: %s", key)
	return key
}
