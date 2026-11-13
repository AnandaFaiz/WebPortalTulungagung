BEGIN;

DELETE FROM users WHERE role = 'admin_opd';

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('super_admin', 'editor'));

CREATE TABLE IF NOT EXISTS user_permissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_key VARCHAR(100) NOT NULL,
    can_create BOOLEAN NOT NULL DEFAULT FALSE,
    can_read BOOLEAN NOT NULL DEFAULT FALSE,
    can_update BOOLEAN NOT NULL DEFAULT FALSE,
    can_delete BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, permission_key)
);

CREATE INDEX IF NOT EXISTS idx_user_permissions_user_id
    ON user_permissions (user_id);

DELETE FROM role_permissions;

INSERT INTO user_permissions (user_id, permission_key, can_create, can_read, can_update, can_delete)
SELECT users.id, resources.permission_key,
    users.role = 'super_admin',
    users.role = 'super_admin' OR users.role = 'editor',
    users.role = 'super_admin' OR users.role = 'editor',
    users.role = 'super_admin'
FROM users
CROSS JOIN (VALUES
    ('cms.halaman'), ('cms.berita'), ('cms.kategori'), ('cms.pengumuman'),
    ('cms.galeri'), ('cms.sambutan'), ('cms.banner'), ('cms.event'),
    ('cms.layanan'), ('cms.pencarian'), ('pemerintahan.profil'),
    ('pemerintahan.struktur'), ('pemerintahan.pejabat'), ('pemerintahan.opd'),
    ('pemerintahan.tautan'), ('layanan.darurat'), ('layanan.kategori'),
    ('layanan.daftar'), ('layanan.mpp'), ('interaksi.medsos'),
    ('interaksi.tiket'), ('interaksi.kategori'), ('aplikasi'), ('media')
) AS resources(permission_key)
ON CONFLICT (user_id, permission_key) DO NOTHING;

COMMIT;