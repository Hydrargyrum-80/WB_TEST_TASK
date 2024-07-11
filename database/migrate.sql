CREATE TABLE if NOT EXISTS city (
		city_id SERIAL PRIMARY KEY,
		city_name CHARACTER VARYING NOT NULL,
		city_lat FLOAT NOT NULL,
		city_lon FLOAT NOT NULL,
		city_country CHARACTER VARYING NOT NULL,
		UNIQUE(city_name, city_country)
);
CREATE TABLE IF NOT EXISTS weather_predict (
    predict_id SERIAL PRIMARY KEY,
    FK_city_id INTEGER NOT NULL,
    temp FLOAT NOT NULL,
    predict_date TIMESTAMP NOT NULL,
    info JSONB NOT NULL,
    UNIQUE(FK_city_id, predict_date),
    FOREIGN KEY (FK_city_id)
        REFERENCES public.city (city_id)
);
