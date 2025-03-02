const API_URL = "http://localhost:8080/api/v1";
let currentCurrency = "";
let currentBase = "";

document.addEventListener("DOMContentLoaded", () => {
    loadRates();
    setInterval(loadRates, 50000);

    document.getElementById("addRate").addEventListener("click", async () => {
        const currencyPair = document.getElementById("currencyPair").value.trim();
        if (!currencyPair) {
            showNotification("Введите валютную пару (например, EUR/USD)", true);
            return;
        }

        const response = await addRate(currencyPair);
        if (response) {
            showNotification("Курс успешно добавлен!");
            loadRates();
        }
    });

    document.getElementById("searchRate").addEventListener("click", async () => {
        const id = document.getElementById("searchId").value.trim();
        if (!id) {
            showNotification("Введите ID курса", true);
            return;
        }

        const rate = await getRateById(id);
        document.getElementById("searchResult").innerText = rate
            ? `Курс: ${rate.currency}/${rate.base} = ${rate.rate}`
            : "Курс не найден";
    });

    document.getElementById("saveEdit").addEventListener("click", async () => {
        const newRate = document.getElementById("editRate").value.trim();
        updateRate(currentCurrency, currentBase, newRate);
        closeModal();
    });

    document.getElementById("cancelEdit").addEventListener("click", () => {
        closeModal();
    });
});

// Функция для получения данных с API
async function fetchData(url, options = {}) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) throw new Error(`Ошибка: ${response.statusText}`);
        return await response.json();
    } catch (error) {
        console.error("Ошибка запроса:", error);
        showNotification("Ошибка загрузки данных", true);
    }
}

// Получение последнего курса валют
async function getLatestRate(currencyPair) {
    return await fetchData(`${API_URL}/last?rate=${currencyPair}`);
}

// Добавление валютного курса
async function addRate(currencyPair) {
    return await fetchData(`${API_URL}/?rate=${currencyPair}`, { method: "PUT" });
}

// Удаление курса по ID
// Удаление курса по валютной паре
async function deleteRate(currency, base) {
    if (!currency || !base) {
        showNotification("Ошибка: Не указана валютная пара", true);
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/delete/${currency}/${base}`, {
            method: "DELETE",
        });

        if (!response.ok) throw new Error("Ошибка при удалении");
        
        showNotification(`Курс ${currency}/${base} удалён`);
        loadRates();
    } catch (error) {
        console.error("Ошибка удаления:", error);
        showNotification(error.message, true);
    }
}

// Открытие модального окна
function openModal(currency, base, oldRate) {
    currentCurrency = currency;
    currentBase = base;
    document.getElementById("editRate").value = oldRate;
    const modal = document.getElementById("editModal");
    modal.classList.add("show"); // Добавляем класс show
}

function closeModal() {
    const modal = document.getElementById("editModal");
    modal.classList.remove("show"); // Удаляем класс show
}


// Получение курса по ID
async function getRateById(id) {
    return await fetchData(`${API_URL}/by-id/${id}`);
}


async function updateRate(currency, base, newRate) {
    if (!newRate || isNaN(newRate) || parseFloat(newRate) <= 0) {
        showNotification("Некорректное значение курса!", true);
        return;
    }

    try {
        const response = await fetch(`${API_URL}/update?currency=${currency}&base=${base}&rate=${newRate}`, {
            method: "PATCH",
        });

        if (!response.ok) {
            throw new Error("Ошибка при обновлении курса");
        }

        showNotification(`Курс ${currency}/${base} обновлён`);
        loadRates(); // Обновляем таблицу после обновления
    } catch (error) {
        console.error("Ошибка обновления курса:", error);
    }
}

// Загрузка курсов в таблицу
async function loadRates() {
    try {
        const response = await fetch("http://localhost:8080/api/v1/all-last");
        if (!response.ok) {
            throw new Error("Ошибка загрузки данных");
        }

        const data = await response.json();
        const ratesTable = document.getElementById("ratesTableBody");
        ratesTable.innerHTML = "";

        data.forEach(rate => {
            let arrow = "";
            if (rate.changePct > 0) {
                arrow = `<span style="color: green; font-size: 16px;">▲ +${rate.changePct.toFixed(2)}%</span>`;
            } else if (rate.changePct < 0) {
                arrow = `<span style="color: red; font-size: 16px;">▼ ${rate.changePct.toFixed(2)}%</span>`;
            }
            const row = document.createElement("tr");
            row.innerHTML = `
            <td>${rate.currency}</td>
            <td>${rate.base}</td>
            <td>${rate.rate} ${arrow}</td>
            <td>${new Date(rate.updateDt).toLocaleString()}</td>
            <td>
                <button class="icon-button" onclick="openModal('${rate.currency}', '${rate.base}', ${rate.rate})" aria-label="Редактировать">
                    <svg class="icon edit-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168l10-10zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325z"/>
                    </svg>
                </button>
                <button class="icon-button" onclick="deleteRate('${rate.currency}', '${rate.base}')" aria-label="Удалить">
                    <svg class="icon trash-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                        <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                    </svg>
                </button>
            </td>
        `;
            ratesTable.appendChild(row);
        });

    } catch (error) {
        console.error("Ошибка загрузки данных:", error);
    }
}



// Показ уведомлений
function showNotification(message, isError = false) {
    const notification = document.getElementById("notification");
    notification.textContent = message;
    
    // Сбрасываем анимацию
    notification.classList.remove('show', 'error');
    void notification.offsetWidth; // Триггер перерисовки
    
    notification.classList.add('show');
    if(isError) notification.classList.add('error');

    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}
