CREATE TABLE IF NOT EXISTS objects (
    -- Primary Key
    id VARCHAR(36) PRIMARY KEY,

    -- Basic Required Fields
    description TEXT,
    disk VARCHAR(255),
    location VARCHAR(255),

    -- Optional Fields (NULLable due to *string pointers in Go)
    bucket VARCHAR(255) NULL,
    region VARCHAR(255) NULL,
    parent VARCHAR(36) NULL,
    encryption_method VARCHAR(50) NULL,

    -- Booleans (Using BOOLEAN or TINYINT(1) depending on DB)
    encrypted BOOLEAN DEFAULT FALSE,
    public BOOLEAN DEFAULT FALSE,
    deleted BOOLEAN DEFAULT FALSE,

    -- File Metadata
    filename VARCHAR(255) NOT NULL,
    extension VARCHAR(50),
    hash VARCHAR(64) NOT NULL, -- Assuming SHA-256 hash length
    size BIGINT, -- Use BIGINT for file size (int64 in Go)
    unit VARCHAR(10),

    -- Timestamps
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,

    -- Foreign Key/Owner
    user_owner VARCHAR(36) NULL

    -- Unique Constraint for Hash (Enforcing uniqueness only for active/undeleted files)
);

-- Separate Index Definitions (For better readability)
-- CREATE INDEX IF NOT EXISTS idx_objects_parent ON objects (parent);
-- CREATE INDEX IF NOT EXISTS idx_objects_public ON objects (public);
-- CREATE INDEX IF NOT EXISTS idx_objects_deleted ON objects (deleted);
-- CREATE INDEX IF NOT EXISTS idx_objects_user_owner ON objects (user_owner);