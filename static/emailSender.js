document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('emailForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const message = document.getElementById('message').value;
        const attachment = document.getElementById('attachment').files[0];

        const formData = new FormData();
        formData.append('message', message);
        if (attachment) {
            formData.append('attachment', attachment);

        }


        document.getElementById('response').innerText = "";

        fetch('/admin/sendEmail', {
            method: 'POST',
            body: formData,
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