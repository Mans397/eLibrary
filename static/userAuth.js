document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm');

    form.addEventListener('submit', function (event) {
        event.preventDefault();

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
                    if (data.role === 'admin') {
                        // Показываем админские кнопки
                        document.getElementById('adminActions').style.display = 'block';
                        document.getElementById('response').innerText = 'Welcome, Admin!';
                    } else if (data.role === 'user') {
                        // Перенаправляем пользователя на страницу с книгами
                        document.getElementById('response').innerText = 'Login successful! Redirecting...';
                        setTimeout(() => {
                            window.location.href = data.redirect || '/books';
                        }, 2000);
                    }
                } else {
                    document.getElementById('response').innerText = data.message || 'Login failed!';
                }
            })
            .catch(error => {
                document.getElementById('response').innerText = 'Error: ' + error.message;
            });
    });
});
