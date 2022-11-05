## Тестовое задание на позицию "Младший разработчик Go"

### Задание

Написать gRPC прокси-сервис для загрузки thumbnail’ов (превью видеоролика) c видео-
роликов YouTube. При повторном запросе на тот же видеоролик, сервис должен отдать
закэшированный ответ (можно использовать рантайм кэш, но будет плюсом если будет
использоваться временное хранилище, например sqlite). Так же предлагается написать
клиентскую часть как утилиту коммандной строки, которой передается в качестве па-
раметров ссылки на видеоролики. В коммандной утилите предусмотреть ключ --async,
который позволяет скачивать большое количество файлов асинхронно.

### Запуск приложения

1. Склонировать репозиторий с помощью команды **git clone https://github.com/Ig0rVItalevich/echelonTestTask.git**
2. Перейти в каталог проекта с помощью команды **cd ./echelonTestTask**
2. Находясь в рабочей дирректории проекта, запустить приложение командой **docker-compose up** 
3. Вызвать команду **docker ps | grep echelon_server** для просмотра имени запустившегося контейнера

### Запуск клиента

Клиент собирается вместе с сервером и запускается командой **docker exec -it *имя_контейнера* ./client [--async] *Список URL-адресов *youtube-видео* в кавычках***

### Примеры запуска

- Получение одного thumbnail

```docker exec -it echelon_server_1 ./client "https://www.youtube.com/watch?v=UC4T6t2s5vA&ab_channel=NatureRelaxationMusic"```

- Получение нескольких thumbnail'ов

```docker exec -it echelon_server_1 ./client "https://www.youtube.com/watch?v=UC4T6t2s5vA&ab_channel=NatureRelaxationMusic" "https://www.youtube.com/watch?v=RK1K2bCg4J8&t=132s&ab_channel=Balu-RelaxingNaturein4K"```

- Получение нескольких thumbnail'ов асинхронно

```docker exec -it echelon_server_1 ./client --async  "https://www.youtube.com/watch?v=UC4T6t2s5vA&ab_channel=NatureRelaxationMusic" "https://www.youtube.com/watch?v=RK1K2bCg4J8&t=132s&ab_channel=Balu-RelaxingNaturein4K" "https://www.youtube.com/watch?v=UV0mhY2Dxr0&ab_channel=RelaxationFilm"```





