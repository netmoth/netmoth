-- +goose Up
-- +goose StatementBegin
INSERT INTO "public"."provider_signature" ("id", "name", "type", "source", "signature", "created", "updated") VALUES
(1, 'Abuse.ch_Feodo_BlockIP', 'ip', 'https://feodotracker.abuse.ch/downloads/ipblocklist.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', '2023-03-24 09:27:16.657029', NULL),
(2, 'Blocklist.de_BlockIP', 'ip', 'https://lists.blocklist.de/lists/all.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', '2023-03-24 09:27:16.657029', NULL),
(3, 'EmergingThreats_BlockIP', 'ip', 'https://rules.emergingthreats.net/fwrules/emerging-Block-IPs.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$|\/\d{2})){4})', '2023-03-24 09:27:16.657029', NULL),
(4, 'EmergingThreats_CompromisedIP', 'ip', 'https://rules.emergingthreats.net/blockrules/compromised-ips.txt', '^(?P<ip>(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.|$)){4})', '2023-03-24 09:27:16.657029', NULL),
(5, 'Abuse.ch_Botnet_C2_IP_Denylist', 'botnet', 'https://sslbl.abuse.ch/blacklist/sslipblacklist.csv', '^(?P<date>.{19}),(?P<ip>.*),(?P<port>.*)', '2023-03-24 09:27:16.657029', NULL),
(6, 'Abuse.ch_SSL_Certificate_BlockHash', 'cert', 'https://sslbl.abuse.ch/blacklist/sslblacklist.csv', '^(?P<date>.{19}),(?P<sha1>.*),(?P<name>.*)', '2023-03-24 09:27:16.657029', NULL),
(7, 'CyberCrime_Tracker_BlockUrl', 'tracker', 'http://cybercrime-tracker.net/all.php', '^(?P<url>.{5,})', '2023-03-24 09:27:16.657029', NULL),
(8, 'IoT365_BlockIP', 'ip', NULL, NULL, '2023-03-24 09:27:16.657029', NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE "public"."provider_signature" RESTART IDENTITY CASCADE;
-- +goose StatementEnd
