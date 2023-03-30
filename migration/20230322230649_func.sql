-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = CURRENT_TIMESTAMP(6);
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_updated
BEFORE UPDATE ON "provider_signature"
FOR EACH ROW EXECUTE PROCEDURE update_updated();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_updated name ON "provider_signature"
DROP FUNCTION IF EXISTS "date_update";
-- +goose StatementEnd
