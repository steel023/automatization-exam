create or replace function fn_uuid_time_ordered() returns uuid as $$
declare
    v_time timestamp with time zone:= null;
    v_secs bigint := null;
    v_usec bigint := null;

    v_timestamp bigint := null;
    v_timestamp_hex varchar := null;

    v_clkseq_and_nodeid bigint := null;
    v_clkseq_and_nodeid_hex varchar := null;

    v_bytes bytea;

    c_epoch bigint := -12219292800; -- RFC-4122 epoch: '1582-10-15 00:00:00'
    c_variant bit(64):= x'8000000000000000'; -- RFC-4122 variant: b'10xx...'
begin

    -- Get seconds and micros
    v_time := clock_timestamp();
    v_secs := EXTRACT(EPOCH FROM v_time);
    v_usec := mod(EXTRACT(MICROSECONDS FROM v_time)::numeric, 10^6::numeric);

    -- Generate timestamp hexadecimal (and set version 6)
    v_timestamp := (((v_secs - c_epoch) * 10^6) + v_usec) * 10;
    v_timestamp_hex := lpad(to_hex(v_timestamp), 16, '0');
    v_timestamp_hex := substr(v_timestamp_hex, 2, 12) || '6' || substr(v_timestamp_hex, 14, 3);

    -- Generate clock sequence and node identifier hexadecimal (and set variant b'10xx')
    v_clkseq_and_nodeid := ((random()::numeric * 2^62::numeric)::bigint::bit(64) | c_variant)::bigint;
    v_clkseq_and_nodeid_hex := lpad(to_hex(v_clkseq_and_nodeid), 16, '0');

    -- Concat timestemp, clock sequence and node identifier hexadecimal
    v_bytes := decode(v_timestamp_hex || v_clkseq_and_nodeid_hex, 'hex');

    return encode(v_bytes, 'hex')::uuid;

end $$ language plpgsql;

-- CREATE TABLE incidents (
--     id uuid DEFAULT fn_uuid_time_ordered() PRIMARY KEY,
--     registration_date TIMESTAMP NOT NULL,
--     summary TEXT,
--     incident_type TEXT
-- );

-- INSERT INTO incidents (registration_date, summary, incident_type) VALUES
--     ('2023-01-01 10:00:00', 'Несчастный случай на производстве', 'Несчастный случай'),
--     ('2023-01-02 15:30:00', 'Кража в магазине', 'Кража'),
--     ('2023-01-03 20:00:00', 'Дорожно-транспортное происшествие', 'ДТП'),
--     ('2023-01-04 09:45:00', 'Нарушение общественного порядка', 'Нарушение порядка'),
--     ('2023-01-05 12:30:00', 'Подозрение на мошенничество', 'Мошенничество');
