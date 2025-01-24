document.addEventListener('DOMContentLoaded', () => {
    const logoutButton = document.getElementById('logoutButton');

    logoutButton.addEventListener('click', () => {
        fetch('/logout', {
            method: 'POST',
        })
            .then(response => {
                if (response.ok) {
                    window.location.href = '/'; // Перенаправляем на главную страницу
                } else {
                    alert('Failed to log out');
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
    });
});
