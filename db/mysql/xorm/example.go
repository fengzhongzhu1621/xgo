package xorm

// CREATE TABLE `xorm_student` (
//
//	`name` varchar(64) NOT NULL,
//	`age` int NOT NULL DEFAULT '0',
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	PRIMARY KEY (`id`),
//	UNIQUE KEY `UQE_xorm_student_name` (`name`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormStudent struct {
	Name string `xorm:"not null unique VARCHAR(64)"`
	Age  int    `xorm:"not null default 0 INT(10)"`
	Id   int64  `xorm:"bigint(20) pk autoincr"`
}
