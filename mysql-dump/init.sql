create table gorm.posts (
	post_id INT(6) unsigned auto_increment primary key, 
	post_content VARCHAR(100), 
	publish_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

insert into gorm.posts (post_content) values ('aaaaaaaaaa');
insert into gorm.posts (post_content) values ('bbbbbbbbbb');

