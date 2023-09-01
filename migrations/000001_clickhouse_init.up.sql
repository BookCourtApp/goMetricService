create table if not exists metrics
(
	TimeStamp 	DateTime,
	IsApp 		UInt8,
	IsAuth 		UInt8,
	IsNew 		UInt8,
	ResWidth 	UInt16,
	ResHeight 	UInt16,
	UserAgent 	String,
	UserId 		String,
	SessionID 	String,
	DeviceType 	String,
	Reffer 		String,
	Stage 		LowCardinality(String),
	Action 		LowCardinality(String),
	ExtraKeys 	Array(String),
	ExtraValues Array(String)
) 
engine = MergeTree() 
order by Action;
