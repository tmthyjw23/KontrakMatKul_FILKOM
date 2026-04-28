-- Migration 004: Add Contract Periods Table
CREATE TABLE IF NOT EXISTS contract_periods (
    id          INT          NOT NULL AUTO_INCREMENT,
    is_open     TINYINT(1)   NOT NULL DEFAULT 1,
    start_date  DATETIME     NOT NULL,
    end_date    DATETIME     NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed an initial contract period if none exists
INSERT INTO contract_periods (id, is_open, start_date, end_date)
SELECT 1, 1, '2026-01-01 00:00:00', '2026-12-31 23:59:59'
FROM (SELECT 1) AS tmp
WHERE NOT EXISTS (SELECT 1 FROM contract_periods WHERE id = 1);
