BEGIN;

INSERT INTO user_permissions (user_id, permission_key, can_create, can_read, can_update, can_delete)
SELECT users.id, resources.permission_key,
    users.role = 'super_admin',
    users.role = 'super_admin' OR users.role = 'editor',
    users.role = 'super_admin' OR users.role = 'editor',
    users.role = 'super_admin'
FROM users
CROSS JOIN (VALUES
    ('pengaturan.identitas'),
    ('pengaturan.kontak'),
    ('pengaturan.menu')
) AS resources(permission_key)
ON CONFLICT (user_id, permission_key) DO NOTHING;

COMMIT;
