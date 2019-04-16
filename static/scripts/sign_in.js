/*
    This function checks the form input on sign in page
*/
function validateSubmit() {
    // Checking username
    var username_obj = document.getElementById("username");

    // Checking if username is entered
    if (username_obj.value == "") {
        username_obj.placeholder = "Username cannot be blank!";
        username_obj.style.boxShadow = "0 0 5px #ff0000";
        username_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        username_obj.style.boxShadow = "0 0 0px #ff0000";
        username_obj.style.border = "1px solid #737373";
    }

    // Checking the username format. It must contain only English letters, numbers and underscores
    re = /^\w+$/;
    if (!re.test(username_obj.value)) {
        username_obj.value = "";
        username_obj.placeholder = "Username must contain only letters, numbers and underscores!";
        username_obj.style.boxShadow = "0 0 5px #ff0000";
        username_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        username_obj.style.boxShadow = "0 0 0px #ff0000";
        username_obj.style.border = "1px solid #737373";
    }

    // Checking password
    var passwd_obj = document.getElementById("password");

    // Password minimal length is 6 characters
    if (passwd_obj.value.length < 6) {
        passwd_obj.value = "";
        passwd_obj.placeholder = "Password must contain at least six characters!";
        passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        passwd_obj.style.border = "1px solid #737373";
    }

    // Password must not match the username
    if (passwd_obj.value == username_obj.value) {
        passwd_obj.value = "";
        passwd_obj.placeholder = "Password must be different from Username!";
        passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {      
        passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        passwd_obj.style.border = "1px solid #737373";
    }

    // Password must contain at least one number
    re = /[0-9]/;
    if (!re.test(passwd_obj.value)) {
        passwd_obj.value = "";
        passwd_obj.placeholder = "Password must contain at least one number (0-9)!";
        passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        passwd_obj.style.border = "1px solid #737373";
    }

    // Password must contain at least one lowercase letter
    re = /[a-z]/;
    if (!re.test(passwd_obj.value)) {
        passwd_obj.value = "";
        passwd_obj.placeholder = "Password must contain at least one lowercase letter (a-z)!";
        passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        passwd_obj.style.border = "1px solid #737373";
    }

    // Password must contain at least one uppercase letter
    re = /[A-Z]/;
    if (!re.test(passwd_obj.value)) {
        passwd_obj.value = "";
        passwd_obj.placeholder = "Password must contain at least one uppercase letter (A-Z)!";
        passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        passwd_obj.style.border = "1px solid #737373";
    }

    return true;
}