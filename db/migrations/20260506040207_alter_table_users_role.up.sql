ALTER TABLE users
    DROP COLUMN role;

CREATE TABLE user_roles (
    user_id VARCHAR(36)  NOT NULL,
    role_id VARCHAR(100) NOT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);
