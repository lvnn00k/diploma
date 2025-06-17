
document.addEventListener('DOMContentLoaded', function() {
    const container = document.querySelector('body');
    if (container) {

        container.addEventListener("click", function(event) {
            if (event.target && event.target.matches('[data-form="file"]')) {
                b = document.querySelector('[data-form="file"]');

                b.setAttribute('data-active', 'true');
                let div = document.createElement('div');
                div.className = "div-form";
                div.innerHTML =
                    `<div id="back">
                        <div class="form-container">
                            <form class="form file">
                                <div >
                                    <h3>
                                        Добавить файл
                                    </h3>
                                </div>
                                <div >
                                    <div >
                                        <input type="file" accept="image/*" name="file" data-input="file-file">
                                    </div>
                                </div>
                                <div >
                                    <div>
                                        <input type="text" name="name" placeholder="Название" data-input="file-name">
                                        <small></small>
                                    </div>
                                </div>
                                <div class="foot">
                                    <a href="#" id="cancelBtn">Отмена</a>
                                    <button type="submit" class="btn" data-active="true" data-btn="new">Добавить файл</button>
                                </div>
                            </form>
                        </div>
                    </div>`;
                if (document.getElementById("back")) {
                    document.getElementById("back").remove();
                }
                document.body.append(div);

                fileInput = document.getElementsByName("file")[0];
                nameInput = document.getElementsByName("name")[0];
                fileInput.addEventListener("change", function() {
                    if (fileInput.files.length > 0) {
                        const fileName = fileInput.files[0].name;
                        nameInput.value = fileName.split('.').slice(0, -1).join('.');
                        nameInput.select();
                    }
                })

                document.getElementById("cancelBtn").addEventListener("click", function() {
                    div.remove();
                    document.querySelector('[data-form="file"]').setAttribute('data-active', 'false');
                });
                document.getElementById("back").addEventListener("click", function(event) {
                    if (event.target === document.getElementById("back")) {
                        div.remove();
                        document.querySelector('[data-form="file"]').setAttribute('data-active', 'false');
                    }
                });
            }

            if (event.target && event.target.matches('[data-form="folder"]')) {
                b = document.querySelector('[data-form="folder"]')

                b.setAttribute('data-active', 'true');
                let div = document.createElement('div');
                div.className = "div-form";
                div.innerHTML =
                    `<div id="back">
                        <div class="form-container">
                            <form class="form folder">
                                <div >
                                    <h3>
                                        Новая папка
                                    </h3>
                                </div>
                                <div >
                                    <div>
                                        <input type="text" name="name" placeholder="Название" data-input="folder-name">
                                        <small></small>
                                    </div>
                                </div>
                                <div class="foot">
                                    <a href="#" id="cancelBtn">Отмена</a>
                                    <button type="submit" class="btn" data-active="true" data-btn="new">Создать папку</button>
                                </div>
                            </form>
                        </div>
                    </div>
                    `;
                if (document.getElementById("back")) {
                    document.getElementById("back").remove();
                }
                document.body.append(div);

                document.getElementById("cancelBtn").addEventListener("click", function() {
                    div.remove();
                    document.querySelector('[data-form="folder"]').setAttribute('data-active', 'false');
                });
                document.getElementById("back").addEventListener("click", function(event) {
                    if (event.target === document.getElementById("back")) {
                        div.remove();
                        document.querySelector('[data-form="folder"]').setAttribute('data-active', 'false');
                    }
                });
            }

        })

    }
})


