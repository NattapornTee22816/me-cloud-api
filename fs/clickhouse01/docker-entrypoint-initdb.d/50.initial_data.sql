insert into organize(uuid, name, contact, website, status, verify_status, created_time, created_by, updated_time, updated_by)
select generateUUIDv4(),
       'ME-Cloud',
       '',
       '',
       'active',
       'verified',
       current_timestamp(),
       toUUID('00000000-0000-0000-0000-000000000000'),
       current_timestamp(),
       toUUID('00000000-0000-0000-0000-000000000000');

insert into application(organize_uuid, app_uuid, app_name, app_description, default_app, status, created_time, created_by, updated_time, updated_by)
select uuid,
       generateUUIDv4(),
       'FrontEnd',
       '',
       'Y',
       'active',
       created_time,
       created_by,
       updated_time,
       updated_by
from organize
where name = 'ME-Cloud';

insert into api_keys(app_uuid, key_uuid, access_key, access_key_expire_time, refresh_key, refresh_key_expire_time, permissions, status, created_time, created_by, updated_time, updated_by)
select app_uuid,
       generateUUIDv4(),
       base58Encode(randomPrintableASCII(32)),
       toDateTime(addYears(current_timestamp('UTC'), 10)),
       base58Encode(randomPrintableASCII(36)),
       toDateTime(addYears(current_timestamp('UTC'), 10)),
       array('*'),
       'active',
       current_timestamp('UTC'),
       toUUID('00000000-0000-0000-0000-000000000000'),
       current_timestamp('UTC'),
       toUUID('00000000-0000-0000-0000-000000000000')
from application
where app_name = 'FrontEnd';

select base58Encode(randomPrintableASCII(32)), base58Encode(randomPrintableASCII(36))
