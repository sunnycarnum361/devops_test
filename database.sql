create schema if not exists devops collate latin1_swedish_ci;

create table if not exists person
(
	id int not null
		primary key,
	name varchar(50) null,
	age int null,
	address varchar(250) null
);

insert into devops.person(id,name,age,address) values (10,"tom",30,"london")