function  ShowHidePass(id,box) {
    let x = document.getElementById(box);
    if (x.checked) {
        document.getElementById(id).type = "text";
    } else {
        document.getElementById(id).type = "password";
    }
}