-- Put your ddl queries here.
-- NO NEED TO PUT CREATE DATABASE statement here (assuming we already create the database when starting postgres docker container)

CREATE TABLE provinces
(
    id         BIGINT PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE cities
(
    id          BIGINT PRIMARY KEY        NOT NULL,
    name        VARCHAR                   NOT NULL,
    province_id BIGINT REFERENCES provinces (id),
    created_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE user_roles
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR                   NOT NULL UNIQUE,
    password     VARCHAR                   NOT NULL,
    user_role_id BIGINT                    NOT NULL REFERENCES user_roles (id),
    is_verified  BOOLEAN                   NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE doctor_specializations
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    image      VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE doctor_profiles
(
    user_id                  BIGSERIAL PRIMARY KEY REFERENCES users (id),
    name                     VARCHAR                   NOT NULL,
    profile_photo            VARCHAR                   NOT NULL,
    starting_year            INTEGER                   NOT NULL,
    doctor_certificate       VARCHAR                   NOT NULL,
    doctor_specialization_id BIGINT                    NOT NULL REFERENCES doctor_specializations (id),
    consultation_fee         NUMERIC                   NOT NULL,
    is_online                BOOL                      NOT NULL,
    created_at               TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at               TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at               TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE user_profiles
(
    user_id       BIGSERIAL PRIMARY KEY REFERENCES users (id),
    name          VARCHAR                   NOT NULL,
    profile_photo VARCHAR                   NOT NULL,
    date_of_birth timestamptz               NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE addresses
(
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR                   NOT NULL,
    address      TEXT                      NOT NULL,
    sub_district VARCHAR                   NOT NULL,
    district     VARCHAR                   NOT NULL,
    city         BIGINT                    NOT NULL REFERENCES cities (id),
    province     BIGINT                    NOT NULL REFERENCES provinces (id),
    postal_code  VARCHAR                   NOT NULL,
    latitude     VARCHAR                   NOT NULL,
    longitude    VARCHAR                   NOT NULL,
    status       INTEGER                   NOT NULL,
    profile_id   BIGINT                    NOT NULL REFERENCES user_profiles (user_id),
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE shipping_methods
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE pharmacies
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  VARCHAR                   NOT NULL,
    address               TEXT                      NOT NULL,
    sub_district          VARCHAR                   NOT NULL,
    district              VARCHAR                   NOT NULL,
    city                  BIGINT                    NOT NULL REFERENCES cities (id),
    province              BIGINT                    NOT NULL REFERENCES provinces (id),
    postal_code           VARCHAR                   NOT NULL,
    latitude              VARCHAR                   NOT NULL,
    longitude             VARCHAR                   NOT NULL,
    pharmacist_name       VARCHAR                   NOT NULL,
    pharmacist_license_no VARCHAR                   NOT NULL,
    pharmacist_phone_no   VARCHAR                   NOT NULL,
    operational_hours     VARCHAR                   NOT NULL,
    operational_days      VARCHAR                   NOT NULL,
    pharmacy_admin_id     BIGINT                    NOT NULL REFERENCES users (id),
    created_at            TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at            TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at            TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE pharmacy_shipping_methods
(
    id                 BIGSERIAL PRIMARY KEY,
    pharmacy_id        BIGINT                    NOT NULL REFERENCES pharmacies (id),
    shipping_method_id BIGINT                    NOT NULL REFERENCES shipping_methods (id),
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE verification_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    token      VARCHAR                   NOT NULL,
    is_valid   BOOLEAN     DEFAULT TRUE  NOT NULL,
    expired_at TIMESTAMPTZ               NOT NULL,
    email      VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE forgot_password_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    token      VARCHAR                   NOT NULL,
    is_valid   BOOLEAN     DEFAULT TRUE  NOT NULL,
    expired_at TIMESTAMPTZ               NOT NULL,
    user_id    BIGINT                    NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE manufacturers
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    image      VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_categories
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR UNIQUE            NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE drug_classifications
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE products
(
    id                     BIGSERIAL PRIMARY KEY,
    name                   VARCHAR                   NOT NULL,
    generic_name           VARCHAR                   NOT NULL,
    content                VARCHAR                   NOT NULL,
    manufacturer_id        BIGINT                    NOT NULL REFERENCES manufacturers (id),
    description            TEXT                      NOT NULL,
    drug_classification_id BIGINT                    NOT NULL REFERENCES drug_classifications (id),
    product_category_id    BIGINT                    NOT NULL REFERENCES product_categories (id),
    drug_form              VARCHAR                   NOT NULL,
    unit_in_pack           VARCHAR                   NOT NULL,
    selling_unit           VARCHAR                   NOT NULL,
    weight                 FLOAT                     NOT NULL, -- in gram
    length                 FLOAT                     NOT NULL, -- in cm
    width                  FLOAT                     NOT NULL, -- in cm
    height                 FLOAT                     NOT NULL, -- in cm
    image                  VARCHAR                   NOT NULL,
    UNIQUE (name, generic_name, content, manufacturer_id),
    created_at             TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at             TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at             TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE cart_items
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT                    NOT NULL REFERENCES users (id),
    product_id BIGINT                    NOT NULL REFERENCES products (id),
    quantity   INT                       NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE pharmacy_products
(
    id          BIGSERIAL PRIMARY KEY,
    pharmacy_id BIGINT                    NOT NULL REFERENCES pharmacies (id),
    product_id  BIGINT                    NOT NULL REFERENCES products (id),
    is_active   BOOL                      NOT NULL,
    price       NUMERIC                   NOT NULL,
    stock       INT                       NOT NULL,
    UNIQUE (pharmacy_id, product_id),
    created_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE product_stock_mutation_types
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_stock_mutation_request_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_stock_mutations
(
    id                             BIGSERIAL PRIMARY KEY,
    pharmacy_product_id            BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    product_stock_mutation_type_id BIGINT                    NOT NULL REFERENCES product_stock_mutation_types (id),
    stock                          INT                       NOT NULL,
    created_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                     TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE consultation_session_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE consultation_sessions
(
    id                             BIGSERIAL PRIMARY KEY,
    user_id                        BIGINT                    NOT NULL REFERENCES user_profiles (user_id),
    doctor_id                      BIGINT                    NOT NULL REFERENCES doctor_profiles (user_id),
    consultation_session_status_id BIGINT                    NOT NULL REFERENCES consultation_session_statuses (id),
    created_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                     TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE consultation_messages
(
    id           BIGSERIAL PRIMARY KEY,
    session_id   BIGINT                    NOT NULL REFERENCES consultation_sessions (id),
    sender_id    BIGINT                    NOT NULL REFERENCES users (id),
    message_type INTEGER                   NOT NULL,
    message      VARCHAR                   NOT NULL,
    attachment   VARCHAR                   NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE prescriptions
(
    id         BIGSERIAL PRIMARY KEY,
    session_id BIGINT UNIQUE             NOT NULL REFERENCES consultation_sessions (id),
    symptoms   VARCHAR                   NOT NULL,
    diagnosis  VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE prescription_products
(
    id              BIGSERIAL PRIMARY KEY,
    prescription_id BIGINT                    NOT NULL REFERENCES prescriptions (id),
    product_id      BIGINT                    NOT NULL REFERENCES products (id),
    note            VARCHAR                   NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE sick_leave_forms
(
    id            BIGSERIAL PRIMARY KEY,
    session_id    BIGINT UNIQUE             NOT NULL REFERENCES consultation_sessions (id),
    starting_date TIMESTAMPTZ               NOT NULL,
    ending_date   TIMESTAMPTZ               NOT NULL,
    description   TEXT                      NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE order_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE payment_methods
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE transaction_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE transactions
(
    id         BIGSERIAL PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL ,
    payment_proof VARCHAR NOT NULL ,
    transaction_status_id  BIGINT                    NOT NULL REFERENCES transaction_statuses (id),
    payment_method_id  BIGINT                    NOT NULL REFERENCES payment_methods (id),
    address       VARCHAR                   NOT NULL ,
    user_id            BIGINT                    NOT NULL REFERENCES users (id),
    total_payment NUMERIC NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE orders
(
    id                 BIGSERIAL PRIMARY KEY,
    date               TIMESTAMPTZ               NOT NULL,
    pharmacy_id        BIGINT                    NOT NULL REFERENCES pharmacies (id),
    no_of_items        INTEGER                   NOT NULL,
    pharmacy_address   VARCHAR                   NOT NULL,
    shipping_method_id BIGINT                    NOT NULL REFERENCES shipping_methods (id),
    shipping_cost      NUMERIC                   NOT NULL,
    total_payment      NUMERIC                   NOT NULL,
    transaction_id  BIGINT                    NOT NULL REFERENCES transactions (id),
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE order_status_logs
(
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT                    NOT NULL REFERENCES orders (id),
    order_status_id BIGINT                    NOT NULL REFERENCES order_statuses (id),
    is_latest BOOL NOT NULL ,
    description TEXT NOT NULL ,
    created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE order_details
(
    id                                BIGSERIAL PRIMARY KEY,
    order_id                          BIGINT                    NOT NULL REFERENCES orders (id),
    product_id                        BIGINT                    NOT NULL REFERENCES products (id),
    quantity                          INTEGER                   NOT NULL,
    name                              VARCHAR                   NOT NULL,
    generic_name                      VARCHAR                   NOT NULL,
    content                           VARCHAR                   NOT NULL,
    description                       VARCHAR                   NOT NULL,
    image                             VARCHAR                   NOT NULL,
    price                             NUMERIC                   NOT NULL,
    created_at                        TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                        TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                        TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_stock_mutation_requests
(
    id                                       BIGSERIAL PRIMARY KEY,
    pharmacy_product_origin_id               BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    pharmacy_product_dest_id                 BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    stock                                    INT                       NOT NULL,
    product_stock_mutation_request_status_id BIGINT                    NOT NULL REFERENCES product_stock_mutation_request_statuses (id),
    order_detail_id                          BIGINT                    NULL REFERENCES order_details (id),
    created_at                               TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                               TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                               TIMESTAMPTZ DEFAULT NULL
);


-- INSERT DATA
INSERT INTO user_roles (name)
values ('Admin'),
       ('Pharmacy Admin'),
       ('Doctor'),
       ('User');

INSERT INTO users (email, password, user_role_id, is_verified)
VALUES ('byebyesick@gmail.com', '$2a$04$MYf2/GkfNPUUZUj8zInF.ej7KqSVO3KlJrbNEwkCtCerFXzqbOsDe', 1, true),
       ('yafi.tamfan08@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('tifan@email.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('random@email.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('wasikamin4@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 3, true),
       ('lumbanraja.boy@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 4, true);

INSERT INTO doctor_specializations (name, image)
values ('General Practitioners',
        'https://byebyesick-bucket.irfancen.com/doctor_specializations/doctor-specs.jpg'),
       ('Pediatric Specialist',
        'https://byebyesick-bucket.irfancen.com/doctor_specializations/doctor-specs.jpg');

INSERT INTO doctor_profiles (user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id,
                             consultation_fee, is_online)
VALUES (5, 'dr. Wasik Amin', '', 2021, '', 1, 10, false);

INSERT INTO user_profiles(user_id, name, profile_photo, date_of_birth)
VALUES (6, 'Benedict Boy', '', '2000-11-25');

INSERT INTO manufacturers (name, image)
values ('Soho Industri Pharmasi', 'https://byebyesick-bucket.irfancen.com/doctor_specializations/soho.png'),
       ('Amarox Pharma Global', 'https://byebyesick-bucket.irfancen.com/doctor_specializations/amarox.jpeg');

INSERT INTO shipping_methods (name)
values ('Official Instant'),
       ('Official Same Day'),
       ('Non Official');

INSERT INTO product_categories (name)
values ('Obat'),
       ('Non Obat');

INSERT INTO drug_classifications (name)
values ('Obat Bebas'),
       ('Obat Keras'),
       ('Obat Bebas Terbatas'),
       ('Non Obat');

INSERT INTO product_stock_mutation_types (name)
values ('Addition'),
       ('Deduction');

INSERT INTO product_stock_mutation_request_statuses (name)
values ('Pending'),
       ('Accepted'),
       ('Rejected');

INSERT INTO payment_methods (name)
values ('Bank Transfer');

INSERT INTO order_statuses (name)
values ('Waiting for Pharmacy'),
       ('Processed'),
       ('Sent'),
       ('Order Confirmed'),
        ('Canceled by Pharmacy'),
        ('Canceled by User');


INSERT INTO transaction_statuses (name)
values ('Unpaid'),
       ('Waiting for Confirmation'),
       ('Payment Rejected'),
       ('Paid'),
       ('Canceled');

INSERT INTO provinces (id, name)
VALUES (1, 'Bali'),
       (2, 'Bangka Belitung'),
       (3, 'Banten'),
       (4, 'Bengkulu'),
       (5, 'DI Yogyakarta'),
       (6, 'DKI Jakarta'),
       (7, 'Gorontalo'),
       (8, 'Jambi'),
       (9, 'Jawa Barat'),
       (10, 'Jawa Tengah'),
       (11, 'Jawa Timur'),
       (12, 'Kalimantan Barat'),
       (13, 'Kalimantan Selatan'),
       (14, 'Kalimantan Tengah'),
       (15, 'Kalimantan Timur'),
       (16, 'Kalimantan Utara'),
       (17, 'Kepulauan Riau'),
       (18, 'Lampung'),
       (19, 'Maluku'),
       (20, 'Maluku Utara'),
       (21, 'Nanggroe Aceh Darussalam (NAD)'),
       (22, 'Nusa Tenggara Barat (NTB)'),
       (23, 'Nusa Tenggara Timur (NTT)'),
       (24, 'Papua'),
       (25, 'Papua Barat'),
       (26, 'Riau'),
       (27, 'Sulawesi Barat'),
       (28, 'Sulawesi Selatan'),
       (29, 'Sulawesi Tengah'),
       (30, 'Sulawesi Tenggara'),
       (31, 'Sulawesi Utara'),
       (32, 'Sumatera Barat'),
       (33, 'Sumatera Selatan'),
       (34, 'Sumatera Utara');

INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Badung', 17, 1),
       ('Kab. Bangli', 32, 1),
       ('Kab. Buleleng', 94, 1),
       ('Kota Denpasar', 114, 1),
       ('Kab. Gianyar', 128, 1),
       ('Kab. Jembrana', 161, 1),
       ('Kab. Karangasem', 170, 1),
       ('Kab. Klungkung', 197, 1),
       ('Kab. Tabanan', 447, 1),
       ('Kab. Bangka', 27, 2),
       ('Kab. Bangka Barat', 28, 2),
       ('Kab. Bangka Selatan', 29, 2),
       ('Kab. Bangka Tengah', 30, 2),
       ('Kab. Belitung', 56, 2),
       ('Kab. Belitung Timur', 57, 2),
       ('Kota Pangkal Pinang', 334, 2),
       ('Kota Cilegon', 106, 3),
       ('Kab. Lebak', 232, 3),
       ('Kab. Pandeglang', 331, 3),
       ('Kab. Serang', 402, 3),
       ('Kota Serang', 403, 3),
       ('Kab. Tangerang', 455, 3),
       ('Kota Tangerang', 456, 3),
       ('Kota Tangerang Selatan', 457, 3),
       ('Kota Bengkulu', 62, 4),
       ('Kab. Bengkulu Selatan', 63, 4),
       ('Kab. Bengkulu Tengah', 64, 4),
       ('Kab. Bengkulu Utara', 65, 4),
       ('Kab. Kaur', 175, 4),
       ('Kab. Kepahiang', 183, 4),
       ('Kab. Lebong', 233, 4),
       ('Kab. Muko Muko', 294, 4),
       ('Kab. Rejang Lebong', 379, 4),
       ('Kab. Seluma', 397, 4),
       ('Kab. Bantul', 39, 5),
       ('Kab. Gunung Kidul', 135, 5),
       ('Kab. Kulon Progo', 210, 5),
       ('Kab. Sleman', 419, 5),
       ('Kota Yogyakarta', 501, 5),
       ('Kota Jakarta Barat', 151, 6),
       ('Kota Jakarta Pusat', 152, 6),
       ('Kota Jakarta Selatan', 153, 6),
       ('Kota Jakarta Timur', 154, 6),
       ('Kota Jakarta Utara', 155, 6),
       ('Kab. Kepulauan Seribu', 189, 6),
       ('Kab. Boalemo', 77, 7),
       ('Kab. Bone Bolango', 88, 7),
       ('Kab. Gorontalo', 129, 7),
       ('Kota Gorontalo', 130, 7),
       ('Kab. Gorontalo Utara', 131, 7),
       ('Kab. Pohuwato', 361, 7),
       ('Kab. Batang Hari', 50, 8),
       ('Kab. Bungo', 97, 8),
       ('Kota Jambi', 156, 8),
       ('Kab. Kerinci', 194, 8),
       ('Kab. Merangin', 280, 8),
       ('Kab. Muaro Jambi', 293, 8),
       ('Kab. Sarolangun', 393, 8),
       ('Kota Sungaipenuh', 442, 8),
       ('Kab. Tanjung Jabung Barat', 460, 8),
       ('Kab. Tanjung Jabung Timur', 461, 8),
       ('Kab. Tebo', 471, 8),
       ('Kab. Bandung', 22, 9),
       ('Kota Bandung', 23, 9),
       ('Kab. Bandung Barat', 24, 9),
       ('Kota Banjar', 34, 9),
       ('Kab. Bekasi', 54, 9),
       ('Kota Bekasi', 55, 9),
       ('Kab. Bogor', 78, 9),
       ('Kota Bogor', 79, 9),
       ('Kab. Ciamis', 103, 9),
       ('Kab. Cianjur', 104, 9),
       ('Kota Cimahi', 107, 9),
       ('Kab. Cirebon', 108, 9),
       ('Kota Cirebon', 109, 9),
       ('Kota Depok', 115, 9),
       ('Kab. Garut', 126, 9),
       ('Kab. Indramayu', 149, 9),
       ('Kab. Karawang', 171, 9),
       ('Kab. Kuningan', 211, 9),
       ('Kab. Majalengka', 252, 9),
       ('Kab. Pangandaran', 332, 9),
       ('Kab. Purwakarta', 376, 9),
       ('Kab. Subang', 428, 9),
       ('Kab. Sukabumi', 430, 9),
       ('Kota Sukabumi', 431, 9),
       ('Kab. Sumedang', 440, 9),
       ('Kab. Tasikmalaya', 468, 9),
       ('Kota Tasikmalaya', 469, 9),
       ('Kab. Banjarnegara', 37, 10),
       ('Kab. Banyumas', 41, 10),
       ('Kab. Batang', 49, 10),
       ('Kab. Blora', 76, 10),
       ('Kab. Boyolali', 91, 10),
       ('Kab. Brebes', 92, 10),
       ('Kab. Cilacap', 105, 10),
       ('Kab. Demak', 113, 10),
       ('Kab. Grobogan', 134, 10),
       ('Kab. Jepara', 163, 10),
       ('Kab. Karanganyar', 169, 10),
       ('Kab. Kebumen', 177, 10),
       ('Kab. Kendal', 181, 10),
       ('Kab. Klaten', 196, 10),
       ('Kab. Kudus', 209, 10),
       ('Kab. Magelang', 249, 10),
       ('Kota Magelang', 250, 10),
       ('Kab. Pati', 344, 10),
       ('Kab. Pekalongan', 348, 10),
       ('Kota Pekalongan', 349, 10),
       ('Kab. Pemalang', 352, 10),
       ('Kab. Purbalingga', 375, 10),
       ('Kab. Purworejo', 377, 10),
       ('Kab. Rembang', 380, 10),
       ('Kota Salatiga', 386, 10),
       ('Kab. Semarang', 398, 10),
       ('Kota Semarang', 399, 10),
       ('Kab. Sragen', 427, 10),
       ('Kab. Sukoharjo', 433, 10),
       ('Kota Surakarta (Solo)', 445, 10),
       ('Kab. Tegal', 472, 10),
       ('Kota Tegal', 473, 10),
       ('Kab. Temanggung', 476, 10),
       ('Kab. Wonogiri', 497, 10),
       ('Kab. Wonosobo', 498, 10),
       ('Kab. Bangkalan', 31, 11),
       ('Kab. Banyuwangi', 42, 11),
       ('Kota Batu', 51, 11),
       ('Kab. Blitar', 74, 11),
       ('Kota Blitar', 75, 11),
       ('Kab. Bojonegoro', 80, 11),
       ('Kab. Bondowoso', 86, 11),
       ('Kab. Gresik', 133, 11),
       ('Kab. Jember', 160, 11),
       ('Kab. Jombang', 164, 11),
       ('Kab. Kediri', 178, 11),
       ('Kota Kediri', 179, 11),
       ('Kab. Lamongan', 222, 11),
       ('Kab. Lumajang', 243, 11),
       ('Kab. Madiun', 247, 11),
       ('Kota Madiun', 248, 11),
       ('Kab. Magetan', 251, 11),
       ('Kota Malang', 256, 11),
       ('Kab. Malang', 255, 11),
       ('Kab. Mojokerto', 289, 11),
       ('Kota Mojokerto', 290, 11),
       ('Kab. Nganjuk', 305, 11),
       ('Kab. Ngawi', 306, 11),
       ('Kab. Pacitan', 317, 11),
       ('Kab. Pamekasan', 330, 11),
       ('Kab. Pasuruan', 342, 11),
       ('Kota Pasuruan', 343, 11),
       ('Kab. Ponorogo', 363, 11),
       ('Kab. Probolinggo', 369, 11),
       ('Kota Probolinggo', 370, 11),
       ('Kab. Sampang', 390, 11),
       ('Kab. Sidoarjo', 409, 11),
       ('Kab. Situbondo', 418, 11),
       ('Kab. Sumenep', 441, 11),
       ('Kota Surabaya', 444, 11),
       ('Kab. Trenggalek', 487, 11),
       ('Kab. Tuban', 489, 11),
       ('Kab. Tulungagung', 492, 11),
       ('Kab. Bengkayang', 61, 12),
       ('Kab. Kapuas Hulu', 168, 12),
       ('Kab. Kayong Utara', 176, 12),
       ('Kab. Ketapang', 195, 12),
       ('Kab. Kubu Raya', 208, 12),
       ('Kab. Landak', 228, 12),
       ('Kab. Melawi', 279, 12),
       ('Kab. Pontianak', 364, 12),
       ('Kota Pontianak', 365, 12),
       ('Kab. Sambas', 388, 12),
       ('Kab. Sanggau', 391, 12),
       ('Kab. Sekadau', 395, 12),
       ('Kota Singkawang', 415, 12),
       ('Kab. Sintang', 417, 12),
       ('Kab. Balangan', 18, 13),
       ('Kab. Banjar', 33, 13),
       ('Kota Banjarbaru', 35, 13),
       ('Kota Banjarmasin', 36, 13),
       ('Kab. Barito Kuala', 43, 13),
       ('Kab. Hulu Sungai Selatan', 143, 13),
       ('Kab. Hulu Sungai Tengah', 144, 13),
       ('Kab. Hulu Sungai Utara', 145, 13),
       ('Kab. Kotabaru', 203, 13),
       ('Kab. Tabalong', 446, 13),
       ('Kab. Tanah Bumbu', 452, 13),
       ('Kab. Tanah Laut', 454, 13),
       ('Kab. Tapin', 466, 13),
       ('Kab. Barito Selatan', 44, 14),
       ('Kab. Barito Timur', 45, 14),
       ('Kab. Barito Utara', 46, 14),
       ('Kab. Gunung Mas', 136, 14),
       ('Kab. Kapuas', 167, 14),
       ('Kab. Katingan', 174, 14),
       ('Kab. Kotawaringin Barat', 205, 14),
       ('Kab. Kotawaringin Timur', 206, 14),
       ('Kab. Lamandau', 221, 14),
       ('Kab. Murung Raya', 296, 14),
       ('Kota Palangka Raya', 326, 14),
       ('Kab. Pulang Pisau', 371, 14),
       ('Kab. Seruyan', 405, 14),
       ('Kab. Sukamara', 432, 14),
       ('Kota Balikpapan', 19, 15),
       ('Kab. Berau', 66, 15),
       ('Kota Bontang', 89, 15),
       ('Kab. Kutai Barat', 214, 15),
       ('Kab. Kutai Kartanegara', 215, 15),
       ('Kab. Kutai Timur', 216, 15),
       ('Kab. Paser', 341, 15),
       ('Kab. Penajam Paser Utara', 354, 15),
       ('Kota Samarinda', 387, 15),
       ('Kab. Bulungan (Bulongan)', 96, 16),
       ('Kab. Malinau', 257, 16),
       ('Kab. Nunukan', 311, 16),
       ('Kab. Tana Tidung', 450, 16),
       ('Kota Tarakan', 467, 16),
       ('Kota Batam', 48, 17),
       ('Kab. Bintan', 71, 17),
       ('Kab. Karimun', 172, 17),
       ('Kab. Kepulauan Anambas', 184, 17),
       ('Kab. Lingga', 237, 17),
       ('Kab. Natuna', 302, 17),
       ('Kota Tanjung Pinang', 462, 17),
       ('Kota Bandar Lampung', 21, 18),
       ('Kab. Lampung Barat', 223, 18),
       ('Kab. Lampung Selatan', 224, 18),
       ('Kab. Lampung Tengah', 225, 18),
       ('Kab. Lampung Timur', 226, 18),
       ('Kab. Lampung Utara', 227, 18),
       ('Kab. Mesuji', 282, 18),
       ('Kota Metro', 283, 18),
       ('Kab. Pesawaran', 355, 18),
       ('Kab. Pesisir Barat', 356, 18),
       ('Kab. Pringsewu', 368, 18),
       ('Kab. Tanggamus', 458, 18),
       ('Kab. Tulang Bawang', 490, 18),
       ('Kab. Tulang Bawang Barat', 491, 18),
       ('Kab. Way Kanan', 496, 18),
       ('Kota Ambon', 14, 19),
       ('Kab. Buru', 99, 19),
       ('Kab. Buru Selatan', 100, 19),
       ('Kab. Kepulauan Aru', 185, 19),
       ('Kab. Maluku Barat Daya', 258, 19),
       ('Kab. Maluku Tengah', 259, 19),
       ('Kab. Maluku Tenggara', 260, 19),
       ('Kab. Maluku Tenggara Barat', 261, 19),
       ('Kab. Seram Bagian Barat', 400, 19),
       ('Kab. Seram Bagian Timur', 401, 19),
       ('Kota Tual', 488, 19),
       ('Kab. Halmahera Barat', 138, 20),
       ('Kab. Halmahera Selatan', 139, 20),
       ('Kab. Halmahera Tengah', 140, 20),
       ('Kab. Halmahera Timur', 141, 20),
       ('Kab. Halmahera Utara', 142, 20),
       ('Kab. Kepulauan Sula', 191, 20),
       ('Kab. Pulau Morotai', 372, 20),
       ('Kota Ternate', 477, 20),
       ('Kota Tidore Kepulauan', 478, 20),
       ('Kab. Aceh Barat', 1, 21),
       ('Kab. Aceh Barat Daya', 2, 21),
       ('Kab. Aceh Besar', 3, 21),
       ('Kab. Aceh Jaya', 4, 21),
       ('Kab. Aceh Selatan', 5, 21),
       ('Kab. Aceh Singkil', 6, 21),
       ('Kab. Aceh Tamiang', 7, 21),
       ('Kab. Aceh Tengah', 8, 21),
       ('Kab. Aceh Tenggara', 9, 21),
       ('Kab. Aceh Timur', 10, 21),
       ('Kab. Aceh Utara', 11, 21),
       ('Kota Banda Aceh', 20, 21),
       ('Kab. Bener Meriah', 59, 21),
       ('Kab. Bireuen', 72, 21),
       ('Kab. Gayo Lues', 127, 21),
       ('Kota Langsa', 230, 21),
       ('Kota Lhokseumawe', 235, 21),
       ('Kab. Nagan Raya', 300, 21),
       ('Kab. Pidie', 358, 21),
       ('Kab. Pidie Jaya', 359, 21),
       ('Kota Sabang', 384, 21),
       ('Kab. Simeulue', 414, 21),
       ('Kota Subulussalam', 429, 21),
       ('Kab. Bima', 68, 22),
       ('Kota Bima', 69, 22),
       ('Kab. Dompu', 118, 22),
       ('Kab. Lombok Barat', 238, 22),
       ('Kab. Lombok Tengah', 239, 22),
       ('Kab. Lombok Timur', 240, 22),
       ('Kab. Lombok Utara', 241, 22),
       ('Kota Mataram', 276, 22),
       ('Kab. Sumbawa', 438, 22),
       ('Kab. Sumbawa Barat', 439, 22),
       ('Kab. Alor', 13, 23),
       ('Kab. Belu', 58, 23),
       ('Kab. Ende', 122, 23),
       ('Kab. Flores Timur', 125, 23),
       ('Kab. Kupang', 212, 23),
       ('Kota Kupang', 213, 23),
       ('Kab. Lembata', 234, 23),
       ('Kab. Manggarai', 269, 23),
       ('Kab. Manggarai Barat', 270, 23),
       ('Kab. Manggarai Timur', 271, 23),
       ('Kab. Nagekeo', 301, 23),
       ('Kab. Ngada', 304, 23),
       ('Kab. Rote Ndao', 383, 23),
       ('Kab. Sabu Raijua', 385, 23),
       ('Kab. Sikka', 412, 23),
       ('Kab. Sumba Barat', 434, 23),
       ('Kab. Sumba Barat Daya', 435, 23),
       ('Kab. Sumba Tengah', 436, 23),
       ('Kab. Sumba Timur', 437, 23),
       ('Kab. Timor Tengah Selatan', 479, 23),
       ('Kab. Timor Tengah Utara', 480, 23),
       ('Kab. Asmat', 16, 24),
       ('Kab. Biak Numfor', 67, 24),
       ('Kab. Boven Digoel', 90, 24),
       ('Kab. Deiyai (Deliyai)', 111, 24),
       ('Kab. Dogiyai', 117, 24),
       ('Kab. Intan Jaya', 150, 24),
       ('Kab. Jayapura', 157, 24),
       ('Kota Jayapura', 158, 24),
       ('Kab. Jayawijaya', 159, 24),
       ('Kab. Keerom', 180, 24),
       ('Kab. Kepulauan Yapen (Yapen Waropen)', 193, 24),
       ('Kab. Lanny Jaya', 231, 24),
       ('Kab. Mamberamo Raya', 263, 24),
       ('Kab. Mamberamo Tengah', 264, 24),
       ('Kab. Mappi', 274, 24),
       ('Kab. Merauke', 281, 24),
       ('Kab. Mimika', 284, 24),
       ('Kab. Nabire', 299, 24),
       ('Kab. Nduga', 303, 24),
       ('Kab. Paniai', 335, 24),
       ('Kab. Pegunungan Bintang', 347, 24),
       ('Kab. Puncak', 373, 24),
       ('Kab. Puncak Jaya', 374, 24),
       ('Kab. Sarmi', 392, 24),
       ('Kab. Supiori', 443, 24),
       ('Kab. Tolikara', 484, 24),
       ('Kab. Waropen', 495, 24),
       ('Kab. Yahukimo', 499, 24),
       ('Kab. Yalimo', 500, 24),
       ('Kab. Fakfak', 124, 25),
       ('Kab. Kaimana', 165, 25),
       ('Kab. Manokwari', 272, 25),
       ('Kab. Manokwari Selatan', 273, 25),
       ('Kab. Maybrat', 277, 25),
       ('Kab. Pegunungan Arfak', 346, 25),
       ('Kab. Raja Ampat', 378, 25),
       ('Kab. Sorong', 424, 25),
       ('Kota Sorong', 425, 25),
       ('Kab. Sorong Selatan', 426, 25),
       ('Kab. Tambrauw', 449, 25),
       ('Kab. Teluk Bintuni', 474, 25),
       ('Kab. Teluk Wondama', 475, 25),
       ('Kab. Bengkalis', 60, 26),
       ('Kota Dumai', 120, 26),
       ('Kab. Indragiri Hilir', 147, 26),
       ('Kab. Indragiri Hulu', 148, 26),
       ('Kab. Kampar', 166, 26),
       ('Kab. Kepulauan Meranti', 187, 26),
       ('Kab. Kuantan Singingi', 207, 26),
       ('Kota Pekanbaru', 350, 26),
       ('Kab. Pelalawan', 351, 26),
       ('Kab. Rokan Hilir', 381, 26),
       ('Kab. Rokan Hulu', 382, 26),
       ('Kab. Siak', 406, 26),
       ('Kab. Majene', 253, 27),
       ('Kab. Mamasa', 262, 27),
       ('Kab. Mamuju', 265, 27),
       ('Kab. Mamuju Utara', 266, 27),
       ('Kab. Polewali Mandar', 362, 27),
       ('Kab. Bantaeng', 38, 28),
       ('Kab. Barru', 47, 28),
       ('Kab. Bone', 87, 28),
       ('Kab. Bulukumba', 95, 28),
       ('Kab. Enrekang', 123, 28),
       ('Kab. Gowa', 132, 28),
       ('Kab. Jeneponto', 162, 28),
       ('Kab. Luwu', 244, 28),
       ('Kab. Luwu Timur', 245, 28),
       ('Kab. Luwu Utara', 246, 28),
       ('Kota Makassar', 254, 28),
       ('Kab. Maros', 275, 28),
       ('Kota Palopo', 328, 28),
       ('Kab. Pangkajene Kepulauan', 333, 28),
       ('Kota Parepare', 336, 28),
       ('Kab. Pinrang', 360, 28),
       ('Kab. Selayar (Kepulauan Selayar)', 396, 28),
       ('Kab. Sidenreng Rappang/Rapang', 408, 28),
       ('Kab. Sinjai', 416, 28),
       ('Kab. Soppeng', 423, 28),
       ('Kab. Takalar', 448, 28),
       ('Kab. Tana Toraja', 451, 28),
       ('Kab. Toraja Utara', 486, 28),
       ('Kab. Wajo', 493, 28),
       ('Kab. Banggai', 25, 29),
       ('Kab. Banggai Kepulauan', 26, 29),
       ('Kab. Buol', 98, 29),
       ('Kab. Donggala', 119, 29),
       ('Kab. Morowali', 291, 29),
       ('Kota Palu', 329, 29),
       ('Kab. Parigi Moutong', 338, 29),
       ('Kab. Poso', 366, 29),
       ('Kab. Sigi', 410, 29),
       ('Kab. Tojo Una-Una', 482, 29),
       ('Kab. Toli-Toli', 483, 29),
       ('Kota Bau-Bau', 53, 30),
       ('Kab. Bombana', 85, 30),
       ('Kab. Buton', 101, 30),
       ('Kab. Buton Utara', 102, 30),
       ('Kota Kendari', 182, 30),
       ('Kab. Kolaka', 198, 30),
       ('Kab. Kolaka Utara', 199, 30),
       ('Kab. Konawe', 200, 30),
       ('Kab. Konawe Selatan', 201, 30),
       ('Kab. Konawe Utara', 202, 30),
       ('Kab. Muna', 295, 30),
       ('Kab. Wakatobi', 494, 30),
       ('Kota Bitung', 73, 31),
       ('Kab. Bolaang Mongondow (Bolmong)', 81, 31),
       ('Kab. Bolaang Mongondow Selatan', 82, 31),
       ('Kab. Bolaang Mongondow Timur', 83, 31),
       ('Kab. Bolaang Mongondow Utara', 84, 31),
       ('Kab. Kepulauan Sangihe', 188, 31),
       ('Kab. Kepulauan Siau Tagulandang Biaro (Sitaro)', 190, 31),
       ('Kab. Kepulauan Talaud', 192, 31),
       ('Kota Kotamobagu', 204, 31),
       ('Kota Manado', 267, 31),
       ('Kab. Minahasa', 285, 31),
       ('Kab. Minahasa Selatan', 286, 31),
       ('Kab. Minahasa Tenggara', 287, 31),
       ('Kab. Minahasa Utara', 288, 31),
       ('Kota Tomohon', 485, 31),
       ('Kab. Agam', 12, 32),
       ('Kota Bukittinggi', 93, 32),
       ('Kab. Dharmasraya', 116, 32),
       ('Kab. Kepulauan Mentawai', 186, 32),
       ('Kab. Lima Puluh Koto/Kota', 236, 32),
       ('Kota Padang', 318, 32),
       ('Kota Padang Panjang', 321, 32),
       ('Kab. Padang Pariaman', 322, 32),
       ('Kota Pariaman', 337, 32),
       ('Kab. Pasaman', 339, 32),
       ('Kab. Pasaman Barat', 340, 32),
       ('Kota Payakumbuh', 345, 32),
       ('Kab. Pesisir Selatan', 357, 32),
       ('Kota Sawah Lunto', 394, 32),
       ('Kab. Sijunjung (Sawah Lunto Sijunjung)', 411, 32),
       ('Kab. Solok', 420, 32),
       ('Kota Solok', 421, 32),
       ('Kab. Solok Selatan', 422, 32),
       ('Kab. Tanah Datar', 453, 32),
       ('Kab. Banyuasin', 40, 33),
       ('Kab. Empat Lawang', 121, 33),
       ('Kab. Lahat', 220, 33),
       ('Kota Lubuk Linggau', 242, 33),
       ('Kab. Muara Enim', 292, 33),
       ('Kab. Musi Banyuasin', 297, 33),
       ('Kab. Musi Rawas', 298, 33),
       ('Kab. Ogan Ilir', 312, 33),
       ('Kab. Ogan Komering Ilir', 313, 33),
       ('Kab. Ogan Komering Ulu', 314, 33),
       ('Kab. Ogan Komering Ulu Selatan', 315, 33),
       ('Kab. Ogan Komering Ulu Timur', 316, 33),
       ('Kota Pagar Alam', 324, 33),
       ('Kota Palembang', 327, 33),
       ('Kota Prabumulih', 367, 33),
       ('Kab. Asahan', 15, 34),
       ('Kab. Batu Bara', 52, 34),
       ('Kota Binjai', 70, 34),
       ('Kab. Dairi', 110, 34),
       ('Kab. Deli Serdang', 112, 34),
       ('Kota Gunungsitoli', 137, 34),
       ('Kab. Humbang Hasundutan', 146, 34),
       ('Kab. Karo', 173, 34),
       ('Kab. Labuhan Batu', 217, 34),
       ('Kab. Labuhan Batu Selatan', 218, 34),
       ('Kab. Labuhan Batu Utara', 219, 34),
       ('Kab. Langkat', 229, 34),
       ('Kab. Mandailing Natal', 268, 34),
       ('Kota Medan', 278, 34),
       ('Kab. Nias', 307, 34),
       ('Kab. Nias Barat', 308, 34),
       ('Kab. Nias Selatan', 309, 34),
       ('Kab. Nias Utara', 310, 34),
       ('Kab. Padang Lawas', 319, 34),
       ('Kab. Padang Lawas Utara', 320, 34),
       ('Kota Padang Sidempuan', 323, 34),
       ('Kab. Pakpak Bharat', 325, 34),
       ('Kota Pematang Siantar', 353, 34),
       ('Kab. Samosir', 389, 34),
       ('Kab. Serdang Bedagai', 404, 34),
       ('Kota Sibolga', 407, 34),
       ('Kab. Simalungun', 413, 34),
       ('Kota Tanjung Balai', 459, 34),
       ('Kab. Tapanuli Selatan', 463, 34),
       ('Kab. Tapanuli Tengah', 464, 34),
       ('Kab. Tapanuli Utara', 465, 34),
       ('Kota Tebing Tinggi', 470, 34),
       ('Kab. Toba Samosir', 481, 34);

INSERT INTO products(name, generic_name, content, manufacturer_id, description, drug_classification_id,
                     product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image)
VALUES ('Panadol 500 mg 10 Kaplet', 'Panadol', 'Paracetamol', 1, 'Obat sakit kepala', 1,
        1, 'Blister', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/panadol.jpg'),
       ('Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 1,
        'Obat sakit kepala', 1,
        1, 'Tablet', 4, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/prx`oducts/saridon.jpg'),
       ('Ultraflu 200 mg 10 Kaplet', 'Paracetamol', 'Paracetamol 200 mg', 1, 'Obat demam dan sakit kepala', 1,
        1, 'Kaplet', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/ultraflu.jpg'),
       ('Amoxan 500 mg 10 Kaplet', 'Amoxicillin', 'Amoxicillin 500 mg', 1, 'Obat infeksi bakteri', 4,
        1, 'Kaplet', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/amoxan.jpg'),
       ('Infacol 100 mg/ml Sirup 60 ml', 'Simeticon', 'Simeticon 100 mg/ml', 1, 'Obat kolik dan kembung untuk bayi', 1,
        1, 'Sirup', 60, 'Botol', 0.2, 100, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/infacol.jpeg'),
       ('Neozep 60 mg/ml Sirup 60 ml', 'Paracetamol', 'Paracetamol 60 mg/ml', 1,
        'Obat demam dan sakit kepala untuk anak-anak', 1,
        1, 'Sirup', 60, 'Botol', 0.2, 100, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/neozep.jpg'),
       ('Oskadon 650 mg 10 Kaplet', 'Paracetamol', 'Paracetamol 650 mg, caffeine 40 mg', 1,
        'Obat sakit kepala, demam, dan nyeri', 1,
        1, 'Kaplet', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/oskadon.jpg'),
       ('Advil 200 mg 10 Tablet', 'Ibuprofen', 'Ibuprofen 200 mg', 1, 'Obat pereda nyeri dan demam', 1,
        1, 'Tablet', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/advil.jpeg'),
       ('Otrivin 0,05% 10 ml', 'Xylometazoline', 'Xylometazoline 0,05%', 1,
        'Obat semprot hidung untuk meredakan hidung tersumbat', 2,
        1, 'Botol', 10, 'Botol', 0.05, 100, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/otrivin.jpeg'),
       ('Mylanta 400 mg/5 ml Sirup 60 ml', 'Magnesium hidroksida, aluminium hidroksida',
        'Magnesium hidroksida 400 mg/5 ml, aluminium hidroksida 500 mg/5 ml', 1, 'Obat gangguan pencernaan', 2,
        1, 'Sirup', 60, 'Botol', 0.2, 100, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/mylanta.jpg'),
       ('Combiflam 500 mg/500 mg 10 Kaplet', 'Paracetamol, ibuprofen', 'Paracetamol 500 mg, ibuprofen 500 mg', 1,
        'Obat pereda nyeri dan demam', 1,
        1, 'Kaplet', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/combiflam.jpeg'),
       ('Promag 400 mg/50 mg 12 Tablet', 'Antasida, simeticon', 'Antasida 400 mg, simeticon 50 mg', 1,
        'Obat sakit maag dan gangguan pencernaan', 2,
        1, 'Tablet', 12, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/promag.jpg'),
       ('Dulcolactol Sirup 60 ml', 'Dulcolactol', 'Per 15 ml : Lactulose 10 gram', 1,
        'DULCOLACTOL merupakan obat konstipasi (sulit buang air besar) yang mengandung Laktulosa yang bekerja menaikkan tekanan osmosa dan suasana asam sehingga feses menjadi lunak. Obat ini dalam penggunaannya dapat dicampur dengan sari buah, air, dan susu.',
        1, 1, 'Liquid', 'Bottle', 'Bottle', 1000, 60, 10, 100,
        'https://byebyesick-bucket.irfancen.com/products/0588003c-bb7c-11ee-bcc7-9d6e84ac2af9.png'),
       ('OBH Combi Plus Batuk Flu Menthol 100 ml', 'OBH',
        'Every 5 ml contain: Succus Liquiritiae 167 mg, Paracetmaol 150 mg, Ammonium Chloride 50 mg, Pseudoephedrin HCl 10 mg, Chlorpheniramin Maleate 1.33 mg.',
        2,
        'OBH COMBI PLUS BATUK FLU merupakan obat batuk dengan kandungan OBH, Paracetamol, Ephedrine HCl, dan Chlorphenamine maleat yang digunakan untuk meredakan batuk disertai gejala-gejala flu seperti demam, sakit kepala, hidung tersumbat, dan bersin-bersin. OBH bekerja sebagai ekspektoran atau peluruh dahak, Paracetamol digunakan sebagai pereda demam dan sakit kepala, Chlorpheniramine maleate bekerja sebagai antihistamin atau anti alergi untuk meredakan gejala alergi, dan Ephedrine HCl sebagai dekongestan hidung atau melonggarkan saluran pernafasan.',
        3, 1, 'Liquid', 'Bottle', 'Bottle', 167, 150, 150, 150,
        'https://byebyesick-bucket.irfancen.com/products/85f5d4e6-bb7d-11ee-bcc7-9d6e84ac2af9.png'),
       ('Panadol Flu Batuk 10 Kaplet', 'Panadol',
        'Paracetamol 500 mg, Phenylephrine HCL 5 mg, dan Dextromethorphan HBr 15 mg.', 1,
        'PANADOL FLU & BATUK merupakan obat batuk dan pereda flu dengan kandungan Paracetamol, Phenylephrine HCI, dan Dextromethorphan HBr. Bekerja sebagai analgesik-antipiretik, masal dekongestan, dan antitusif. Obat ini dapat digunakan untuk meredakan gejala flu seperti: demam, sakit kepala, hidung tersumbat dan batuk tidak berdahak.',
        3, 1, 'tablet', 'strip', 'strip', 100, 40, 10, 10,
        'https://byebyesick-bucket.irfancen.com/products/14bd52fd-bb7e-11ee-bcc7-9d6e84ac2af9.png'),
       ('Sanmol Sirup 60 ml', 'Paracetamol', 'Every 5 ml contain : Paracetamol 120 mg', 1,
        'SANMOL SIRUP merupakan obat dengan kandungan Paracetamol. Obat ini digunakan untuk meringankan rasa sakit pada keadaan sakit kepala, sakit gigi dan menurunkan demam. Sanmol bekerja pada pusat pengatur suhu di hipotalamus untuk menurunkan suhu tubuh (antipiretik) serta menghambat sintesis prostaglandin sehingga dapat mengurangi nyeri ringan sampai sedang',
        1, 1, 'Liquid', 'Bottle', 'Bottle', 120, 30, 60, 30,
        'https://byebyesick-bucket.irfancen.com/products/8008dc3d-bb7e-11ee-bcc7-9d6e84ac2af9.png');


INSERT INTO pharmacies(name, address, sub_district, district, city, province, postal_code, latitude, longitude,
                       pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days,
                       pharmacy_admin_id)
VALUES ('Kimia Farma Kuningan', 'Jalan Gatau', 'Kuningan', 'Setia Budi', 153, 6, '12950', '-6.230060', '106.827363',
        'M. Irfan Junaidi',
        '69696969', '08123456789', '0-20', 'mon,tue,wed,thu,fri', 2),
       ('Kimia Farma Pasar Minggu', 'Jalan Jalan', 'Ragunan', 'Pasar Minggu', 153, 6, '12560', '-6.290963',
        '106.817317',
        'M. Yafi Al Hakim',
        '42042042', '08998239082', '0-22', 'mon,tue,wed,thu,fri,sat,sun', 2),
       ('Apotek Sinar Jaya', 'Jalan Jaya', 'Merdeka', 'Medan Baru', 278, 34, '20222', '3.576816', '98.659355',
        'Victor Castor',
        '10090900', '0892308932', '0-18', 'mon,tue,wed,thu,fri', 4);

INSERT INTO pharmacy_products(pharmacy_id, product_id, is_active, price, stock)
VALUES (1, 1, true, '12000', 100),
       (1, 2, true, '2500', 95),
       (2, 3, true, '5000', 87),
       (2, 4, true, '15000', 95),
       (3, 5, true, '20000', 65),
       (2, 6, true, '10000', 34),
       (3, 7, true, '10000', 273),
       (3, 8, true, '10000', 22),
       (1, 9, true, '10000', 93),
       (2, 10, true, '10000', 77),
       (2, 2, true, '5000', 95),
       (3, 1, true, '15000', 87);

INSERT INTO pharmacy_shipping_methods(pharmacy_id, shipping_method_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (2, 2);

INSERT INTO addresses(name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id)
values ('rumah', 'jl kripat', 'ciangasna', 'gunung putri', 78,9, '16968', '-6.354846', '106.952082', 1, 6);

INSERT INTO transactions(date, payment_proof, transaction_status_id, payment_method_id, address, user_id, total_payment)
values (now(), '',1,1,'jl kripat',6,20000);

INSERT INTO orders(date, pharmacy_id, no_of_items, pharmacy_address, shipping_method_id, shipping_cost, total_payment, transaction_id)
values (now(), 1, 2, 'Jalan Gatau', 1, 5000, 10000, 1),
       (now(), 2, 1, 'Jalan Jalan', 1, 1000, 10000, 1),
       (now(), 3, 1, 'Jalan Jaya', 1, 0, 0, 1);

INSERT INTO order_details(order_id, product_id, quantity, name, generic_name, content, description, image, price)
values (1, 1, 1, 'Panadol 500 mg 10 Kaplet', 'Panadol', 'Paracetamol', 'Obat sakit kepala', '', '12000'),
       (1, 2, 1, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', '2500'),
       (2, 2, 3, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', '2500'),
       (3, 2, 3, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', '0');


INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
values (1,1,false,''), (1,2,true,''), (2,1,true,''),(3,1,true,'');

INSERT INTO consultation_session_statuses(name)
VALUES ('Ongoing'),
       ('Ended');


-- CREATE FUNCTIONS --

-- SQL code to create a function that calculates the distance in kilometers
-- using the haversine formula

-- Define the radius of the Earth in kilometers
CREATE
OR REPLACE FUNCTION earth_radius()
    RETURNS DECIMAL AS
$$
BEGIN
RETURN 6371;
END;
$$
LANGUAGE plpgsql;

-- Create a function that takes two pairs of coordinates as input
-- and returns the distance between them as output
CREATE
OR REPLACE FUNCTION distance(lat1 VARCHAR, lon1 VARCHAR, lat2 VARCHAR, lon2 VARCHAR) RETURNS DECIMAL AS
$$
DECLARE
    -- Convert degrees to radians
radLat1 DECIMAL := RADIANS(lat1::DECIMAL);
    radLon1
DECIMAL := RADIANS(lon1::DECIMAL);
    radLat2
DECIMAL := RADIANS(lat2::DECIMAL);
    radLon2
DECIMAL := RADIANS(lon2::DECIMAL);
    -- Calculate the difference between the coordinates
    dLat
DECIMAL := radLat2 - radLat1;
    dLon
DECIMAL := radLon2 - radLon1;
    -- Apply the haversine formula
    a
DECIMAL := SIN(dLat / 2) ^ 2 + COS(radLat1) * COS(radLat2) * SIN(dLon / 2) ^ 2;
    c
DECIMAL := 2 * ATAN2(SQRT(a), SQRT(1 - a));
BEGIN
    -- Calculate the distance in kilometers
RETURN earth_radius() * c;
END
$$
LANGUAGE plpgsql;
