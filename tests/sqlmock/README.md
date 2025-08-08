# 替换空sql
使用sqlmock替换一条空sql，然后选择一个一个样例点执行，根据报错信息找的需要替换的sql
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 替换空sql
        mock.ExpectQuery("")

        ro := RoleServiceImpl{}
        got, got1 := ro.GetRole(tt.args.id)
        assert.Equalf(t, tt.want, got, "GetRole(%v)", tt.args.id)
        assert.Equalf(t, tt.want1, got1, "GetRole(%v)", tt.args.id)
    })
}
```

# ExpectQuery

```go
//这个是使用正则来匹配sql语句，需要使用“\\”来转义字符
mock.ExpectQuery("SELECT \\* FROM `cmdb_app` WHERE app_code = \\? AND is_deleted = \\? ORDER BY `cmdb`\\.`id` limit 1")
mock.ExpectQuery("^SELECT \\* FROM `agents` WHERE voip_number = \\?").WillReturnRows(rows)
// 返回结果为被更改的数据行数
mock.ExpectExec("^INSERT INTO `apps` \\(.*\\) VALUES \\(.*\\)").WillReturnResult(sqlmock.NewResult(1, 1))
```

# regexp.QuoteMeta
```go
//直接匹配sql语句
mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cmdb_app` WHERE app_code = ? AND is_deleted = ? ORDER BY `cmdb`.`id` limit 1"))
```

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        mock.ExpectQuery("SELECT * FROM `role` WHERE `role`.`id` = ? AND `role`.`is_delete` = ? LIMIT 1")

        ro := RoleServiceImpl{}
        got, got1 := ro.GetRole(tt.args.id)
        assert.Equalf(t, tt.want, got, "GetRole(%v)", tt.args.id)
        assert.Equalf(t, tt.want1, got1, "GetRole(%v)", tt.args.id)
    })
}
```

# AddRow
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        mock.ExpectQuery("SELECT * FROM `role` WHERE `role`.`id` = ? AND `role`.`is_delete` = ? LIMIT 1").
            WithArgs(tt.args.id, 0).  // 补充占位符缺失的数据
            WillReturnRows(
                sqlmock.NewRows([]string{"id", "name", "is_delete"}).  // 补充sql返回的字段
                AddRow(tt.want.ID, tt.want.Name, tt.want.IsDelete))  // 补充字段内容
        got, got1 := GetUserService().GetRole(tt.args.id)
        assert.Equalf(t, tt.want, got, "GetRole(%v)", tt.args.id)
        assert.Equalf(t, tt.want1, got1, "GetRole(%v)", tt.args.id)
    })
}
```

# WillReturnRows
```go
// 模拟 SQL 查询并设置预期结果
rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Alice").AddRow(2, "Bob")
mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)
```

# WillReturnResult
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO `role` (`name`) VALUES (?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`)").
            WithArgs(tt.args.role.Name).
                WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectCommit()
        ro := RoleServiceImpl{}
        assert.Equalf(t, tt.want, ro.AddRole(tt.args.role), "AddRole(%v)", tt.args.role)
    })
}
```

# WillReturnError
```go
// 返回error
mock.ExpectExec("^INSERT INTO `apps` \\(.*\\) VALUES \\(.*\\)").WillReturnError(errors.New("some error"))
```

# 测试事务
```go
mock.ExpectBegin() // 匹配begin
mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
mock.ExpectCommit() // 匹配commit
```

# 查询列表的测试
```go
func TestRoleServiceImpl_GetRoleList(t *testing.T) {
    tests := []struct {
        name string
        want *[]*model.Role
        want1 *errno.Errno
    }{
        {
            name: "获取角色列表测试点",
            want: &[]*model.Role{
                {ID: 1, Name: "学生", IsDelete: 0},
                {ID: 2, Name: "教师", IsDelete: 0},
                {ID: 3, Name: "专业管理员", IsDelete: 0},
                {ID: 4, Name: "院管理员", IsDelete: 0},
                {ID: 5, Name: "校管理员", IsDelete: 0},
                {ID: 6, Name: "管理员", IsDelete: 0},
                {ID: 7, Name: "超级管理员", IsDelete: 0},
            },
            want1: nil,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 匹配查询sql 并模拟查询到的数据
            mock.ExpectQuery("SELECT * FROM `role` WHERE `role`.`is_delete` = ?").
                WithArgs(0).
                WillReturnRows(
                    sqlmock.NewRows([]string{"id", "name", "is_delete"}).
                    AddRow((*tt.want)[0].ID, (*tt.want)[0].Name, (*tt.want)[0].IsDelete).
                    AddRow((*tt.want)[1].ID, (*tt.want)[1].Name, (*tt.want)[1].IsDelete).
                    AddRow((*tt.want)[2].ID, (*tt.want)[2].Name, (*tt.want)[2].IsDelete).
                    AddRow((*tt.want)[3].ID, (*tt.want)[3].Name, (*tt.want)[3].IsDelete).
                    AddRow((*tt.want)[4].ID, (*tt.want)[4].Name, (*tt.want)[4].IsDelete).
                    AddRow((*tt.want)[5].ID, (*tt.want)[5].Name, (*tt.want)[5].IsDelete).
                    AddRow((*tt.want)[6].ID, (*tt.want)[6].Name, (*tt.want)[6].IsDelete))
            ro := RoleServiceImpl{}
            got, got1 := ro.GetRoleList()
            assert.Equalf(t, tt.want, got, "GetRoleList()")
            assert.Equalf(t, tt.want1, got1, "GetRoleList()")
        })
    }
}
```
