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