# s21auto-client-go 🍻🫃

Клиент для внутреннего GQL API платформы platform.21-school.ru. 
Предназначен для использования с [s21auto](https://github.com/s21toolkit/s21auto) для генерации запросов (`requests/`) из HAR логов платформы.

> [!IMPORTANT]
> Готовые версии автоклиента (со сгенерированным `requests/`) не публикуются на гитхабе.
> Если вам нужен простой доступ к платформе, используйте [s21introspector](https://github.com/s21toolkit/s21introspector) вместе с любым GQL клиентом для Go.
> Если же нужен именно автоклиент, его нужно склонить и собрать самостоятельно по [инструкции](#генерация-методов) ниже.

Пример использования:

```golang
client := s21client.New(
  s21client.DefaultAuth(
    os.Getenv("S21_USERNAME"),
    os.Getenv("S21_PASSWORD")
  )
)

user, err := client.R().GetCurrentUser(requests.GetCurrentUser_Variables{})
if err != nil {
  t.Error(err)
}

fmt.Println(user)
```

## Генерация методов

Методы клиента генерируются автоматически на основе запросов платформы к бекенду.

Для генерации запросов используется [s21auto](https://github.com/s21toolkit/s21auto):

```sh
s21auto client generate log.har -o s21client/requests
```

> Если какие-то методы не нужны, из папки requests можно удалить всё кроме [`context.go`](/requests/context.go).
