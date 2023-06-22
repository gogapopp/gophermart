# go-musthave-diploma-tpl

**users**                       
| id | login | password |   
|----|-------|----------|   
id - serial not null unique 
login - varchar(256)        
password - varchar(256)  (hashed)                  

**orders**                                                
| id | number | status | accrual | uploaded_at |
|----|--------|--------|---------|-------------|
id - serial not null unique                     
number - varchar(256)                           
status - varchar(256)                           
accrual - decimal                               
uploaded_at - timestamptz

**users_orders**                                               
| id | user_id | order_id |                                
|----|---------|----------|                                
id - serial not null unique                                
user_id - int references users (id) on delete cascade      
order_id - int references orders (id) on delete cascade    

**user_balance**
| id | user_id | current_balance | withdrawn |
|----|---------|-----------------|-----------|
id - serial not null unique
user_id - int
current_balance - decimal default 0
withdrawn  - decimal default 0                                
**withdrawals**
| id | user_id | order_id | sum | processed_at |
|----|---------|----------|-----|--------------|
id - serial not null unique
user_id - int
order_id - varchar(256)
sum - decimal
processed_at - timestamptz