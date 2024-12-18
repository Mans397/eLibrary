document.addEventListener('DOMContentLoaded', () => {
    // Обработчик для кнопки Create user
    document.getElementById('DB-Create').addEventListener('click', function () {
        const name = document.getElementById('nameInput').value.trim();
        const email = document.getElementById('emailInput').value.trim();

        if (!name || !email) {
            alert('Введите имя и email!');
            return;
        }
        console.log('Кнопка "Create User" нажата');
        console.log('Email:', email);

        fetch('/db/createUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: name, email: email }),
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('output').innerText = `${data.message}`;
            })
            .catch(error => {
                document.getElementById('output').innerText = 'Error: ' + error.message;
            });
    });

    // Обработчик для кнопки Read user
    document.getElementById('DB-Read').addEventListener('click', function () {
        const email = document.getElementById('emailInput').value.trim();

        console.log('Кнопка "Read User" нажата');
        console.log('Email:', email);

        if (!email) {
            alert('Введите email!');
            return;
        }

        const url = `/db/readUser?email=${encodeURIComponent(email)}`;

        fetch(url, { method: 'GET' })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Ошибка: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                document.getElementById('output').innerHTML = `
                    <p>User Name: ${data.name}</p>
                    <p>User Email: ${data.email}</p>
                `;
            })
            .catch(error => {
                document.getElementById('output').innerHTML = `
                    <p style="color: red;">${error.message}</p>
                `;
            });
    });
});