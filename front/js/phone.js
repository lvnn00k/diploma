function phoneBtn() {
    let phonebuttons = document.createElement("div");
    phonebuttons.className = "phonebutton";
    phonebuttons.innerHTML = "&#43";
    document.body.append(phonebuttons);

    phonebuttons.addEventListener("click", function() {
        let btnblock = document.createElement("div");
        btnblock.className = "btn-block";
        btnblock.innerHTML = 
            `<div id="back">
                <div class="form-container">
                    <div class="addforms">
                        <div class="addbtns">
                            <button class="btn" data-form="folder" data-active="false">Создать папку</button>
                            <button class="btn" data-form="file" data-active="false">Загрузить файл</button>
                        </div>
                    </div>
                </div>
            </div>
            `;
        document.body.append(btnblock);
        back.addEventListener("click", function(event) {
            if (event.target === back) {
                btnblock.remove();
            }
        });
        

    })
}

window.addEventListener('resize', () => {
    if (window.innerWidth <= 415) {
        phoneBtn();
    } else {
        if (document.getElementsByClassName("btn-block")[0]) {
            document.getElementsByClassName("btn-block")[0].remove();
        }
        if (document.getElementsByClassName("phonebutton")[0]) {
            document.getElementsByClassName("phonebutton")[0].remove();
        }
    }
});

if (window.innerWidth <= 415) {
    phoneBtn();
}