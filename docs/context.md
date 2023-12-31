# Контекст решения
<!-- Окружение системы (роли, участники, внешние системы) и связи системы с ним. Диаграмма контекста C4 и текстовое описание. 
-->
```plantuml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

Person(user, "Пользователь")
Person_Ext(courier, "Курьер")

System_Ext(webui, "Клинтский веб-сайт", "Взаимодействие пользователя через сайт")
Rel(user, webui, "Использует службу доставки с помощью сайта")

System(delivery_system, "Служба доставки", "Бэкенд для обработки пользователей, посылок и доставок")
Rel(webui, delivery_system, "API запросы")

System_Ext(logistic_system,"Система логистики", "Выполняет доставки")
Rel(delivery_system, logistic_system, "Передает доставки, которые нужно выполнить")
Rel(logistic_system, delivery_system, "Передает состояние доставки")
Rel(courier, logistic_system, "Выполняет доставку")

' System(user_service, "Сервис пользователей", "- Создание нового пользователя\n\n- Поиск пользователя по логину/по маске")
' System(item_service, "Сервис посылок", "- Создание посылки\n\n- Получение посылок")
' System(delivery_service, "Сервис доставки", "- Создание доставки\n\n- Получение информации о доставке по получателю/отправителю")

' Rel(webui, user_service, "API запросы")
' Rel(webui, item_service, "API запросы")
' Rel(webui, delivery_service, "API запросы")



@enduml
```
## Назначение систем
| Система            | Описание                                                                      |
| ------------------ | ----------------------------------------------------------------------------- |
| Клинтский веб-сайт | Веб-интерфейс, обеспечивающий взаимодействие между пользователем и сервисами. |
| Служба доставки    | система, отвечающая за работу с пользовательским фронтендом.                  |
| Система логистики  | Отвечает за логистику доставок. Внешняя система, которая их выполняет         |

