function validateSubmit() {
    var username_obj = document.getElementById("username");

    if (username_obj.value == "") {
        username_obj.placeholder = "Username cannot be blank!";
        username_obj.style.boxShadow = "0 0 5px #ff0000";
        username_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        username_obj.style.boxShadow = "0 0 0px #ff0000";
        username_obj.style.border = "1px solid #737373";
    }

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

    var email_obj = document.getElementById("email");
    var at = email_obj.value.indexOf("@");
    var dot = email_obj.value.lastIndexOf(".");

    if (email_obj.value == "") {
        email_obj.placeholder = "Email cannot be blank!"
        email_obj.style.boxShadow = "0 0 5px #ff0000";
        email_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        email_obj.style.boxShadow = "0 0 0px #ff0000";
        email_obj.style.border = "1px solid #737373";
    }

    if (at < 0 || dot < 0 || at > dot) {
        email_obj.value = "";
        email_obj.placeholder = "Wrong email format!"
        email_obj.style.boxShadow = "0 0 5px #ff0000";
        email_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        email_obj.style.boxShadow = "0 0 0px #ff0000";
        email_obj.style.border = "1px solid #737373";
    }

    var passwd_obj = document.getElementById("password");
    var conf_passwd_obj = document.getElementById("confirm_password");

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

    if (passwd_obj.value != conf_passwd_obj.value) {
        conf_passwd_obj.value = "";
        conf_passwd_obj.placeholder = "Passwords do not match!";
        conf_passwd_obj.style.boxShadow = "0 0 5px #ff0000";
        conf_passwd_obj.style.border = "1px solid #ff0000";
        return false;
    } else {           
        conf_passwd_obj.style.boxShadow = "0 0 0px #ff0000";
        conf_passwd_obj.style.border = "1px solid #737373";
    }

    var f_name_obj = document.getElementById("f_name");
    
    if (f_name_obj.value == "") {
        f_name_obj.placeholder = "First name cannot be blank!";
        f_name_obj.style.boxShadow = "0 0 5px #ff0000";
        f_name_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        f_name_obj.style.boxShadow = "0 0 0px #ff0000";
        f_name_obj.style.border = "1px solid #737373";
    }

    var l_name_obj = document.getElementById("l_name");

    if (l_name_obj.value == "") {
        l_name_obj.placeholder = "Last name cannot be blank!";
        l_name_obj.style.boxShadow = "0 0 5px #ff0000";
        l_name_obj.style.border = "1px solid #ff0000";
        return false;
    } else {
        l_name_obj.style.boxShadow = "0 0 0px #ff0000";
        l_name_obj.style.border = "1px solid #737373";
    }

    return true;
}
