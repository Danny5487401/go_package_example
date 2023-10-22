<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [RBAC](#rbac)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# RBAC 



casbin rbac 使用注意


- token处理： sub 默认是第一个位置。可以使用e.SetFieldIndex()修改
- Casbin doesn't distinguish between roles and users in its RBAC.建议加前缀使用  role::admin