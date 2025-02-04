document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const otpForm = document.getElementById('otpForm');
    const responseDiv = document.getElementById('response');
    const emailField = document.getElementById('email');

    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const email = emailField.value.trim();
        const password = document.getElementById('password').value.trim();

        console.log("Attempting login with email:", email);

        fetch('/auth/userLogin', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        })
            .then(response => response.json())
            .then(data => {
                console.log("Server response for login:", data);

                if (data.status === 'success') {
                    if (data.redirect) {
                        // ✅ Администратор: редиректим в /admin
                        window.location.href = data.redirect;
                    } else {
                        // ✅ Обычный пользователь: показываем форму OTP
                        responseDiv.innerText = 'OTP sent to your email. Please enter OTP.';
                        setTimeout(() => {
                            loginForm.style.display = 'none';
                            otpForm.removeAttribute('hidden');
                            otpForm.style.display = 'block';
                        }, 500);
                    }
                } else {
                    responseDiv.innerText = data.message || 'Login failed!';
                }
            })
            .catch(error => {
                console.error("Error during login:", error);
                responseDiv.innerText = 'Error: ' + error.message;
            });
    });

    otpForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const email = emailField.value.trim();
        const otp = document.getElementById('otp').value.trim();

        console.log("Attempting OTP verification for email:", email, "with OTP:", otp);

        fetch('/auth/verifyOTP', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, otp }),
        })
            .then(response => response.json())
            .then(data => {
                console.log("Server response for OTP verification:", data);

                if (data.status === 'success') {
                    responseDiv.innerText = 'Login successful! Redirecting...';

                    // ✅ Сохраняем токен в localStorage
                    localStorage.setItem('token', data.token);

                    // ✅ Мгновенно редиректим
                    window.location.href = '/books';
                } else {
                    responseDiv.innerText = 'Invalid OTP!';
                }
            })
            .catch(error => {
                console.error("Error during OTP verification:", error);
                responseDiv.innerText = 'Error: ' + error.message;
            });
    });
});