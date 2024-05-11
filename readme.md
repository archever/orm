
# orm

## TODO
- sql builder
    - [x] builder 和 action 分离
    - [x] 支持 join
    - [x] 如果没有 join, 不需要使用 `表名.字段名`
    - [x] 支持 subquery
    - [x] 无需指定表名的场景, 不用加表名称
- 完善 scan
    - [ ] 完善 支持 json
    - [ ] 支持同一个表的字段被多次bind的场景
- 完善 payload
    - [ ] payload 自动更新
- 补充测试用例
