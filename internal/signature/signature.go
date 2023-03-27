package signature

import (
	"database/sql"

	"github.com/netmoth/netmoth/internal/storage/postgres"
)

// Detector is ...
type Detector struct {
	postgres.Connect
}

// Request is ...
type Request struct {
	IP         string
	Port       int
	TrackerURL string
	CertSHA1   string `json:",omitempty"`
}

// Detect is ...
type Detect struct {
	Type        string
	Provider    string
	SignatureID int
}

/*
type sqls struct {
	ip      string
	botnet  string
	tracker string
	cert    string
}
*/

// New is ...
func New(conn postgres.Connect) Detector {
	return Detector{
		conn,
	}
}

// Scan is ...
func (d *Detector) Scan(req *Request) ([]Detect, error) {
	var resp []Detect

	// Check IP address
	if req.IP != "" {
		rows, err := d.Conn.Query(`SELECT "signature_ip"."id", "provider_signature"."name", "provider_signature"."type" FROM "signature_ip" 
		INNER JOIN "provider_signature" ON "provider_signature"."id" = "signature_ip"."provider"
		WHERE "signature_ip"."ip" = $1`, req.IP)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			detect := new(Detect)
			err = rows.Scan(&detect.SignatureID, &detect.Provider, &detect.Type)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			resp = append(resp, *detect)
		}
		defer rows.Close()
	}

	// Check BotNet
	if req.IP != "" && req.Port != 0 {
		rows, err := d.Conn.Query(`SELECT "signature_botnet"."id", "provider_signature"."name", "provider_signature"."type" FROM "signature_botnet" 
		INNER JOIN "provider_signature" ON "provider_signature"."id" = "signature_botnet"."provider"
		WHERE "signature_botnet"."ip" = $1 AND "signature_botnet"."port" = $2`, req.IP, req.Port)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			detect := new(Detect)
			err = rows.Scan(&detect.SignatureID, &detect.Provider, &detect.Type)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			resp = append(resp, *detect)
		}
		defer rows.Close()
	}

	// Check tracker
	if req.TrackerURL != "" {
		rows, err := d.Conn.Query(`SELECT "signature_tracker"."id", "provider_signature"."name", "provider_signature"."type" FROM "signature_tracker" 
		INNER JOIN "provider_signature" ON "provider_signature"."id" = "signature_tracker"."provider"
		WHERE "signature_tracker"."url" = $1`, req.TrackerURL)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			detect := new(Detect)
			err = rows.Scan(&detect.SignatureID, &detect.Provider, &detect.Type)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			resp = append(resp, *detect)
		}
		defer rows.Close()
	}

	// Check Certificat
	if req.TrackerURL != "" {
		rows, err := d.Conn.Query(`SELECT "signature_cert"."id", "provider_signature"."name", "provider_signature"."type" FROM "signature_cert" 
		INNER JOIN "provider_signature" ON "provider_signature"."id" = "signature_cert"."provider"
		WHERE "signature_cert"."sha1" = $1`, req.CertSHA1)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			detect := new(Detect)
			err = rows.Scan(&detect.SignatureID, &detect.Provider, &detect.Type)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			resp = append(resp, *detect)
		}
		defer rows.Close()
	}

	return resp, nil
}