function MoreImage(images) {
    images.forEach(image => {
        image.addEventListener('click', function() {
            link = this.getAttribute('data-link')+"/"+this.querySelector(".name").textContent;
            let div = document.createElement('div');
            div.innerHTML = '<div id="back"><div class="image-container"><img id="panImage" src="/files/'+link+'"><div id="closeBtn">&#10006</div></div></div>';
            document.body.append(div);

            const imageContainer = document.querySelector('.image-container');
            const panImage = document.getElementById('panImage');

            let isDragging = false;
            let startX, startY, translateX = 0, translateY = 0;
            let scale = 1;

            const startDragging = (e) => {
                isDragging = true;
                startX = e.pageX;
                startY = e.pageY;
                imageContainer.style.cursor = 'grabbing';
            };

            const stopDragging = () => {
                isDragging = false;
                imageContainer.style.cursor = 'grab';
            };

            const drag = (e) => {
                if (!isDragging) return;
                e.preventDefault();
                const x = e.pageX;
                const y = e.pageY;
                const walkX = (x - startX);
                const walkY = (y - startY);
                translateX += walkX;
                translateY += walkY;
                panImage.style.transform = `scale(${scale}) translate(${translateX}px, ${translateY}px)`;
                startX = x;
                startY = y;
            };

            imageContainer.addEventListener('mousedown', startDragging);
            imageContainer.addEventListener('mouseup', stopDragging);
            imageContainer.addEventListener('mouseleave', stopDragging);
            imageContainer.addEventListener('mousemove', drag);

            imageContainer.addEventListener('touchstart', startDragging);
            imageContainer.addEventListener('touchend', stopDragging);
            imageContainer.addEventListener('touchcancel', stopDragging);
            imageContainer.addEventListener('touchmove', drag);

            imageContainer.addEventListener('wheel', function(e) {
                e.preventDefault();
                const delta = e.deltaY;
                const zoomSpeed = 0.1;

                if (delta < 0) {
                    scale += zoomSpeed;
                } else {
                    scale -= zoomSpeed;
                }

                scale = Math.max(1, scale);
                panImage.style.transform = `scale(${scale}) translate(${translateX}px, ${translateY}px)`;
            });

            imageContainer.addEventListener('gesturechange', function(e) {
                e.preventDefault();
                const delta = e.deltaY;
                const zoomSpeed = 0.1;

                if (delta < 0) {
                    scale += zoomSpeed;
                } else {
                    scale -= zoomSpeed;
                }

                scale = Math.max(1, scale);
                panImage.style.transform = `scale(${scale}) translate(${translateX}px, ${translateY}px)`;
            });

            closeBtn.addEventListener("click", function(event) {
                if (event.target === closeBtn) {
                    div.remove();
                }
            });
            back.addEventListener("click", function(event) {
                const rect = panImage.getBoundingClientRect();
                const x = event.clientX;
                const y = event.clientY;

                if (!(x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom)) {
                    div.remove();
                } 
            });

        });
    });
}

function showImage() {
    const images = document.querySelectorAll('.file-info');
    MoreImage(images);
}

function MoreInfo() {
    const infos = document.querySelectorAll('.more-info');
    infos.forEach(info => {
        
        info.addEventListener('click', function() {
            let hidden = document.getElementById("hidden");

            if (!hidden) {
                hidden = document.createElement('div');
                hidden.setAttribute("id", "hidden");
                document.body.append(hidden);
            }

            let div = document.getElementsByClassName("more-info-block")[0];

            if (!div) {
                line = this.closest("tr");
                line.style.background = "rgb(206, 206, 206)";
                div = document.createElement('div');
                div.className = "more-info-block";
                if (line.getAttribute("data-type") === "file") {
                    fileName = line.getElementsByClassName("name")[0].textContent;
                    fileLink = line.getElementsByClassName("file-info")[0].getAttribute("data-link");
                    div.innerHTML = `<div><a href="/files/`+fileLink+`/`+fileName+`" download="`+fileName+`">Скачать</a></div><div class="delete" style="color:red;">Удалить</div>`;
                } else {
                    div.innerHTML = `<div class="delete" style="color:red;">Удалить</div>`;
                }
                this.append(div)
            }

            
            hidden.addEventListener('click', function(event) {
                if (event.target === hidden) {
                    div.remove();
                    hidden.remove();
                    line.style.background = "none";
                } 
            });

        });
    });
}