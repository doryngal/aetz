function updateTimeLeftForAllLots() {
    var lotElements = document.querySelectorAll('[id^="time-left-"]');

    lotElements.forEach(function(lotElement) {
        var endDate = new Date(lotElement.getAttribute('data-enddate'));
        var now = new Date();
        var timeDiff = endDate - now;

        var parentDiv = document.getElementById('lot-' + lotElement.id.split('-')[2]);

        if (timeDiff <= 0) {
            lotElement.textContent = "Время истекло";
            parentDiv.style.backgroundColor = "gray";
        } else {
            // Вычисляем количество оставшихся дней
            var days = Math.floor(timeDiff / (1000 * 60 * 60 * 24));
            lotElement.textContent = `Осталось ${days} ${days === 1 ? 'день' : 'дня'}`;
            parentDiv.style.backgroundColor = "green";
        }
    });
}

// Первоначальный вызов функции для отображения времени
updateTimeLeftForAllLots();

function updateCountdown() {
    // Получаем дату окончания из input
    const endDateInput = document.getElementById('endDate').value;
    const endDate = new Date(endDateInput);

    const currentTime = new Date();
    const difference = endDate - currentTime;

    // Если время вышло
    if (difference <= 0) {
        // Показываем, что тендер завершен
        document.querySelector('.lot__time span').textContent = "Тендер завершен!";

        // Обновляем значения в HTML, чтобы показывать нули
        document.getElementById('days').querySelector('.countdown__value').textContent = "00";
        document.getElementById('hours').querySelector('.countdown__value').textContent = "00";
        document.getElementById('minutes').querySelector('.countdown__value').textContent = "00";
        document.getElementById('seconds').querySelector('.countdown__value').textContent = "00";

        clearInterval(intervalId);
        return;
    }

    // Вычисляем дни, часы, минуты и секунды
    const days = Math.floor(difference / (1000 * 60 * 60 * 24));
    const hours = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((difference % (1000 * 60)) / 1000);

    // Обновляем значения в HTML
    document.getElementById('days').querySelector('.countdown__value').textContent = String(days).padStart(2, '0');
    document.getElementById('hours').querySelector('.countdown__value').textContent = String(hours).padStart(2, '0');
    document.getElementById('minutes').querySelector('.countdown__value').textContent = String(minutes).padStart(2, '0');
    document.getElementById('seconds').querySelector('.countdown__value').textContent = String(seconds).padStart(2, '0');
}

// Обновляем обратный отсчет каждую секунду
const intervalId = setInterval(updateCountdown, 1000);

