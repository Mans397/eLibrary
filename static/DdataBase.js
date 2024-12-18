document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('DB-Create').addEventListener('click', function () {
        // Получаем данные из формы
        const name = document.getElementById('nameInput').value;
        const email = document.getElementById('emailInput').value;

        // Отправляем запрос на сервер
        fetch('/db/createUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({name: name, email: email}), // Отправляем данные в формате JSON
        })
            .then(response => response.json())
            .then(data => {
                // Обрабатываем ответ от сервера
                document.getElementById('output').innerText = `${data.message}`;

            })
            .catch(error => {
                document.getElementById('output').innerText = 'Error: ' + error.message;
            });
    });
})


document.addEventListener('DOMContentLoaded', () => {
document.getElementById("DB-Read").addEventListener("click", function () {
    // Получаем значение email из текстового поля
    const email = document.getElementById("emailInput").value;

    if (!email) {
        alert("Введите email!");
        return;
    }

    // Формируем URL с параметром email
    const url = `/db/readUser?email=${encodeURIComponent(email)}`;

    // Отправляем GET-запрос
    fetch(url, {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Ошибка: ${response.status}`);
            }
            return response.json();
        })
        .then((data) => {
            // Выводим полученные данные
            const output = document.getElementById("output");
            output.innerHTML = `
                <p>User Name: ${data.name}</p>
                <p>User Email: ${data.email}</p>
            `;
        })
        .catch((error) => {
            // Обработка ошибок
            const output = document.getElementById("output");
            output.innerHTML = `<p style="color: red;">${error.message}</p>`;
        });
});})