function objt(data) {
    html = '';
    if (data['folders'] != null) {
        data['folders'].forEach(function(item) {
            fileName = item['folder']['name'];
            link_name = item['path'];
            if (document.documentElement.clientWidth <= 415) {
                link_objs = link_name.split("/");
                if (link_objs.length > 3) {
                    link_name = "/"+link_objs[1]+"/.../"+link_objs[link_objs.length - 2];
                } else if (link_objs.length == 3) {
                    link_name = "/"+link_objs[1];
                } else {
                    link_name = "/";
                }
            } 
            id = item['folder']['id'];
            html +=
                `<div class="result-obj">
                    <div class="obj-type"></div>
                    <div class="obj-name">
                        <a href='/`+id+`'>`+fileName+`</a>
                        <div class="obj-path">`+link_name+`</div>
                    </div>
                </div>
                `;
        });
    }
    if (data['files'] != null) {
        data['files'].forEach(function(item) {
            fileName = item['file']['name'];
            link = item['path'];
            link_name = link;
            if (document.documentElement.clientWidth <= 415) {
                link_objs = link_name.split("/");
                if (link_objs.length > 2) {
                    link_name = "/"+link_objs[1]+"/.../"+link_objs[link_objs.length - 1];
                } else if (link_objs.length == 2) {
                    link_name = "/"+link_objs[1];
                } else {
                    link_name = "/";
                }
            } 
            html += 
                `<div class="result-obj">
                    <div class="obj-type"><div class="file-icon"></div></div>
                    <div class="obj-name search-file-info" data-link="`+link+`">
                        <div class='name'>`+fileName+`</div>
                        <div class="obj-path">`+link_name+`</div>
                    </div>
                </div>
                `;
        });
    }
    if (html === '') {
        html += `<div class="res-null">Совпадений не найдено</div>`;
    }
    return html;
}

function search(input, query) {
    var url = "/api/dirs/search/" + query;
    fetch(url)
        .then(response => { 
            if (response.status === 401) {
                refreshAccessToken();
    
                return search(input, query);
            }
            return response.json();
        })
        .then(data => {
            div = input.closest(".search-container");
            input.style.zIndex = "10";
            html = objt(data);
            let result = document.getElementsByClassName("result")[0];
            if (!result) {
                result = document.createElement('div');
                result.className = "result";
                div.append(result);
            }

            backs = document.getElementsByClassName("res-back")[0];
            if (!backs) {
                backs = document.createElement('div');
                backs.className = "res-back";
                document.body.append(backs);
            }
            
            result.innerHTML = html;
            let images = document.querySelectorAll('.search-file-info');
            MoreImage(images);

            backs.addEventListener("click", function() {
                input.style.zIndex = "auto";
                result.remove();
                backs.remove();
            });
        })
        .catch(error => {
            console.log(error);
        });

}

var input = document.getElementsByName('search')[0];
input.addEventListener("input", function() {
    var query = input.value.trim();
    if (query != '') {
        search(this, query);
    } else {
        let result = document.getElementsByClassName("result")[0];
        if (result) {
            this.style.zIndex = "0";
            result.remove();
            document.getElementsByClassName("res-back")[0].remove();
        }
    }  
})

input.addEventListener("click", function() {
    var query = input.value.trim();
    if (query != '') {
        search(this, query);
    } else {
        let result = document.getElementsByClassName("result")[0];
        if (result) {
            this.style.zIndex = "auto";
            result.remove();
            document.getElementsByClassName("res-back")[0].remove();
        }
    } 
})

document.querySelector(".search").addEventListener("search", function(event) {
    event.preventDefault();
})
document.querySelector(".search").addEventListener("submit", function(event) {
    event.preventDefault();
})


