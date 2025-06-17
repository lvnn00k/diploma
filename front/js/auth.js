document.querySelector(".auth").addEventListener("submit", function(event) {
    event.preventDefault();

    active = document.getElementsByTagName("button")[0].getAttribute("data-active");
    if (active === "false") {
        return;
    }

    const formData = new FormData(this); 

    const request = {
        login: formData.get('login'),
        password: formData.get('password')
    }

    fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
    })
    .then(response =>{
        if (response.ok) {
            window.location.href = "/1"; 
            return;        
        } else {
            document.querySelector("small").textContent="Неверный логин или пароль";
        }
    })
    .catch(error => {
        console.log(error);
    })
    
});

function check() {
    login = document.getElementsByName("login")[0].value.trim();
    button = document.getElementsByTagName("button")[0];
    if (login === '') {
        button.setAttribute('data-active', 'false');
    } else {
        password = document.getElementsByName("password")[0].value.trim();
        if (password === '') {
            button.setAttribute('data-active', 'false');
        } else {
            button.setAttribute('data-active', 'true');
        }
    }
}


document.getElementsByName("login")[0].addEventListener("input", function() {
    check();
});

document.getElementsByName("password")[0].addEventListener("input", function() {
    check();
});




