create table users (
    'id' integer primary key autoincrement,
    'name' text not null,
    'surname' text not null,
    'username' text not null,
    'passport_series' text not null,
    'passport_number' text not null,
    'phone' text not null,
    'email' text not null,
    'role' text not null
);

insert into users (name, surname,username, passport_series, passport_number, phone, email, role) 
values ('Elisey', 'Akulich', 'loshok', '22', '333333', '+79534523344', 'loh@ebaniy.com', 'admin');