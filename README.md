[README.md](https://github.com/user-attachments/files/23018291/README.md)
# Практическое задание № 6 Борисов Денис Александрович ЭФМО-01-25

Цели:

1.	Понять, что такое ORM и чем удобен GORM..
2.	Научиться описывать модели Go-структурами и автоматически создавать таблицы (миграции через AutoMigrate).
3.	Освоить базовые связи: 1:N и M:N + выборки с Preload.
4.	Написать короткий REST (2–3 ручки) для проверки результата.

Выполнение практического задания.

1.	Установить и настроить PostgreSQL локально.
   Для выполнения задания на сервер был установлен postgresSQL и Go
  	
<img width="558" height="122" alt="Снимок экрана 2025-10-21 021739" src="https://github.com/user-attachments/assets/c5184ba4-e795-41e3-8b60-602c29dd7963" />

2.	Подключиться к БД из Go с помощью database/sql и драйвера PostgreSQL.
  Поссле был выполнен вход в пространство PostgresSQL

<img width="854" height="175" alt="Снимок экрана 2025-10-21 120634" src="https://github.com/user-attachments/assets/0746ff18-7fd0-46dc-9981-6a1e7395bca7" />


3.	Старт проекта и зависимости.
   Для выполнения задания был сформирован проект и установленны нужные зависимости

<img width="1068" height="610" alt="Снимок экрана 2025-10-21 120738" src="https://github.com/user-attachments/assets/bac30ed7-4dea-4b11-b129-ba81054381c4" />

  После была выполнена подготовка проекта для выполнения практики
  
<img width="703" height="535" alt="Снимок экрана 2025-10-21 121510" src="https://github.com/user-attachments/assets/9fbed819-8e10-4f78-a44e-b87233d77082" />

  После был написан код в файле для postgres.go, в котором будет происходит подключениие к БД и при помощи библиотеки GORM
  
<img width="686" height="860" alt="Снимок экрана 2025-10-21 124434" src="https://github.com/user-attachments/assets/6a5d55ae-0274-48d8-bcbe-15015a811c65" />

  После был написан код в файле для models.go, в котором будут происходить создания моделей
  
<img width="712" height="741" alt="Снимок экрана 2025-10-21 124531" src="https://github.com/user-attachments/assets/da1cdb5a-2a40-4517-b1fe-e0e90fd40a6d" />

  После был написан код в файле для router.go, в котором будут прописаны машруты, адрес http://denisb1995.fvds.ru
  
![Screenshot_20251021_195148_Termius](https://github.com/user-attachments/assets/da04e222-10c6-4972-a073-9a6e66b6a426)


  А так же написан код в файле для handlers.go
  
  Затем для проверки выполнения практики был осуществлен запуск проекта и проверка CRUD

Проверка на доступ к серверу

  <img width="493" height="131" alt="изображение" src="https://github.com/user-attachments/assets/cbbb128a-9f15-40d4-be0d-608f50a41b98" />

Проверка на создание пользователя

<img width="525" height="434" alt="изображение" src="https://github.com/user-attachments/assets/eee1bc5d-f523-4c2b-a62a-068015f39317" />

Вывод созданного пользователя

<img width="695" height="224" alt="изображение" src="https://github.com/user-attachments/assets/adbfc1be-0147-4314-99e1-f7f3015a19df" />

Удаление созданного пользователя 

<img width="484" height="427" alt="изображение" src="https://github.com/user-attachments/assets/520484d0-6f3a-4dd5-bf9c-084ec0318c6c" />

Создание заметки 

<img width="472" height="815" alt="изображение" src="https://github.com/user-attachments/assets/97a352ef-da79-4669-80e3-8d19b1ab464c" />


Вывод всех созданных заметок

<img width="527" height="637" alt="изображение" src="https://github.com/user-attachments/assets/acac76b6-19f5-4f63-86ce-0580ceeb3f85" />

Удаление заметки

<img width="424" height="598" alt="изображение" src="https://github.com/user-attachments/assets/7d945567-f90a-4a4a-9797-d811e23afcf0" />






