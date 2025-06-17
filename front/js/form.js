function AddFolder() {
    const formData = new FormData(document.querySelector('.form')); 
    const name = formData.get('name');
    const path = Number(window.location.pathname.substring(1));

    const request = {
        name: name,
        parent_id: path
    };

    fetch('/api/dirs/add/folder', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
    })
        .then(response => { 
            if (response.status === 401) {
                refreshAccessToken();
                return AddFolder();
            }
            if (response.status === 409) {
                document.querySelector(".div-form").querySelector("small").textContent = "Папка уже существует";
                return;
            }
            if (response.status === 500 || response.status === 400) {
                document.querySelector(".div-form").querySelector("small").textContent = "Неудалось создать папку";
                return;
            }
            document.querySelector(".div-form").remove();
            document.querySelector('[data-form="folder"]').setAttribute('data-active', 'false');
            load(); 
        })
        .catch(error => {
            console.log(error);
        });
}

function AddFile() {
    const formData = new FormData(document.querySelector('.form'));
    formData.append("parent_id", Number(window.location.pathname.substring(1)))
    fetch('/api/dirs/add/file', {
        method: 'POST',
        body: formData
    })
        .then(response => { 
            if (response.status === 401) {
                refreshAccessToken();
                return AddFile();
            }
            if (response.status === 409) {
                document.querySelector(".div-form").querySelector("small").textContent = "Файл уже существует";
                return;
            }
            if (response.status === 413) {
                document.querySelector(".div-form").querySelector("small").textContent = "Файл слишком большой";
                return;
            }
            if (response.status === 500 || response.status === 400) {
                document.querySelector(".div-form").querySelector("small").textContent = "Не удалось загрузить файл";
                return;
            }
            document.querySelector(".div-form").remove();
            document.querySelector('[data-form="file"]').setAttribute('data-active', 'false');
            load(); 
        })
        .catch(error => {
            console.log(error);
        });
        
}

function DeleteItem(block) {
    tr = block.closest("tr");
    type = tr.getAttribute("data-type");
    id = tr.getAttribute("data-id");

    fetch(`/api/dirs/delete/${type}/${id}`, {
        method: 'DELETE'
    })
        .then(response => { 
            if (response.status === 401) {
                refreshAccessToken();
                return DeleteItem(block);
            }
            load();
            document.getElementById("hidden").remove();
            document.getElementsByClassName("div-form")[0].remove();
        })
        .catch(error => {
            console.log(error);
        });
}
 

document.addEventListener('DOMContentLoaded', function() {
    const container = document.querySelector('body');
    if (container) {
        container.addEventListener("change", function(event) {
            if (event.target && event.target.matches('[data-input="file-file"]')) {
                fileName = document.getElementsByName("name")[0].value.trim();
                button = document.querySelector('[data-btn="new"]');
                if (fileName === '') {
                    button.setAttribute('data-active', 'true');
                } else {
                    file = document.getElementsByName("file")[0];
                    if (file.files.length === 0) {
                        button.setAttribute('data-active', 'true');
                    } else {
                        button.setAttribute('data-active', 'false');
                    }
                }
            }
        })

        container.addEventListener("input", function(event) {
            if (event.target && event.target.matches('[data-input="folder-name"]')) {
                folder = document.getElementsByName("name")[0].value.trim();
                button = document.querySelector('[data-btn="new"]');
                if (folder === '') {
                    button.setAttribute('data-active', 'true');
                } else {
                    button.setAttribute('data-active', 'false');
                }
            }
            if (event.target && event.target.matches('[data-input="file-name"]')) {
                fileName = document.getElementsByName("name")[0].value.trim();
                button = document.querySelector('[data-btn="new"]');
                if (fileName === '') {
                    button.setAttribute('data-active', 'true');
                } else {
                    file = document.getElementsByName("file")[0];
                    if (file.files.length === 0) {
                        button.setAttribute('data-active', 'true');
                    } else {
                        button.setAttribute('data-active', 'false');
                    }
                }
            }
            
        });
        container.addEventListener('submit', function(event) {
            if (event.target && event.target.matches('.folder')) {
                event.preventDefault(); 
                active = document.querySelector('[data-btn="new"]').getAttribute("data-active");
                if (active === "false") {
                    AddFolder();
                }
            }

            if (event.target && event.target.matches('.file')) {
                event.preventDefault();
                active = document.querySelector('[data-btn="new"]').getAttribute("data-active");
                if (active === "false") {
                    AddFile();
                }   
            }
        });
        container.addEventListener("click", function(event) {
            if (event.target && event.target.matches('.delete')) {
                block = event.target;
                let delblock = document.createElement("div");
                delblock.className = "div-form";
                delblock.innerHTML = 
                    `<div id="back">
                        <div class="form-container">
                            <div class="delete-block">
                                <div >
                                    <h3>
                                        Вы действительно хотите удалить этот элемент?
                                    </h3>
                                </div>
                                <div class="foot">
                                    <a href="#" id="cancelBtn">Отмена</a>
                                    <button class="btn-delete">Удалить</button>
                                </div>
                            </div>
                        </div>
                    </div>
                    `;

                    document.body.append(delblock);

                    cancelBtn.addEventListener("click", function() {
                        delblock.remove();
                    });
                    back.addEventListener("click", function(event) {
                        if (event.target === back) {
                            delblock.remove();
                        }
                    });

                    document.getElementsByClassName("btn-delete")[0].addEventListener("click", function() {
                        DeleteItem(block);
                    })
            }
        })
    } 
});

document.querySelector(".logout").addEventListener("click", function(event) {
    event.preventDefault;
    fetch('api/auth/logout')
        .then(response => {
            if (response.status === 200) {
                window.location.href = "/login";
            }
        })
        .then(error => {
            console.log(error);
        })
})