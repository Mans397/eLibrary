document.addEventListener('DOMContentLoaded', () => {
    const button = document.getElementById('fetchDataButton');
    if (button) {
        button.addEventListener('click', async () => {
            try {
                // Отправляем GET-запрос на /json
                const response = await fetch('/json');
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }

                // Получаем JSON-ответ
                const data = await response.json();

                // Отображаем данные в div
                document.getElementById('result').innerText = JSON.stringify(data, null, 2);
            } catch (error) {
                console.error('Error fetching JSON:', error);
            }
        });
    }
});



document.addEventListener('DOMContentLoaded', () => {
    // Обработчик для кнопки "Send POST"
    const sendButton = document.getElementById('sendButton');
    if (sendButton) {
        sendButton.addEventListener('click', async () => {
            // Получаем значение из поля ввода
            const message = document.getElementById('textInput').value;

            if (message) {
                try {
                    // Отправляем POST-запрос с данными в формате JSON
                    const response = await fetch('/', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({ message: message })
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }

                    // Получаем ответ от сервера
                    const data = await response.json();

                    // Находим элемент с id "result" и выводим в него ответ
                    document.getElementById('result').innerText = JSON.stringify(data, null, 2);
                } catch (error) {
                    console.error('Error sending message:', error);
                    document.getElementById('result').innerText = 'Error sending message: ' + error.message;
                }
            } else {
                alert('Введите текст перед отправкой!');
            }
        });
    }
})