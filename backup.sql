--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2
-- Dumped by pg_dump version 15.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plata_currency_rates; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA plata_currency_rates;


ALTER SCHEMA plata_currency_rates OWNER TO postgres;

--
-- Name: rate; Type: TYPE; Schema: plata_currency_rates; Owner: postgres
--

CREATE TYPE plata_currency_rates.rate AS (
	id uuid,
	currency character(3),
	base character(3),
	rate real
);


ALTER TYPE plata_currency_rates.rate OWNER TO postgres;

--
-- Name: add_to_queue(uuid, character, character, numeric); Type: FUNCTION; Schema: plata_currency_rates; Owner: postgres
--

CREATE FUNCTION plata_currency_rates.add_to_queue(_id uuid, _currency character, _base character, _rate numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$

BEGIN

    INSERT INTO plata_currency_rates.rates_queue(id, currency, base, rate, date)

    VALUES (_id, _currency, _base, _rate, current_timestamp);

END;

$$;


ALTER FUNCTION plata_currency_rates.add_to_queue(_id uuid, _currency character, _base character, _rate numeric) OWNER TO postgres;


CREATE FUNCTION plata_currency_rates.get_previous_rate(_currency character, _base character) 
RETURNS TABLE(ret_currency character, ret_base character, ret_rate numeric, ret_date timestamp without time zone) 
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY 
    SELECT currency, base, rate, date 
    FROM plata_currency_rates.rates 
    WHERE currency = _currency AND base = _base 
    ORDER BY date DESC OFFSET 1 LIMIT 1;
END;
$$;



--
-- Name: add_to_rates(uuid, character, character, numeric); Type: FUNCTION; Schema: plata_currency_rates; Owner: postgres
--

CREATE FUNCTION plata_currency_rates.add_to_rates(_id uuid, _currency character, _base character, _rate numeric) RETURNS TABLE(id uuid, currency character, base character, rate numeric, date timestamp without time zone)
    LANGUAGE plpgsql
    AS $$

DECLARE

    inserted_row plata_currency_rates.rates%ROWTYPE;

BEGIN

    INSERT INTO plata_currency_rates.rates(id, currency, base, rate, date)

    VALUES (_id, _currency, _base, _rate, current_timestamp)

    RETURNING * INTO inserted_row;



    RETURN QUERY SELECT inserted_row.*;

END;

$$;


ALTER FUNCTION plata_currency_rates.add_to_rates(_id uuid, _currency character, _base character, _rate numeric) OWNER TO postgres;

--
-- Name: confirm_queue(); Type: FUNCTION; Schema: plata_currency_rates; Owner: postgres
--

CREATE FUNCTION plata_currency_rates.confirm_queue() RETURNS TABLE(ret_id uuid, ret_currency character, ret_base character, ret_rate numeric)
    LANGUAGE plpgsql
    AS $$

DECLARE

    deleted_row plata_currency_rates.rates_queue%ROWTYPE;

BEGIN

    DELETE FROM plata_currency_rates.rates_queue

    WHERE id = (SELECT id FROM plata_currency_rates.rates_queue ORDER BY date ASC LIMIT 1)

    RETURNING * INTO deleted_row;



    IF NOT FOUND THEN

        RAISE NOTICE 'No records found in rates_queue';

        RETURN;

    END IF;



    RETURN QUERY SELECT deleted_row.id, deleted_row.currency, deleted_row.base, deleted_row.rate;

END;

$$;


ALTER FUNCTION plata_currency_rates.confirm_queue() OWNER TO postgres;

--
-- Name: get_by_id(uuid); Type: FUNCTION; Schema: plata_currency_rates; Owner: postgres
--

CREATE FUNCTION plata_currency_rates.get_by_id(_id uuid) RETURNS TABLE(ret_id uuid, ret_currency character, ret_base character, ret_rate numeric, ret_date timestamp without time zone)
    LANGUAGE plpgsql
    AS $$

BEGIN

    RETURN QUERY SELECT id, currency, base, rate, date

                 FROM plata_currency_rates.rates

                 WHERE id = _id;

END;

$$;


ALTER FUNCTION plata_currency_rates.get_by_id(_id uuid) OWNER TO postgres;

--
-- Name: get_last_rate(character, character); Type: FUNCTION; Schema: plata_currency_rates; Owner: postgres
--

CREATE FUNCTION plata_currency_rates.get_last_rate(_currency character, _base character) RETURNS TABLE(ret_currency character, ret_base character, ret_rate numeric, ret_date timestamp without time zone)
    LANGUAGE plpgsql
    AS $$

BEGIN

    RETURN QUERY SELECT currency, base, rate, date

                 FROM plata_currency_rates.rates

                 WHERE currency = _currency and base = _base ORDER BY date DESC LIMIT 1;

END;

$$;


ALTER FUNCTION plata_currency_rates.get_last_rate(_currency character, _base character) OWNER TO postgres;

--
-- Name: add_to_rates(uuid, character, character, real); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_to_rates(_id uuid, _currency character, _base character, _rate real) RETURNS TABLE(id uuid, currency character, base character, rate numeric, date timestamp without time zone)
    LANGUAGE plpgsql
    AS $$

DECLARE

    inserted_row plata_currency_rates.rates%ROWTYPE;

BEGIN

    INSERT INTO plata_currency_rates.rates(id, currency, base, rate, date)

    VALUES (_id, _currency, _base, _rate, current_timestamp)

    RETURNING * INTO inserted_row;



    RETURN QUERY SELECT inserted_row.*;

END;

$$;


ALTER FUNCTION public.add_to_rates(_id uuid, _currency character, _base character, _rate real) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: rates; Type: TABLE; Schema: plata_currency_rates; Owner: postgres
--

CREATE TABLE plata_currency_rates.rates (
    id uuid NOT NULL,
    currency character(3) NOT NULL,
    base character(3) NOT NULL,
    rate numeric NOT NULL,
    date timestamp without time zone NOT NULL
);


ALTER TABLE plata_currency_rates.rates OWNER TO postgres;

--
-- Name: rates_queue; Type: TABLE; Schema: plata_currency_rates; Owner: postgres
--

CREATE TABLE plata_currency_rates.rates_queue (
    id uuid NOT NULL,
    currency character(3) NOT NULL,
    base character(3) NOT NULL,
    rate numeric NOT NULL,
    date timestamp without time zone NOT NULL
);


ALTER TABLE plata_currency_rates.rates_queue OWNER TO postgres;

--
-- Data for Name: rates; Type: TABLE DATA; Schema: plata_currency_rates; Owner: postgres
--

COPY plata_currency_rates.rates (id, currency, base, rate, date) FROM stdin;

\.


--
-- Data for Name: rates_queue; Type: TABLE DATA; Schema: plata_currency_rates; Owner: postgres
--

COPY plata_currency_rates.rates_queue (id, currency, base, rate, date) FROM stdin;
\.


--
-- Name: rates firstkey; Type: CONSTRAINT; Schema: plata_currency_rates; Owner: postgres
--

ALTER TABLE ONLY plata_currency_rates.rates
    ADD CONSTRAINT firstkey PRIMARY KEY (id);


--
-- Name: rates_queue rates_queue_pkey; Type: CONSTRAINT; Schema: plata_currency_rates; Owner: postgres
--

ALTER TABLE ONLY plata_currency_rates.rates_queue
    ADD CONSTRAINT rates_queue_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

