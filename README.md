# go-musthave-diploma-tpl

**users**                       
| id | login | password |   
|----|-------|----------|
| serial not null unique | varchar(256) | varchar(256) (hashed) |

**orders**                                                
| id | number | status | accrual | uploaded_at |
|----|--------|--------|---------|-------------|
| serial not null unique | varchar(256) | varchar(256) | decimal | timestamptz |

**users_orders**                                               
| id | user_id | order_id |                                
|----|---------|----------|                                
| serial not null unique | int references users (id) on delete cascade | int references orders (id) on delete cascade |

**user_balance**
| id | user_id | current_balance | withdrawn |
|----|---------|-----------------|-----------|
| serial not null unique | int | decimal default 0 | decimal default 0 |
                     
**withdrawals**
| id | user_id | order_id | sum | processed_at |
|----|---------|----------|-----|--------------|
| serial not null unique | int | varchar(256) | decimal | timestamptz |