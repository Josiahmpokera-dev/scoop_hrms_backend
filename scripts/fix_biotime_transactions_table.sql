-- Fix biotime_transactions table column name
-- Run this if the table was created with the wrong column name

-- Check if column exists with wrong name and rename it
DO $$
BEGIN
    -- If bio_time_transaction_id exists, rename it to biotime_transaction_id
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'biotime_transactions' 
        AND column_name = 'bio_time_transaction_id'
    ) THEN
        ALTER TABLE biotime_transactions 
        RENAME COLUMN bio_time_transaction_id TO biotime_transaction_id;
        RAISE NOTICE 'Renamed column bio_time_transaction_id to biotime_transaction_id';
    END IF;
END $$;

-- Or if the table doesn't exist at all, it will be created by GORM on next migration
-- Just restart your application and GORM will handle it
