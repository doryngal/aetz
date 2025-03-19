const passwordInput = document.getElementById('password');
const showPasswordCheckbox = document.getElementById('showPassword');

showPasswordCheckbox.addEventListener('change', function() {
    // Toggle the type of the password input
    if (showPasswordCheckbox.checked) {
        passwordInput.type = 'text';
    } else {
        passwordInput.type = 'password';
    }
});