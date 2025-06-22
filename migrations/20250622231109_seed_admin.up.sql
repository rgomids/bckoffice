/* bcrypt hash da senha "admin123"
  gerado com: echo -n 'admin123' | bcrypt
  $2a$10$0uv5XvSLV3zOcnpMj8/RMeM8E81lxZiRvXj9l5fNfuD4.DkJrOe2O
*/

INSERT INTO users (id, email, password_hash, full_name)
VALUES ('01HX0000000000000000000010', 'admin@rcm.tech',
        '$2a$10$0uv5XvSLV3zOcnpMj8/RMeM8E81lxZiRvXj9l5fNfuD4.DkJrOe2O',
        'Admin User');

INSERT INTO user_roles (user_id, role_id)
VALUES ('01HX0000000000000000000010', '01HX0000000000000000000000'); -- role admin