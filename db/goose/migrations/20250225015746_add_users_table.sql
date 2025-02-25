-- +goose Up
-- +goose StatementBegin
CREATE TABLE goose_post (
  id int NOT NULL,
  title text,
  body text,
  PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE goose_post;
-- +goose StatementEnd

