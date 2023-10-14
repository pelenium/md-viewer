const fileUploadForm = document.getElementById('fileUploadForm');
const fileInput = document.getElementById('fileInput');

fileUploadForm.addEventListener('submit', function (event) {
    event.preventDefault();

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        const fileReader = new FileReader();

        fileReader.onload = function (event) {
            const fileData = event.target.result;

            sendFileToServer(fileData);
        };

        fileReader.readAsText(selectedFile);
    }
});

function sendFileToServer() {
    const fileInput = document.getElementById("fileInput");
    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        const reader = new FileReader();

        reader.onload = function (e) {
            const fileData = e.target.result;
            console.log(CryptoJS.MD5(fileData).toString());

            fetch("/render", {
                method: "POST",
                body: JSON.stringify({
                    data: fileData,
                    name: CryptoJS.MD5(fileData).toString()
                })
            })
                .then(response => {
                    if (response.ok) {
                        status.textContent = "Файл успешно отправлен на сервер.";
                        setTimeout(() => {
                            console.log(CryptoJS.MD5(fileData).toString());
                            var hash = CryptoJS.MD5(fileData).toString();
                            window.location.replace(`/${hash}`)
                        }, 1500);

                    } else {
                        status.textContent = "Ошибка при отправке файла на сервер.";
                    }
                })
                .catch(error => {
                    console.error("Ошибка при отправке файла:", error);
                    status.textContent = "Произошла ошибка.";
                });
        };

        reader.readAsText(selectedFile);
    }
}