function displayAuth(id) {
    let auth = document.getElementById(id)

    if (auth.style.display === "none") {
        document.getElementById('github-logo').style.display = 'none'
        document.getElementById('google-logo').style.display = 'none'
        auth.style.display = "block"
    } else check()

}


function check() {
    let form = document.getElementById("auth_form")

    let mail = document.getElementById("mail-check")
    mail.style.display = "none";

    let pass = document.getElementById("pass-check")
    pass.style.display = "none";

    switch (true) {
        // Mail
        case !/^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(form.elements["mail"].value):

            if (form.elements["mail"].value === "") return;
            mail.innerHTML = "Wrong Email Format!";
            mail.style.display = "block";
            return;
        // Password
        case !/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,128}$/.test(form.elements["password"].value):
            if (form.elements["password"].value === "") return;
            pass.innerHTML = "Wrong Password Format!";
            pass.style.display = "block";
            return;

        default:
            form.submit()
            break;
    }
}