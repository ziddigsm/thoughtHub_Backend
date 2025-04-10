
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_on = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS MENU (
    id SERIAL PRIMARY KEY,
    options VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_navbar BOOLEAN DEFAULT FALSE,
    created_on TIMESTAMP DEFAULT NOW(),
    updated_on TIMESTAMP DEFAULT NOW()
);


DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp' AND tgrelid = 'menu'::regclass
    ) THEN
        CREATE TRIGGER set_timestamp
        BEFORE UPDATE ON MENU
        FOR EACH ROW
        EXECUTE FUNCTION update_timestamp();
    END IF;
END $$;
