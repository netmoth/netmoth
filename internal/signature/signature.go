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
	TrackerURL string
	CertSHA1   string `json:",omitempty"`
	Port       int
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

/*
func (d *Detector) Scan2(req *Request) ([]Detect, error) {
	var resp []Detect

	// Create a map to store query parameters based on request type
	params := make(map[string]interface{})
	if req.IP != "" {
		params["ip"] = req.IP
	}
	if req.Port != 0 {
		params["port"] = req.Port
	}
	if req.TrackerURL != "" {
		params["url"] = req.TrackerURL
	}
	if req.CertSHA1 != "" {
		params["sha1"] = req.CertSHA1
	}

	// Create a prepared statement with conditional clauses
	stmt, err := d.Conn.Prepare(`
			SELECT s."id", p."name", p."type"
			FROM "provider_signature" AS p
			INNER JOIN %s AS s ON p."id" = s."provider"
			WHERE %s`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query for each request type
	for table, param := range map[string]string{
		"signature_ip":      "ip = :ip",
		"signature_botnet":  "ip = :ip AND port = :port",
		"signature_tracker": "url = :url",
		"signature_cert":    "sha1 = :sha1",
	} {
		if _, ok := params[strings.Split(param, " ")[0][1:]]; !ok {
			continue // Skip if the parameter is not present
		}
		rows, err := stmt.Query(table, param)
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
		rows.Close()
	}

	return resp, nil
}
*/
