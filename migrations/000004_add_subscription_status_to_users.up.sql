-- subscription_status: 'ACTIVE' | 'INACTIVE'
ALTER TABLE "users" ADD COLUMN subscription_status TEXT NOT NULL DEFAULT 'INACTIVE';