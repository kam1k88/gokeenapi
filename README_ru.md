## Утилита для добавления\удаления статичных маршрутов в роутерах Кинетик через REST API

---

#### Настройка REST API на роутере
* Прежде чем появится возможность использовать эту утилиту, Вы должны настроить REST API на роутере. Как это сделать можно почитать в официальной документации [здесь](https://help.keenetic.com/hc/ru/articles/11282223272092-%D0%9F%D1%80%D0%B8%D0%BC%D0%B5%D0%BD%D0%B5%D0%BD%D0%B8%D0%B5-%D0%BC%D0%B5%D1%82%D0%BE%D0%B4%D0%BE%D0%B2-API-%D0%BF%D0%BE%D1%81%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B2%D0%BE%D0%BC-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81%D0%B0-HTTP-Proxy). Не рекомендуется открывать REST API в интернет
* Ссылка на Keenetic REST API в конфигурации должна содержать в себе путь /rci, например через флаг`--api https://super-api.my-super-keenetic.keenetic.pro/rci` 
---

#### Что утилита уже умеет:
* Выводить список интерфейсов которые уже добавлены в роутер - для удобного поиска ID интерфейса для которого нужно добавить\удалить маршруты
* Удалять статичные маршруты только для указанного интерфейса. В Веб-конфигураторе роутера на текущий момент можно только удалить все созданные статичные маршруты для всех интерфейсов сразу
* Добавлять\обновлять статичные маршруты для указанного интерфейса из bat файлов с диска
* Добавлять\обновлять статичные маршруты для указанного интерфейса из ссылок, которые ведут на bat файл (например [отсюда](https://iplist.opencck.org/?format=bat&data=cidr4&site=youtube.com))
---

#### Конфигурация

Утилиту можно сконфигурировать несколькими путями:
* Через конфигурационный файл YAML
* Через переменные окружения
* Через файл с переменными окружениями которые надо загрузить
* Через флаги в командной строке

Все варианты можно совмещать - например логин\пароль и API URL можно хранить в переменных окружения, а список файлов откуда нужно добавить маршруты можно добавить в yaml конфиг файл или передать через флаги

---

#### Примеры использования

Самый простой способ начать пользоваться утилитой через docker контейнеры

---

#### Docker 

* Посмотреть интерфейсы на роутере - передача логин\пароля\апи через флаги
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:latest"
docker pull "${GOKEENAPI_IMAGE}"
docker run --rm -ti "${GOKEENAPI_IMAGE}" show-interfaces --api https://api.mykeenetic.keenetic.pro/rci --login admin --password admin
```

* Посмотреть интерфейсы на роутере - передача логин\пароля\апи через переменные окружения
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:latest"
docker run --rm -ti -e GOKEENAPI_API="https://api.mykeenetic.keenetic.pro/rci" -e GOKEENAPI_LOGIN="admin" -e OKEENAPI_PASSWORD="admin" "${GOKEENAPI_IMAGE}" show-interfaces
```

* Посмотреть интерфейсы на роутере - передача логин\пароля\апи через файл с переменными окружениями
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:latest"
touch .gokeenapienv
echo -e "GOKEENAPI_API=https://api.mykeenetic.keenetic.pro/rci\n" >> .gokeenapienv
echo -e "GOKEENAPI_LOGIN=admin\n" >> .gokeenapienv
echo -e "GOKEENAPI_PASSWORD=admin\n" >> .gokeenapienv
docker run --rm -ti -v "$(pwd)/.gokeenapienv":"/gokeenapi/.gokeenapienv" "${GOKEENAPI_IMAGE}" show-interfaces
```

* Посмотреть интерфейсы на роутере - передача логин\пароля\апи через YAML конфиг файл
```shell
export GOKEENAPI_IMAGE="noksa/gokeenapi:latest"
docker run --rm -ti -v "$(pwd)/config_example.yaml":"/gokeenapi/config.yaml" "${GOKEENAPI_IMAGE}" show-interfaces --config "/gokeenapi/config.yaml"
```