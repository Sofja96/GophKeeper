create table if not exists users
(
      id bigserial primary key,
      username varchar unique not null,
      password varchar not null,
      created_at timestamp with time zone default now() not null,
      updated_at timestamp with time zone default now() not null
);

create table if not exists data (
    id bigserial  primary key,
    user_id bigint not null references users(id) on delete cascade,
    data_type varchar not null,        -- Пароль, бинарные данные и т.д.
    data_content BYTEA not null,       -- Храним зашифрованные данные
    metadata JSONB,                    -- Метаданные (опционально)
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);