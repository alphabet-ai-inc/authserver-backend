INSERT INTO apps(name, path, release, init, web, title, created, updated) VALUES
('conta','conta','','conta.php','','BANK-ZERO Contabilidad en la nube', EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-03 17:47:19'),EXTRACT('EPOCH' FROM TIMESTAMP'2022-05-03 17:47:19')),
('AutServer','AutServer','','admincp.php','','AutServer',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-03 13:00:09'), EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-03 13:00:09')),
('appadmin','appadmin','','appadmin.php','','BANK-ZERO Administración de aplicaciones',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-06 19:21:53'), EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-06 19:21:53'));

INSERT INTO companies(name, created, updated) VALUES
('BiblioData', EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'), EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05')),
('Test company', EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 20:48:38'),EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 20:48:38'));

INSERT INTO config(name, value) VALUES
('cookie_prefix','qls3_'),
('max_username','15'),
('min_username','2'),
('max_password','15'),
('min_password','4'),
('cookie_path','/AutServer/'),
('cookie_secure','0'),
('cookie_length','1209600'),
('cookie_domain','dev.bank-zero.com'),
('max_tries','5'),
('user_regex','/^[a-zA-Z_\\x7f-\\xff][a-zA-Z0-9_\\x7f-\\xff]*$/'),
('security_image','yes'),
('activation_type','0'),
('login_redirect','/AutServer/admincp.php'),
('logout_redirect','/AutServer/login.php'),
('max_upload_size','1048576'),
('auth_registration','0'),
('current_version','3.1.11'),
('redirect_type','1'),
('online_users_format','{username}'),
('online_users_separator',',');
INSERT INTO dbs(dsn, options, alt_dsn, created, updated) VALUES
('pgsql:host=127.0.0.1;port=5432;dbname=bzero_bz1001','','pgsql:host=$host;port=5432;dbname=$db',EXTRACT('EPOCH' FROM TIMESTAMP '2022-04-28 16:39:59'), EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('pgsql:host=127.0.0.1;port=5432;dbname=bzero_bz2','','pgsql:host=$host3;port=5432;dbname=$db',EXTRACT('EPOCH' FROM TIMESTAMP '2022-04-28 6:40:14'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('pgsql:host=127.0.0.1;port=5432;dbname=bzero_bz1002','','',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-18 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));

INSERT INTO dbsauth(name, db_id, grants, user_id, group_id, grants_options, created, updated) VALUES
('superuser',1,'SELECT,INSERT,UPDATE,DELETE',1,1,1,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('ContaDev',1,'select,insert',2,0,0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('ContaDevAdmin',1,'select, insert, update',3,2,0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('ContaDev2',3,'select,insert,update',3,0,1,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'), EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));
INSERT INTO groups(name, profile_id, dbsauth_id, is_public, leader, lan, company_id, created, updated) VALUES
('Admin',1,0,true,0,'es',0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('autuser',2,1,true,0,'es',0, EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));
INSERT INTO profiles(name, app_ids, db_ids,created, updated) VALUES
('admin', 'all', 'all',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('director', '1', '1',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('manager', '1', '1',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('salesrep', '1', '1',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('accountant', '1', '1',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('project_leader', '1', '1',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));
INSERT INTO sessions(value, time, created, updated) VALUES
('dca7e953034d6ffda3006b79cb2da5a12cddda1f','2022-05-07 20:48:38',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));
INSERT INTO users(username, password, code, active, last_login, last_session, blocked, tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action, last_app, last_db, lan, company_id,created, updated ) VALUES
('jpassano','4d2ed2743381d892c7d172c4fa0ac70b82db70ad','b3059fc8a92a6f3cbc863be325dafcc067a8caf1','yes',EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'a898daa560f744136438654fbcc8fa2c25416376','no',0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'jpassano@aztechventures.com',1,1,0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'2022-05-07 17:39:05',1,2,'es',1,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('felipep','f0ebddc6b86d257c16fb2386457cec05ab8bbe1a','888401144a858f50321cdc3913944eb0987dc902','yes',EXTRACT('EPOCH' FROM TIMESTAMP '2022-06-20 17:39:05'),'a898daa560f744136438654fbcc8fa2c25416376','no',0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'fperez@aztechventures.com',6,2,0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'2022-05-07 17:39:05',2,1,'es',1,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57')),
('aescala','31c79dd16050c7cfc64d8dc8a27c1e606a9bfd36','afbb9fad6cc1822db8b6aa262504221680a67a7e','yes',EXTRACT('EPOCH' FROM TIMESTAMP '2022-07-08 17:39:05'),'dca7e953034d6ffda3006b79cb2da5a12cddda1f','no',0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'), 'anaescala@gmail.com',6,2,0,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-07 17:39:05'),'2022-05-07 17:39:05',1,1,'es', 2,EXTRACT('EPOCH' FROM TIMESTAMP '2022-05-28 15:47:57'),EXTRACT('EPOCH' FROM TIMESTAMP '2024-05-28 15:47:57'));



update users set email='admin@example.com' where id=1;


update users set password='$2a$10$5z20bnzXmir2Tkx3R5tlkufgqkmmZHzuAWni80LVAy.ch.L9lkBvG' where id=1;
update users set password='$2a$10$O0m5UkQUgFrKfhJhlvCYGe3z1cQv.i1Rb04f0MKFYJyQA4Ymju6wu' where id=2;
update users set password='$2a$10$tO.JMVcRgEktbTpARXNz1enaTHx2A1JfaWkbxT/oz/kkKFpmEVHrS' where id=3;

