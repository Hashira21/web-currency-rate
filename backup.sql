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
6b44eae3-859a-4998-b495-1d0200d13206	EUR	USD	0.91853	2024-01-20 15:42:12.383064
6b44eae3-259a-4998-b495-1d0200d13206	EUR	USD	0.91853	2024-01-20 15:44:06.939855
b56c42bb-c20c-4c91-bf38-b1eda1cb015e	EUR	USD	0.91853	2024-01-20 18:45:51.005755
2b83e7a3-0cbb-4703-81e5-0e61f83c549b	EUR	USD	0.91853	2024-01-20 18:46:06.008022
4fd47ca9-59a5-445c-990e-b6f05d5df570	EUR	USD	0.91853	2024-01-20 18:46:21.010514
d02d349e-bc6b-4fe7-b69b-94aaa0144613	EUR	USD	0.91853	2024-01-20 18:46:36.001156
e7443d4b-f27b-4dca-9492-8cdf80410797	EUR	USD	0.91853	2024-01-20 18:46:51.006264
8e5a0fff-c6d6-4c71-9e32-d13e2fe82e3f	EUR	USD	0.91853	2024-01-20 18:47:06.01387
573405f9-a377-497c-9809-caff00e71bb2	EUR	USD	0.91853	2024-01-20 18:47:21.01072
43aaa931-5bc1-4e4a-9db3-49af55b357c1	EUR	USD	0.91853	2024-01-20 19:01:23.008535
945b1ad4-f88d-4983-a64d-08e25f076a69	EUR	USD	0.91853	2024-01-20 19:05:46.009603
417815e4-6016-4830-9de9-7160753d833b	EUR	USD	0.91853	2024-01-20 19:20:15.013552
f2fb7654-0e17-4e8f-afe2-121c754f5d7c	EUR	USD	0.91853	2024-01-20 19:20:30.010315
da340d00-5af0-4a50-b579-479938e848de	EUR	USD	0.91853	2024-01-20 19:22:21.000607
706a2a69-5976-4e6e-8300-81eddb53573d	EUR	USD	0.91853	2024-01-20 19:22:36.010035
8241552c-7bb9-4965-8fe9-0d484eca5ee8	EUR	USD	0.91853	2024-01-20 19:22:51.009119
560a1acc-4454-4601-b6e0-feda69d0bdef	EUR	USD	0.91853	2024-01-20 19:23:06.005021
5e3f50ec-f546-4b60-a374-f12f46cd2175	EUR	USD	0.91853	2024-01-20 19:23:21.007739
4ac6c779-69a4-45d4-8e5a-1a7a53589ed5	EUR	USD	0.91853	2024-01-20 19:23:36.010801
1fe63796-1a8f-4a01-8f38-fc4f94ed6d89	EUR	USD	0.91853	2024-01-20 19:23:51.009467
cc3604db-c56c-4386-92e2-41b430755e0e	EUR	USD	0.91853	2024-01-20 19:24:06.00132
267172a5-e46f-44e8-a3d0-cfbb14fdfa8f	EUR	USD	0.91853	2024-01-20 19:30:50.006049
00eca0f9-fee1-4573-aca2-4169feeacc73	EUR	USD	0.91853	2024-01-20 19:52:44.005052
f18c048d-0761-4999-bc95-6e3dd63df6e1	EUR	USD	0.91853	2024-01-20 20:19:30.012102
5af1a030-d342-4467-aa55-f2528fccf493	EUR	USD	0.91853	2024-01-20 20:46:29.010515
e0dc014a-704b-4faa-8a4f-d42acb48f049	EUR	USD	0.91853	2024-01-20 21:25:05.013845
a432c2d0-8a8c-4e58-83f0-3dd93d83ddf2	EUR	JPY	0.0062	2024-01-20 21:30:44.0094
37edaf79-e7a1-4828-aa61-a48fa67eba2f	EUR	USD	0.91853	2024-01-21 04:34:08.001187
4d68d255-3d57-4328-ab8b-44309de3ccb2	MXN	USD	17.13	2024-01-21 04:34:23.002119
1b7c7b42-a2f6-4edb-93c2-cffa2ef4cbcf	EUR	JPY	0.0062	2024-01-21 14:07:22.015196
4910c3d3-b785-4bec-9b4b-91639fa63ef9	EUR	MXN	0.05362	2024-01-21 14:49:52.006875
f37324a3-b7c3-4845-8115-3517b276d154	MXN	EUR	18.5897	2024-01-23 14:02:07.011773
e91b64c6-bf46-4858-b7f9-eb707fbaa015	MXN	EUR	18.5897	2024-01-23 14:02:22.009796
a9f16f38-a22f-42c3-a215-4061d831a098	MXN	EUR	18.58970069885254	2024-01-23 14:02:37.015888
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

