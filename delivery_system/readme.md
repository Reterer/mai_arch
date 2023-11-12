# Delivery system

### auth
Я нашел [статью](https://habr.com/ru/companies/spectr/articles/715290/), где описываются паттерны авторизации для ms архитектур.<br>
Я хочу использовать вариант, где другие сервисы (в будущем nginx) ходят в client_serivce, что бы проверить доступ.<br>
Сейчас используется basic авторизация. Однако хочу добавить access токены, что бы сделать проверку быстрее.

### client_serivce
- [x] Client Service
- [x] Client Repo Memory
- [x] Auth Package
  - [x] pass hashed
  - [x] external (логика хождения в client_service)
- [~] Auth логика (В планах использовать access токен)
  - [x] register
  - [~] login (выдавать access токен)
  - [~] check (проверять access токен)
  - [ ] logout (отзывать access токен) 
- [x] Docker
- [ ] Service tests
### item_service
- [ ] Api
- [ ] Service
- [ ] Repo Memory
- [ ] Docker
- [ ] Service tests
### delivery_service
- [ ] Api
- [ ] Service
- [ ] Repo Memory
- [ ] Docker
- [ ] Service tests
### db_service
- [ ] Docker
- [ ] SQL init
- [ ] Client Repo
- [ ] Delivery Repo
- [ ] Item Repo
### General
- [ ] docker-compose
- [ ] tests