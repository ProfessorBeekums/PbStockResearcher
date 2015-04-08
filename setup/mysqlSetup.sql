create database pb_stock_researcher;

create table company (
  cik bigint unsigned not null
, name varchar(255) not null
, PRIMARY KEY(cik)
) ENGINE=INNODB;

# wide table for ease of screening queries
# this will need to be frequently updated
create table financial_report (
  financial_report_id int unsigned auto_increment
, cik bigint unsigned not null
, year int unsigned not null
, quarter tinyint unsigned not null
, report_file_id int unsigned not null
, revenue bigint not null default 0
, operating_expense bigint not null default 0
, net_income bigint not null default 0
, current_assets bigint not null default 0
, total_assets bigint not null default 0
, current_liabilities bigint not null default 0
, total_liabilities bigint not null default 0
, operating_cash bigint not null default 0
, capital_expenditures bigint not null default 0
, PRIMARY KEY(financial_report_id)
, KEY (cik)
, KEY (year, quarter)
, KEY (report_file_id)
) ENGINE=INNODB;

create table financial_report_raw_fields (
  financial_report_raw_field_id int unsigned auto_increment
, financial_report_id int unsigned not null
, field_name varchar(255) not null
, field_value bigint not null default 0
, PRIMARY KEY(financial_report_raw_field_id)
, KEY(financial_report_id)
) ENGINE=INNODB;

create table report_file(
  report_file_id int unsigned auto_increment
, cik bigint unsigned not null
, year int unsigned not null
, quarter int unsigned not null
, filepath varchar(255) not null
, form_type varchar(64) not null
, parsed bool not null default false
, parse_error bool not null default false
, PRIMARY KEY(report_file_id)
, KEY(cik)
, KEY(year, quarter)
) ENGINE=INNODB;