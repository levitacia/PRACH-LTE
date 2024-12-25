
# PRACH 

Упрощенное моделирование случайного доступа в LTE

Были введены следующие допущения

1. Все 64 преамбулы выделены под Connection-Based форму случайного доступа
2. Процедура случайного доступа занимает один фрейм
3. Начальное число мобильных устройств M можно задать в программе, также для каждого фрейма будут добавляться новые UE в случайном количестве от 0 до N, где N можно задать в программе
4. Коллизии, возникающие в случае, если 2 или более мобильных устройств случайно выберут одну и ту же преамбулу, будут решаться так: все UE, выбравшие одинаковую преамбулу, отклоняются и повторяют процедуру случайного доступа в следующем фрейме
5. Смоделирован только переход устройства из состояния IDLE в состояние CONNECTED. Мобильное устройство получает состояние CONNECTED, если за одну процедуру случайного доступа только он занял выбранную преамбулу, в таком случае он выходит из системы
6. Моделирующая команда работает в цикле по фреймам

