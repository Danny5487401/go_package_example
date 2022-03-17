# 映射关系  

## Go与字段类型对应表  

|go类型 |映射方法 |xorm类型 |
|----- |----- |----- |
implemented Conversion|	Conversion.ToDB / Conversion.FromDB|	Text|
int, int8, int16, int32, uint, uint8, uint16, uint32|	|	Int|
int64, uint64| |		BigInt|
float32|	|	Float|
float64	|	|Double|
complex64, complex128|	json.Marshal / json.UnMarshal|	Varchar(64)|
[]uint8	|	|Blob|
array, slice, map except []uint8	|json.Marshal / json.UnMarshal|	Text|
bool|	1 or 0|	Bool|
string|		|Varchar(255)|
time.Time|		|DateTime|
cascade struct|	primary key field valu|	BigInt|
struct	|json.Marshal / json.UnMarshal|	Text|
Others	|	|Text|

## xorm与数据库类型对照  
| xorm | mysql | sqlite3 | postgres | 备注 | 
|-----|-----|------|-----|-----|
BIT	|BIT	|INTEGER|BIT	|
TINYINT|	TINYINT	|INTEGER|	SMALLINT|
SMALLINT|	SMALLINT|	INTEGER|	SMALLINT	|
MEDIUMINT|	MEDIUMINT|	INTEGER	|INTEGER	|
INT|	INT|	INTEGER|	INTEGER	|
INTEGER	|INTEGER|	INTEGER|	INTEGER	|
BIGINT	|BIGINT|	INTEGER	|BIGINT	|
CHAR	|CHAR|	TEXT	|CHAR	|
VARCHAR	|VARCHAR	|TEXT	|VARCHAR|	
TINYTEXT|	TINYTEXT|	TEXT|	TEXT|	
TEXT|	TEXT|	TEXT	|TEXT|	
MEDIUMTEXT	|MEDIUMTEXT	|TEXT	|TEXT|	
LONGTEXT	|LONGTEXT|	TEXT	|TEXT|	
BINARY|	BINARY	|BLOB|	BYTEA|	
VARBINARY	|VARBINARY|	BLOB	|BYTEA	|
DATE	|DATE	|NUMERIC|DATE|	
DATETIME|	DATETIME|	NUMERIC	|TIMESTAMP	|
TIME	|TIME	|NUMERIC|	TIME	|
TIMESTAMP	|TIMESTAMP|	NUMERIC|	TIMESTAMP|	
TIMESTAMPZ	|TEXT|	TEXT	|TIMESTAMP with zone	|timestamp with zone info|
FLOAT|	FLOAT	|REAL	|REAL	|
DOUBLE|	DOUBLE|	REAL	DOUBLE| PRECISION	|
DECIMAL|	DECIMAL|	NUMERIC|	DECIMAL	|
NUMERIC|	NUMERIC	|NUMERIC	|NUMERIC	|
TINYBLOB|	TINYBLOB|	BLOB	|BYTEA	|
BLOB|	BLOB	|BLOB|	BYTEA	|
MEDIUMBLOB|	MEDIUMBLOB|	BLOB|	BYTEA	|
LONGBLOB|	LONGBLOB|	BLOB	|BYTEA|	
BYTEA|	BLOB|	BLOB	|BYTEA|	
BOOL|	TINYINT	|INTEGER	|BOOLEAN|	
SERIAL|	INT|	INTEGER	|SERIAL	|auto increment|
BIGSERIAL|	BIGINT|	INTEGER|	BIGSERIAL	|auto increment
