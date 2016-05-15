# yandex

Talks to the Yandex Machine Translation API https://tech.yandex.com/translate.

## Getting started
- Create a .JSON file containing your Yandex API key e.g
```json
{
  "api_key": "ThisIsMyYandexAPIKey"
}
```

- Next when invoking the API, set env variable `YANDEX\_API\_CREDENTIALS`
- When you invoke `yandex.New` it will auto-discover by the above mentioned
criteria.

Alternatively you can use `yandex.NewWithCredentials`.
