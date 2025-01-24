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
            body: JSON.stringify({name: name, email: email}),
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
        const name = document.getElementById('nameInput').value.trim();

        console.log('Кнопка "Read User" нажата');
        console.log('Email:', email, name);

        const url = `/db/readUser?email=${encodeURIComponent(email)}&name=${encodeURIComponent(name)}`;
        console.log(url)

        fetch(url, {method: 'GET'})
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Ошибка: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (Array.isArray(data) && data.length > 0) {
                    document.getElementById('output').innerHTML = data.map(user => `
                <p>User Name: ${user.name}</p>
                <p>User Email: ${user.email}</p>
                <hr>
            `).join('');

                } else {
                    document.getElementById('output').innerHTML = `
                    <p>User Name: ${data.name}</p>
                    <p>User Email: ${data.email}</p>
                `;
                }
            })
            .catch(error => {
                document.getElementById('output').innerHTML = `
                    <p style="color: red;">${error.message}</p>
                `;
            });
    });

    document.getElementById('DB-Update').addEventListener('click', function () {
        const name = document.getElementById('nameInput').value.trim();
        const email = document.getElementById('emailInput').value.trim();

        console.log('Кнопка "Update User" нажата');

        if (!name || !email) {
            alert('Введите имя и email!');
            return;
        }

        fetch('/db/updateUser', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({email: email, name: name}),
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('output').innerText = data.message || 'User updated successfully';
            })
            .catch(error => {
                document.getElementById('output').innerText = 'Error: ' + error.message;
            });
    });

    document.getElementById('DB-Delete').addEventListener('click', function () {
        const email = document.getElementById('emailInput').value.trim();

        if (!email) {
            alert('Введите email!');
            return;
        }

        console.log('Кнопка "Delete User" нажата');

        fetch('/db/deleteUser', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({email: email}),
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('output').innerText = data.message || 'User deleted successfully';
            })
            .catch(error => {
                document.getElementById('output').innerText = 'Error: ' + error.message;
            });
    });

    // Обработчик для логина администратора
    document.getElementById('DB-AdminLogin').addEventListener('click', function () {
        const email = document.getElementById('emailInput').value.trim();
        const password = document.getElementById('nameInput').value.trim();

        if (!email || !password) {
            alert('Введите email и пароль!');
            return;
        }

        console.log('Кнопка "Admin Login" нажата');
        console.log('Email:', email);

        fetch('/auth/adminLogin', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email: email, password: password }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    alert(data.message);
                    window.location.href = data.redirect; // Перенаправление на страницу админа
                } else {
                    alert('Login failed: ' + (data.message || 'Invalid credentials'));
                }
            })
            .catch(error => {
                alert('Error: ' + error.message);
            });
    });

    // Обработчик для кнопки поиска пользователей
    document.getElementById('DB-Search').addEventListener('click', () => {
        const name = document.getElementById('searchInput').value;

        if (!name) {
            alert("Введите имя для поиска");
            return;
        }

        document.addEventListener("DOMContentLoaded", function () {
            document.getElementById("DB-GetAll").addEventListener("click", getAllUsers);
        });

        function getAllUsers() {
            fetch("/db/getAllUsers", { method: "GET" })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.Status === "fail") {
                        alert(`Error: ${data.Message}`);
                    } else {
                        displayUsers(data);
                    }
                })
                .catch(error => {
                    console.error("Error fetching users:", error);
                });
        }

        function displayUsers(users) {
            const output = document.getElementById("output");
            output.innerHTML = "";

            const table = document.createElement("table");
            table.style.border = "1px solid black";

            const headerRow = document.createElement("tr");
            headerRow.innerHTML = `
        <th>ID</th>
        <th>Name</th>
        <th>Email</th>
    `;
            table.appendChild(headerRow);

            users.forEach(user => {
                const row = document.createElement("tr");
                row.innerHTML = `
            <td>${user.ID}</td>
            <td>${user.Name}</td>
            <td>${user.Email}</td>
        `;
                table.appendChild(row);
            });

            output.appendChild(table);
        }
    });
});
