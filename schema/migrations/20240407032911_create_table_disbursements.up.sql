CREATE TABLE IF NOT EXISTS "disbursements" (
    "id" SERIAL PRIMARY KEY,
    "amount" BIGINT NOT NULL,
    "source_bank_code" VARCHAR(50) NOT NULL,
    "source_account" VARCHAR(100) NOT NULL,
    "destination_bank_code" VARCHAR(50) NOT NULL,
    "destination_account" VARCHAR(100) NOT NULL,
    "reference_id" VARCHAR(100) NOT NULL,
    "remarks" TEXT,
    "status" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);