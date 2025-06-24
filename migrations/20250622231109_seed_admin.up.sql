/* bcrypt hash da senha "admin123"
  gerado com: echo -n 'admin123' | bcrypt
  $2a$10$iu30miJvgCprRhh2tzbQzO2lF2l7OD4ZOrv3oj7LNZ3K2Qt9nwmMq
*/

INSERT INTO users (id, email, password_hash, full_name)
VALUES ('01HX0000000000000000000010', 'admin@rgps.tech',
        '$2a$10$iu30miJvgCprRhh2tzbQzO2lF2l7OD4ZOrv3oj7LNZ3K2Qt9nwmMq',
        'Admin User');

INSERT INTO user_roles (user_id, role_id)
VALUES ('01HX0000000000000000000010', '01HX0000000000000000000000'); -- role admin