-- Удаляем таблицы в обратном порядке их создания
DROP TABLE IF EXISTS incident_individual_link;
DROP TABLE IF EXISTS individuals;
DROP TABLE IF EXISTS decisions;
DROP TABLE IF EXISTS incidents;

-- Удаляем функцию
DROP FUNCTION IF EXISTS fn_uuid_time_ordered;
