# DisNote

Disnote - это Discord бот для ведения заметок 

```
-help - покажет помощь
-add - добавить заметку
-del <num> - удалить заметку
-clear - очистить канал и заного написать заметки
```

Конфигурация бота должна находиться в `~/.config/disnote.yml`
Пример конфигурации находится в файле `disnote.yml.example`

```yaml
discord:
  token: discord token
  channelid: id channel
redis:
  address: address redis server
  password: password redis server
```

