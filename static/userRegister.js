document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('registerForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        fetch('/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, email, password }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    document.getElementById('response').innerText = 'Registration successful!';
                    setTimeout(() => {
                        window.location.href = '/userLogin'; // Перенаправление на страницу логина
                    }, 2000);
                } else {
                    document.getElementById('response').innerText = data.message || 'Registration failed!';
                }
            })
            .catch(error => {
                document.getElementById('response').innerText = 'Error: ' + error.message;
            });
    });
});
