create table user (
    'id' integer primary key autoincrement,
    'name' text not null,
    'middlename' text,
    'surname' text not null,
    'password' text not null,
    'passport_series' text not null,
    'passport_number' text not null,
    'phone' text not null,
    'email' text not null unique,
    'role' text not null
        check (role in ('client', 'manager', 'operator', 'admin')), 
    unique(passport_series, passport_number)
);
create index idx_user_email on user(email);

create table bank (
    'id' integer primary key autoincrement,
    'name' text not null,
    'description' text not null,
    'bic' text not null,
    'address' text not null,
    'rating' integer not null
);

create table enterprise (
    'id' integer primary key autoincrement,
    'name' text not null,
    'unp' text not null,
    'address' text not null,
    'bank_id' integer not null,
    foreign key (bank_id) references bank(id)
);

create table enterprise_manager (
    'enterprise_id' integer not null,
    'user_id' integer not null,
    primary key (enterprise_id, user_id),
    foreign key (enterprise_id) references enterprise(id),
    foreign key (user_id) references user(id)
);

create table user_account (
    'id' integer primary key autoincrement,
    'number' text not null,
    'balance' integer not null,
    'currency' text not null,
    'user_id' integer not null,
    'bank_id' integer not null,
    'hold_balance' integer not null default 0,
    'freezing' integer not null default 0,
    foreign key (bank_id) references bank(id)
);
create unique index idx_user_account_number on user_account(number);

create table enterprise_account (
    'id' integer primary key autoincrement,
    'number' text not null,
    'balance' integer not null,
    'currency' text not null,
    'enterprise_id' integer not null,
    'bank_id' integer not null,
    foreign key (bank_id) references bank(id)
);
create unique index idx_enterprise_account_number on enterprise_account(number);

create table system_transaction (
    'id' integer primary key autoincrement,
    'amount' integer not null,
    'currency' text not null,
    'timestamp' datetime not null default current_timestamp,
    'description' text,
    'status' text not null default 'pending' 
        check (status in ('pending', 'completed', 'cancelled')),
    
    'source_account_id' integer not null,
    'destination_account_id' integer not null,
    
    'source_account_type' text not null 
        check (source_account_type in ('user', 'enterprise')),
    'destination_account_type' text not null 
        check (destination_account_type in ('user', 'enterprise')),
    
    'type' text not null 
        check (type in ('transfer', 'salary')),
    
    'source_bank_id' integer not null,
    'destination_bank_id' integer not null,
    'initiated_by_user_id' integer not null,
    
    foreign key (source_bank_id) references bank(id),
    foreign key (destination_bank_id) references bank(id),
    foreign key (initiated_by_user_id) references user(id)
);

create table system_credit (
    'id' integer primary key autoincrement,
    'amount' integer not null,
    'term' integer not null,
    'currency' text not null,
    'timestamp' datetime not null default current_timestamp,
    'description' text,
    'status' text not null default 'pending' 
        check (status in ('pending', 'completed', 'cancelled')),
    
    'source_account_id' integer not null,
    'source_bank_id' integer not null,
    'initiated_by_user_id' integer not null,
    
    foreign key (source_bank_id) references bank(id),
    foreign key (initiated_by_user_id) references user(id)
);
create index idx_transaction_source_account on system_transaction(source_account_id, source_account_type);
create index idx_transaction_destination_account on system_transaction(destination_account_id, destination_account_type);
create index idx_transaction_source_bank on system_transaction(source_bank_id);
create index idx_transaction_destination_bank on system_transaction(destination_bank_id);
create index idx_transaction_initiator on system_transaction(initiated_by_user_id);
