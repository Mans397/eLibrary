document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('emailForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();
        const message = document.getElementById('message').value;
        document.getElementById('response').innerText = "";
        fetch('/admin/sendEmail', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ message }),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Failed to send message');
                }
                return response.json();
            })
            .then((result) => {
                document.getElementById('response').innerText = result.message;
            })
            .catch(() => {
                document.getElementById('response').innerText = 'Error sending message';
            });
    });
});