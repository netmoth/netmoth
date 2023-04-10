package signature

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
)

// Provider is ...
type Provider struct {
	Type      string
	Source    string
	Signature string
	ID        int
}

// Update is ...
func (d *Detector) Update() error {
	rows, err := d.Conn.Query(`SELECT "id", "type", "source", "signature" FROM "provider_signature"`)
	if err != nil {
		return err
	}

	for rows.Next() {
		provider := new(Provider)
		err = rows.Scan(&provider.ID, &provider.Type, &provider.Source, &provider.Signature)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		signature, err := regexp.Compile(provider.Signature)
		if err != nil {
			return err
		}

		resp, err := http.Get(provider.Source)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}

			matches := signature.FindStringSubmatch(scanner.Text())
			if len(matches) == 0 {
				continue
			}

			result := make(map[string]string)
			for i, name := range signature.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = matches[i]
				}
			}

			switch provider.Type {
			case "ip":
				go d.updateIP(provider.ID, result)
			case "botnet":
				go d.updateBotnet(provider.ID, result)
			case "cert":
				go d.updateCert(provider.ID, result)
			case "tracker":
				go d.updateTracker(provider.ID, result)
			}
		}

		if _, err := d.Conn.Exec(`UPDATE "provider_signature" SET "updated" = NOW() WHERE "id" = $1`, provider.ID); err != nil {
			fmt.Print(err)
		}

	}
	defer rows.Close()
	return nil
}

func (d *Detector) updateIP(providerID int, result map[string]string) {
	var keyID string
	err := d.Conn.QueryRow(`SELECT "id" FROM "signature_ip" WHERE "ip" = $1`, result["ip"]).Scan(&keyID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Print(err)
	}
	if keyID == "" {
		_, err := d.Conn.Exec(`INSERT INTO "signature_ip" ("ip", "provider") VALUES ($1, $2)`, result["ip"], providerID)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (d *Detector) updateBotnet(providerID int, result map[string]string) {
	var keyID string
	err := d.Conn.QueryRow(`SELECT "id" FROM "signature_botnet" WHERE "ip" = $1 AND "port" = $2`, result["ip"], result["port"]).Scan(&keyID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Print(err)
	}
	if keyID == "" {
		_, err := d.Conn.Exec(`INSERT INTO "signature_botnet" ("ip", "port", "provider") VALUES ($1, $2, $3)`, result["ip"], result["port"], providerID)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (d *Detector) updateCert(providerID int, result map[string]string) {
	var keyID string
	err := d.Conn.QueryRow(`SELECT "id" FROM "signature_cert" WHERE "sha1" = $1`, result["sha1"]).Scan(&keyID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Print(err)
	}
	if keyID == "" {
		_, err := d.Conn.Exec(`INSERT INTO "signature_cert" ("sha1", "name", "provider") VALUES ($1, $2, $3)`, result["sha1"], result["name"], providerID)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (d *Detector) updateTracker(providerID int, result map[string]string) {
	var keyID string
	err := d.Conn.QueryRow(`SELECT "id" FROM "signature_tracker" WHERE "url" = $1`, result["url"]).Scan(&keyID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Print(err)
	}
	if keyID == "" {
		_, err := d.Conn.Exec(`INSERT INTO "signature_tracker" ("url", "provider") VALUES ($1, $2)`, result["url"], providerID)
		if err != nil {
			fmt.Print(err)
		}
	}
}
