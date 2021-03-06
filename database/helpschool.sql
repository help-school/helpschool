--
-- PostgreSQL database dump
--

-- Dumped from database version 13.2
-- Dumped by pg_dump version 13.2

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
-- Name: helpschool; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA helpschool;


ALTER SCHEMA helpschool OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: countries; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.countries (
    country_id uuid NOT NULL,
    name character varying(255) NOT NULL COLLATE pg_catalog."en-US-x-icu",
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.countries OWNER TO postgres;

--
-- Name: districts; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.districts (
    district_id uuid NOT NULL,
    state_id uuid NOT NULL,
    name character varying(512) NOT NULL,
    govt_id character varying(1024),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.districts OWNER TO postgres;

--
-- Name: school_supplies; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.school_supplies (
    school_id uuid NOT NULL,
    supply_id uuid NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    fulfilled_count integer DEFAULT 0,
    extra_info jsonb
);


ALTER TABLE helpschool.school_supplies OWNER TO postgres;

--
-- Name: schools; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.schools (
    school_id uuid NOT NULL,
    district_id uuid NOT NULL,
    name character varying(512) NOT NULL,
    address character varying(4096),
    govt_id character varying(2048),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    place character varying(512)
);


ALTER TABLE helpschool.schools OWNER TO postgres;

--
-- Name: states; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.states (
    state_id uuid NOT NULL,
    name character varying(512) NOT NULL COLLATE pg_catalog."en_US.UTF-8",
    country_id uuid NOT NULL,
    govt_id character varying(1024),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.states OWNER TO postgres;

--
-- Name: supplies; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.supplies (
    supply_id uuid NOT NULL,
    title character varying(1024) NOT NULL,
    url character varying(4096) NOT NULL,
    description character varying(4096),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    country_id uuid NOT NULL
);


ALTER TABLE helpschool.supplies OWNER TO postgres;

--
-- Name: teacher_requests; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.teacher_requests (
    id bigint NOT NULL,
    teacher_name character varying(256) NOT NULL,
    teacher_phone character varying(32),
    url character varying(4096) NOT NULL,
    quantity_needed integer DEFAULT 0 NOT NULL,
    address character varying(4096) NOT NULL,
    place character varying(256),
    district character varying(256) NOT NULL,
    state character varying(256) NOT NULL,
    country character varying(256),
    zipcode character varying(128) DEFAULT 0,
    extra_info jsonb,
    photo_link character varying(4096),
    teacher_email character varying(256)
);


ALTER TABLE helpschool.teacher_requests OWNER TO postgres;

--
-- Name: TABLE teacher_requests; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON TABLE helpschool.teacher_requests IS 'teachers requests will be saved here ';


--
-- Name: teacher_requests_id_seq; Type: SEQUENCE; Schema: helpschool; Owner: postgres
--

ALTER TABLE helpschool.teacher_requests ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME helpschool.teacher_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users_donations; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.users_donations (
    user_email character varying(256) NOT NULL,
    user_id character varying(128),
    user_name character varying(256),
    school_id uuid NOT NULL,
    supply_id uuid NOT NULL,
    quantity integer,
    status character varying(128) DEFAULT 'Ordered'::character varying NOT NULL,
    tracking_url character varying(1024),
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    extra_info jsonb
);


ALTER TABLE helpschool.users_donations OWNER TO postgres;

--
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (country_id);


--
-- Name: countries country_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.countries
    ADD CONSTRAINT country_id UNIQUE (country_id);


--
-- Name: districts district_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT district_pkey PRIMARY KEY (district_id);


--
-- Name: states name_country_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT name_country_id UNIQUE (country_id, name);


--
-- Name: schools name_place_address_district_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT name_place_address_district_id UNIQUE (district_id, name, place, address);


--
-- Name: districts name_state_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT name_state_id UNIQUE (state_id, name);


--
-- Name: schools school_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT school_pkey PRIMARY KEY (school_id);


--
-- Name: school_supplies school_supplies_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT school_supplies_pkey PRIMARY KEY (school_id, supply_id);


--
-- Name: states states_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT states_pkey PRIMARY KEY (state_id);


--
-- Name: supplies supplies_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.supplies
    ADD CONSTRAINT supplies_pkey PRIMARY KEY (supply_id);


--
-- Name: users_donations users_donations_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT users_donations_pkey PRIMARY KEY (user_email);


--
-- Name: fki_countries_province; Type: INDEX; Schema: helpschool; Owner: postgres
--

CREATE INDEX fki_countries_province ON helpschool.states USING btree (country_id);


--
-- Name: supplies country_country_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.supplies
    ADD CONSTRAINT country_country_id FOREIGN KEY (country_id) REFERENCES helpschool.countries(country_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: states country_state; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT country_state FOREIGN KEY (country_id) REFERENCES helpschool.countries(country_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: schools district_district_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT district_district_id FOREIGN KEY (district_id) REFERENCES helpschool.districts(district_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: school_supplies schools_school_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT schools_school_id FOREIGN KEY (school_id) REFERENCES helpschool.schools(school_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: users_donations schools_school_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT schools_school_id FOREIGN KEY (school_id) REFERENCES helpschool.schools(school_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: districts state_state_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT state_state_id FOREIGN KEY (state_id) REFERENCES helpschool.states(state_id) NOT VALID;


--
-- Name: school_supplies supplies_supply_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT supplies_supply_id FOREIGN KEY (supply_id) REFERENCES helpschool.supplies(supply_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: users_donations supplies_supply_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT supplies_supply_id FOREIGN KEY (supply_id) REFERENCES helpschool.supplies(supply_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 13.2
-- Dumped by pg_dump version 13.2

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
-- Name: helpschool; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA helpschool;


ALTER SCHEMA helpschool OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: countries; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.countries (
    country_id uuid NOT NULL,
    name character varying(255) NOT NULL COLLATE pg_catalog."en-US-x-icu",
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.countries OWNER TO postgres;

--
-- Name: districts; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.districts (
    district_id uuid NOT NULL,
    state_id uuid NOT NULL,
    name character varying(512) NOT NULL,
    govt_id character varying(1024),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.districts OWNER TO postgres;

--
-- Name: school_supplies; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.school_supplies (
    school_id uuid NOT NULL,
    supply_id uuid NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    fulfilled_count integer DEFAULT 0,
    extra_info jsonb
);


ALTER TABLE helpschool.school_supplies OWNER TO postgres;

--
-- Name: schools; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.schools (
    school_id uuid NOT NULL,
    district_id uuid NOT NULL,
    name character varying(512) NOT NULL,
    address character varying(4096),
    govt_id character varying(2048),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    place character varying(512)
);


ALTER TABLE helpschool.schools OWNER TO postgres;

--
-- Name: states; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.states (
    state_id uuid NOT NULL,
    name character varying(512) NOT NULL COLLATE pg_catalog."en_US.UTF-8",
    country_id uuid NOT NULL,
    govt_id character varying(1024),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now()
);


ALTER TABLE helpschool.states OWNER TO postgres;

--
-- Name: supplies; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.supplies (
    supply_id uuid NOT NULL,
    title character varying(1024) NOT NULL,
    url character varying(4096) NOT NULL,
    description character varying(4096),
    extra_info jsonb,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    country_id uuid NOT NULL
);


ALTER TABLE helpschool.supplies OWNER TO postgres;

--
-- Name: teacher_requests; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.teacher_requests (
    id bigint NOT NULL,
    teacher_name character varying(256) NOT NULL,
    teacher_phone character varying(32),
    url character varying(4096) NOT NULL,
    quantity_needed integer DEFAULT 0 NOT NULL,
    address character varying(4096) NOT NULL,
    place character varying(256),
    district character varying(256) NOT NULL,
    state character varying(256) NOT NULL,
    country character varying(256),
    zipcode character varying(128) DEFAULT 0,
    extra_info jsonb,
    photo_link character varying(4096),
    teacher_email character varying(256)
);


ALTER TABLE helpschool.teacher_requests OWNER TO postgres;

--
-- Name: TABLE teacher_requests; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON TABLE helpschool.teacher_requests IS 'teachers requests will be saved here ';


--
-- Name: teacher_requests_id_seq; Type: SEQUENCE; Schema: helpschool; Owner: postgres
--

ALTER TABLE helpschool.teacher_requests ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME helpschool.teacher_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.users (
    id uuid NOT NULL,
    user_email character varying(256) NOT NULL,
    user_id character varying(128),
    user_name character varying(256),
    extra_info jsonb,
    created_date time with time zone DEFAULT now(),
    modified_date timestamp with time zone
);


ALTER TABLE helpschool.users OWNER TO postgres;

--
-- Name: COLUMN users.id; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON COLUMN helpschool.users.id IS 'This will be generated by postgres when we first insert and will be used in APIs with front end , other column data will come from auth0 or anyother services, since other column data is PII , we will use this ID for URL''s ';


--
-- Name: COLUMN users.user_email; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON COLUMN helpschool.users.user_email IS 'comes from auth0 initially';


--
-- Name: COLUMN users.user_id; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON COLUMN helpschool.users.user_id IS 'comes from auth0 initially';


--
-- Name: COLUMN users.user_name; Type: COMMENT; Schema: helpschool; Owner: postgres
--

COMMENT ON COLUMN helpschool.users.user_name IS 'comes from auth0 initially';


--
-- Name: users_donations; Type: TABLE; Schema: helpschool; Owner: postgres
--

CREATE TABLE helpschool.users_donations (
    user_id uuid NOT NULL,
    school_id uuid NOT NULL,
    supply_id uuid NOT NULL,
    quantity integer,
    status character varying(128) DEFAULT 'Ordered'::character varying NOT NULL,
    tracking_url character varying(1024),
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    modified_date timestamp with time zone DEFAULT now(),
    extra_info jsonb
);


ALTER TABLE helpschool.users_donations OWNER TO postgres;

--
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (country_id);


--
-- Name: countries country_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.countries
    ADD CONSTRAINT country_id UNIQUE (country_id);


--
-- Name: districts district_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT district_pkey PRIMARY KEY (district_id);


--
-- Name: states name_country_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT name_country_id UNIQUE (country_id, name);


--
-- Name: schools name_place_address_district_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT name_place_address_district_id UNIQUE (district_id, name, place, address);


--
-- Name: districts name_state_id; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT name_state_id UNIQUE (state_id, name);


--
-- Name: schools school_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT school_pkey PRIMARY KEY (school_id);


--
-- Name: school_supplies school_supplies_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT school_supplies_pkey PRIMARY KEY (school_id, supply_id);


--
-- Name: states states_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT states_pkey PRIMARY KEY (state_id);


--
-- Name: supplies supplies_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.supplies
    ADD CONSTRAINT supplies_pkey PRIMARY KEY (supply_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: fki_countries_province; Type: INDEX; Schema: helpschool; Owner: postgres
--

CREATE INDEX fki_countries_province ON helpschool.states USING btree (country_id);


--
-- Name: fki_user_user_id; Type: INDEX; Schema: helpschool; Owner: postgres
--

CREATE INDEX fki_user_user_id ON helpschool.users_donations USING btree (user_id);


--
-- Name: user_email_idx; Type: INDEX; Schema: helpschool; Owner: postgres
--

CREATE INDEX user_email_idx ON helpschool.users USING btree (user_email) INCLUDE (user_email);


--
-- Name: supplies country_country_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.supplies
    ADD CONSTRAINT country_country_id FOREIGN KEY (country_id) REFERENCES helpschool.countries(country_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: states country_state; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.states
    ADD CONSTRAINT country_state FOREIGN KEY (country_id) REFERENCES helpschool.countries(country_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: schools district_district_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.schools
    ADD CONSTRAINT district_district_id FOREIGN KEY (district_id) REFERENCES helpschool.districts(district_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: school_supplies schools_school_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT schools_school_id FOREIGN KEY (school_id) REFERENCES helpschool.schools(school_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: users_donations schools_school_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT schools_school_id FOREIGN KEY (school_id) REFERENCES helpschool.schools(school_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: districts state_state_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.districts
    ADD CONSTRAINT state_state_id FOREIGN KEY (state_id) REFERENCES helpschool.states(state_id) NOT VALID;


--
-- Name: school_supplies supplies_supply_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.school_supplies
    ADD CONSTRAINT supplies_supply_id FOREIGN KEY (supply_id) REFERENCES helpschool.supplies(supply_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: users_donations supplies_supply_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT supplies_supply_id FOREIGN KEY (supply_id) REFERENCES helpschool.supplies(supply_id) ON UPDATE SET NULL ON DELETE SET NULL NOT VALID;


--
-- Name: users_donations users_id; Type: FK CONSTRAINT; Schema: helpschool; Owner: postgres
--

ALTER TABLE ONLY helpschool.users_donations
    ADD CONSTRAINT users_id FOREIGN KEY (user_id) REFERENCES helpschool.users(id) NOT VALID;


--
-- PostgreSQL database dump complete
--

