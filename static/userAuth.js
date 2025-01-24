document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm');

    form.addEventListener('submit', function (event) {
        event.preventDefault(); // Предотвращаем перезагрузку страницы

        const email = document.getElementById('email').value.trim();
        const password = document.getElementById('password').value.trim();

        fetch('/auth/userLogin', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    // Показать сообщение об успешном входе
                    document.getElementById('response').innerText = 'Login successful!';

                    // Перенаправить пользователя в зависимости от роли
                    setTimeout(() => {
                        window.location.href = data.redirect; // Сервер возвращает /admin или /books
                    }, 2000);
                } else {
                    // Если вход не удался, показать сообщение
                    document.getElementById('response').innerText = data.message || 'Login failed!';
                }
            })
            .catch(error => {
                // Обработка ошибки
                document.getElementById('response').innerText = 'Error: ' + error.message;
            });
    });
});
