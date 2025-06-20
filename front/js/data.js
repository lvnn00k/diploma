function table(data) {
    document.querySelector("tbody").innerHTML = ''; 
    document.querySelector(".breadcrumbs").innerHTML = `<a href="/1">Главная</a>`;
    var link = "";
    if (data['path'].length > 1) {
        data['path'].forEach(function(item) {
            folderName = item['name'];
            folderId = item['id'];
            if (folderName != "") {
                document.querySelector(".breadcrumbs").insertAdjacentHTML('beforeend', 
                    `<div class='division'></div>
                    <a href="/${folderId}">${folderName}</a>
                    `
                );
            }
            link += `${folderName}/`;
        });
    }

    if (document.documentElement.clientWidth <= 415 && document.querySelector(".breadcrumbs").clientWidth >= 330) {
        document.querySelector(".breadcrumbs").innerHTML = `<a href="/1">Главная</a>
            <div class='division'></div>
            <a>...</a>
            <div class='division'></div>
            <a href="/${data['path'][data['path'].length - 1]['id']}">${data['path'][data['path'].length - 1]['name']}</a>
        `;
    }
    if (data['folders'] === null && data['files'] === null) {
        document.querySelector("tbody").insertAdjacentHTML('beforeend',
            `<tr><td colspan="4" text-align="center">Папка пуста</td></tr>`
        );
    }
    if (data['folders'] != null) {
        data['folders'].forEach(function(item) {
            fileName = item['name'];
            id = item['id'];
            date = item['date_create'];
            document.querySelector("tbody").insertAdjacentHTML('beforeend',
                `
                <tr data-type="folder" data-id="`+id+`">
                    <td></td>
                    <td class="table-name">
                        <div class='name'><a class="dir-link" href='/`+id+`'>`+fileName+`</a></div>
                    </td>
                    <td>`+date+`</td>
                    <td>
                        <div class="file-icon-block">
                            <div class="more-info"></div>
                        </div>
                    </td>
                </tr>
                `
            );
        });
    }
    if (data['files'] != null) {
        data['files'].forEach(function(item) {
            fileName = item['name'];
            id = item['id'];
            date = item['date_create'];
            document.querySelector("tbody").insertAdjacentHTML('beforeend',
                `
                <tr data-type="file" data-id="`+id+`">
                    <td>
                        <div class="file-icon-block"">
                            <div class="file-icon"></div>
                        </div>
                    </td>
                    <td class="table-name file-info" data-link="`+link+`">
                        <div class='name'>`+fileName+`</div>
                    </td>
                    <td>`+date+`</td>
                    <td>
                        <div class="file-icon-block">
                            <div class="more-info"></div>
                        </div>
                    </td>
                </tr>
                    `
            );
        });
    }
    
}


function load() {
    url = "/api/dirs/d"+window.location.pathname;

    fetch(url)
        .then(response => { 
            if (response.status === 401) {
                refreshAccessToken();
                return load();
            }
            return response.json();
        })
        .then(data => {
            table(data);
            showImage();
            MoreInfo();
        })
        .catch(error => {
            console.log(error);
        });
}


function refreshAccessToken() {
    fetch('/api/auth/refresh-token')
        .then(response => {
            if (response.status !== 200) {
                window.location.href = '/login'; 
                return;
            }
        })
        .catch(error => {
            console.log(error);
        });
}

load();


