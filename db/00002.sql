BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'surf' 
        AND table_name = 'records' 
        AND column_name = 'created_at'
    ) THEN
        ALTER TABLE surf.records ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
    END IF;
END $$;

COMMIT;

