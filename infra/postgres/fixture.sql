create table public.company (id SERIAL PRIMARY KEY,name varchar(300),duns varchar(27),spin bigint,city varchar(200),country varchar(90),type varchar(255) DEFAULT 'REGULAR');
create table public.iam_user (id SERIAL PRIMARY KEY,iam_id varchar(50),email varchar(200),name varchar(200),company_id bigint,FOREIGN KEY (company_id) REFERENCES company(id));

INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('Ford','114720025',1,'Detroit','USA','REGULAR');
INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('Opel','214720025',2,'Detroit','Germany','REGULAR');
INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('VW','314720025',3,'Wolfsburg','Germany','EMAIL');
INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('Toyota','414720025',4,'Tokio','Japan','REGULAR');
INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('Audi','514720025',5,'Sindelfingen','Germany','EMAIL');
INSERT INTO public.company (name, duns, spin, city, country, type) VALUES ('Crysler','614720025',6,'Leeds','England','REGULAR');

INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_1','mustermann@ford.com','Mustermann',1);
INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_2','mueller@ford.com','Müller',1);
INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_3','meier@ford.com','Meier',1);
INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_4','rudolf@opel.com','Rudolf',2);
INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_5','herbert@vw.com','Herbert',3);
INSERT INTO public.iam_user (iam_id, email, name, company_id) VALUES ('ID_ABCD_6','schmidt@toyota.com','Schmidt',4);