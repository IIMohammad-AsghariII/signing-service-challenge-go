package persistence

import (
	"database/sql"
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDeviceRepository implements the DeviceRepository interface for SQLite
type SQLiteDeviceRepository struct {
	db *sql.DB
}

// NewSQLiteDeviceRepository creates a new instance of SQLiteDeviceRepository
func NewSQLiteDeviceRepository(dataSourceName string) (DeviceRepository, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the devices table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		label TEXT,
		algorithm TEXT,
		publicKey TEXT,
		privateKey TEXT,
		lastSignature TEXT,
		signatureCount INTEGER DEFAULT 0
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &SQLiteDeviceRepository{db: db}, nil
}

// AddDevice saves a new SignatureDevice to the repository
func (repo *SQLiteDeviceRepository) AddDevice(id, label string, algorithm domain.AlgorithmType, publicKey, privateKey, lastSignature string) (*domain.SignatureDevice, error) {
	device := domain.NewSignatureDevice(id, label, algorithm, publicKey, privateKey, lastSignature)

	insertSQL := `INSERT INTO devices (id, label, algorithm, publicKey, privateKey, lastSignature, signatureCount) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(insertSQL, device.GetID(), device.GetLabel(), device.GetAlgorithm(), device.GetPublicKey(), device.GetPrivateKey(), device.GetLastSignature(), device.GetSignatureCount())
	if err != nil {
		return nil, err
	}

	return device, nil
}

// GetDevice retrieves a SignatureDevice by its ID
func (repo *SQLiteDeviceRepository) GetDevice(id string) (*domain.SignatureDevice, error) {
	querySQL := `SELECT id, label, algorithm, publicKey, privateKey, lastSignature, signatureCount FROM devices WHERE id = ?`
	row := repo.db.QueryRow(querySQL, id)

	var label, algorithm, publicKey, privateKey, lastSignature string
	var signatureCount uint64

	err := row.Scan(&id, &label, &algorithm, &publicKey, &privateKey, &lastSignature, &signatureCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("device not found")
		}
		return nil, err
	}

	// Create a new SignatureDevice using the retrieved values
	device := domain.NewSignatureDevice(id, label, domain.AlgorithmType(algorithm), publicKey, privateKey, lastSignature)
	device.SetSignatureCount(signatureCount)

	return device, nil
}

// ListDevices returns all SignatureDevices in the repository
func (repo *SQLiteDeviceRepository) ListDevices() ([]*domain.SignatureDevice, error) {
	var devices []*domain.SignatureDevice

	querySQL := `SELECT id, label, algorithm, publicKey, privateKey, lastSignature, signatureCount FROM devices`
	rows, err := repo.db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, label, algorithm, publicKey, privateKey, lastSignature string
		var signatureCount uint64

		if err := rows.Scan(&id, &label, &algorithm, &publicKey, &privateKey, &lastSignature, &signatureCount); err != nil {
			return nil, err
		}

		device := domain.NewSignatureDevice(id, label, domain.AlgorithmType(algorithm), publicKey, privateKey, lastSignature)
		device.SetSignatureCount(signatureCount)

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// IncrementSignatureCount updates the signature count of a device
func (repo *SQLiteDeviceRepository) IncrementSignatureCount(id string) error {
	updateSQL := `UPDATE devices SET signatureCount = signatureCount + 1 WHERE id = ?`
	_, err := repo.db.Exec(updateSQL, id)
	return err
}

// UpdateLastSignature updates the last signature of a device
func (repo *SQLiteDeviceRepository) UpdateLastSignature(id string, lastSignature string) error {
	updateSQL := `UPDATE devices SET lastSignature = ? WHERE id = ?`
	_, err := repo.db.Exec(updateSQL, lastSignature, id)
	return err
}

// Close closes the database connection
func (repo *SQLiteDeviceRepository) Close() error {
	return repo.db.Close()
}
