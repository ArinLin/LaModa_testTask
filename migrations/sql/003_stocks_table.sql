CREATE TABLE stocks (
    good_id         INT REFERENCES goods (id) NOT NULL,
    warehouse_id    INT REFERENCES warehouses (id) NOT NULL,
    amount          INT NOT NULL,
    reserved        INT NOT NULL DEFAULT 0,
    PRIMARY KEY(good_id, warehouse_id)
);
CREATE INDEX idx_stocks_good_id ON stocks (good_id);
CREATE INDEX idx_stocks_warehouse_id ON stocks (warehouse_id);

---- create above / drop below ----
DROP INDEX idx_stocks_warehouse_id;
DROP INDEX idx_stocks_good_id;
DROP TABLE stocks;