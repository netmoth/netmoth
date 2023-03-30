-- +goose Up
-- +goose StatementBegin
INSERT INTO "provider_signature" ("id", "name", "type", "source", "signature", "active", "created", "updated") VALUES
(1, 'Abuse.ch_Feodo_BlockIP', 'ip', 'https://feodotracker.abuse.ch/downloads/ipblocklist.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', TRUE, CURRENT_TIMESTAMP(6), NULL),
(2, 'Blocklist.de_BlockIP', 'ip', 'https://lists.blocklist.de/lists/all.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', TRUE, CURRENT_TIMESTAMP(6), NULL),
(3, 'EmergingThreats_BlockIP', 'ip', 'https://rules.emergingthreats.net/fwrules/emerging-Block-IPs.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$|\/\d{2})){4})', TRUE, CURRENT_TIMESTAMP(6), NULL),
(4, 'EmergingThreats_CompromisedIP', 'ip', 'https://rules.emergingthreats.net/blockrules/compromised-ips.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', TRUE, CURRENT_TIMESTAMP(6), NULL),
(5, 'Abuse.ch_Botnet_C2_IP_Denylist', 'botnet', 'https://sslbl.abuse.ch/blacklist/sslipblacklist.csv', '^(?P<date>.{19}),(?P<ip>.*),(?P<port>.*)', TRUE, CURRENT_TIMESTAMP(6), NULL),
(6, 'Abuse.ch_SSL_Certificate_BlockHash', 'cert', 'https://sslbl.abuse.ch/blacklist/sslblacklist.csv', '^(?P<date>.{19}),(?P<sha1>.*),(?P<name>.*)', TRUE, CURRENT_TIMESTAMP(6), NULL),
(7, 'CyberCrime_Tracker_BlockUrl', 'tracker', 'http://cybercrime-tracker.net/all.php', '^(?P<url>.{5,})', TRUE, CURRENT_TIMESTAMP(6), NULL),
(8, 'Netmoth_BlockIP', 'ip', NULL, NULL, TRUE, CURRENT_TIMESTAMP(6), NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE "provider_signature" RESTART IDENTITY CASCADE;
-- +goose StatementEnd
