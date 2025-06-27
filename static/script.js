document.addEventListener('DOMContentLoaded', function() {
    // Fun loading messages
    const funnyLoadingMessages = [
        "Brewing coffee for the server hamsters...",
        "Warming up the login processors...",
        "Dusting off your account details...",
        "Searching for your data in the cloud...",
        "Convincing the AI not to take over the world...",
        "Untangling the internet tubes...",
        "Herding bytes into their proper places..."
    ];

    // Random witty messages - these change on reload
    const loginWittyMessages = [
        "Passwords are like underwear... make them mysterious but don't share them!",
        "Login attempts are like first dates - sometimes they need a second try.",
        "Your password is like a toothbrush. Choose a good one, change it regularly, and don't share it!",
        "Hackers don't break in through the front door - unless you leave it unlocked.",
        "If your password is 'password', maybe reconsider your life choices."
    ];

    const registerWittyMessages = [
        "Don't worry, we won't judge your password strength (but our algorithm might)",
        "Creating an account is like planting a tree - do it once, enjoy it for years.",
        "We promise to keep your data safer than your secrets in a group chat",
        "Your email is safe with us (we won't sell it... for less than a million dollars)",
        "Your account will be more secure than that sandwich in the office fridge"
    ];

    // Set random witty messages
    document.querySelector('#loginForm .witty-message').textContent =
        loginWittyMessages[Math.floor(Math.random() * loginWittyMessages.length)];

    document.querySelector('#registerForm .witty-message').textContent =
        registerWittyMessages[Math.floor(Math.random() * registerWittyMessages.length)];

    // Tab switching functionality
    const tabBtns = document.querySelectorAll('.tab-btn');
    const forms = document.querySelectorAll('.form-container form');

    tabBtns.forEach(button => {
        button.addEventListener('click', function() {
            // Update active tab button
            tabBtns.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');

            // Show corresponding form
            const formId = this.getAttribute('data-form') + 'Form';
            forms.forEach(form => {
                form.classList.remove('active');
                if (form.id === formId) {
                    form.classList.add('active');
                }
            });

            // Clear previous form messages when switching tabs
            document.querySelectorAll('.form-message').forEach(msg => {
                msg.textContent = '';
                msg.className = 'form-message';
            });
        });
    });

    // Fix placeholder issue with material design input
    document.querySelectorAll('input').forEach(input => {
        input.setAttribute('placeholder', ' '); // Empty space for placeholder
    });

    // Registration form submission
    document.getElementById('registerForm').addEventListener('submit', async function(event) {
        event.preventDefault();

        const messageElement = document.getElementById('registerMessage');
        messageElement.textContent = funnyLoadingMessages[Math.floor(Math.random() * funnyLoadingMessages.length)];
        messageElement.className = 'form-message';

        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        if (!name || !email || !password) {
            messageElement.textContent = 'All fields are required (yes, every single one)';
            messageElement.className = 'form-message error';
            return;
        }

        try {
            // Simulate loading time for the joke message to be seen
            await new Promise(resolve => setTimeout(resolve, 1500));

            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, email, password })
            });

            const result = await response.json();
            console.log('Register response:', result);

            if (result.success) {
                // Clear form
                this.reset();

                // Show success message
                messageElement.textContent = result.message || "Account created! You're officially part of the club now.";
                messageElement.className = 'form-message success';

                // Switch to login tab after successful registration
                setTimeout(() => {
                    document.querySelector('[data-form="login"]').click();
                    showSuccessBanner('Account created! Time to show off those login skills!');
                }, 1500);
            } else {
                // Show error message with a fun twist
                const errorMessages = {
                    "Email already exists": "This email is already taken. Are you trying to clone yourself?",
                    "Invalid email format": "That doesn't look like an email. We need something with an @ symbol (and preferably a dot).",
                    "Password too weak": "That password is weaker than a paper umbrella. Try something stronger!"
                };

                messageElement.textContent = errorMessages[result.message] || result.message || "Something went wrong. Technology, am I right?";
                messageElement.className = 'form-message error';
            }
        } catch (error) {
            console.error('Registration error:', error);
            messageElement.textContent = 'Our servers are having an existential crisis. Please try again.';
            messageElement.className = 'form-message error';
        }
    });

    // Login form submission
    document.getElementById('loginForm').addEventListener('submit', async function(event) {
        event.preventDefault();

        const messageElement = document.getElementById('loginMessage');
        messageElement.textContent = funnyLoadingMessages[Math.floor(Math.random() * funnyLoadingMessages.length)];
        messageElement.className = 'form-message';

        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;

        if (!email || !password) {
            messageElement.textContent = 'Email and password are both required (shocking, we know)';
            messageElement.className = 'form-message error';
            return;
        }

        try {
            // Simulate loading time for the joke message to be seen
            await new Promise(resolve => setTimeout(resolve, 1500));

            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            const result = await response.json();
            console.log('Login response:', result);

            if (result.success) {
                // Clear form
                this.reset();

                // Show success message
                messageElement.textContent = "Login successful! Hold on to your hat!";
                messageElement.className = 'form-message success';

                // Display welcome banner
                const welcomeMessages = [
                    `Welcome back, ${result.name || 'User'}! We've missed you terribly.`,
                    `${result.name || 'User'} has entered the chat!`,
                    `The prodigal ${result.name || 'User'} returns!`,
                    `Welcome back! Your inbox has been crying without you, ${result.name || 'User'}.`
                ];

                showSuccessBanner(welcomeMessages[Math.floor(Math.random() * welcomeMessages.length)]);

                // Store user info in localStorage
                localStorage.setItem('userLoggedIn', 'true');
                localStorage.setItem('userName', result.name || 'User');

                // In a real app, you might redirect to a dashboard or homepage after a delay
                // setTimeout(() => {
                //     window.location.href = '/dashboard';
                // }, 2000);
            } else {
                // Show error message with a fun twist
                const errorMessages = {
                    "Invalid credentials": "That email/password combo is like oil and water - they don't mix.",
                    "User not found": "We looked everywhere, but couldn't find you. Are you sure you exist?",
                    "Account locked": "Your account is on time-out. It needs to think about what it did."
                };

                messageElement.textContent = errorMessages[result.message] || result.message || "Login failed. The keys to the kingdom remain elusive.";
                messageElement.className = 'form-message error';
            }
        } catch (error) {
            console.error('Login error:', error);
            messageElement.textContent = 'Our servers are currently taking a coffee break. Please try again.';
            messageElement.className = 'form-message error';
        }
    });

    // Success banner display function
    function showSuccessBanner(message) {
        const banner = document.getElementById('successBanner');
        const bannerMessage = document.getElementById('successMessage');

        bannerMessage.textContent = message;
        banner.classList.remove('hidden');

        setTimeout(() => {
            banner.classList.add('hidden');
        }, 5000);
    }

    // Animate logo letters on hover
    const logoLetters = document.querySelectorAll('.letter');
    logoLetters.forEach(letter => {
        letter.addEventListener('mouseover', function() {
            this.style.transform = 'translateY(-10px) rotate(' + (Math.random() * 20 - 10) + 'deg)';
        });

        letter.addEventListener('mouseout', function() {
            this.style.transform = 'translateY(0) rotate(0)';
        });
    });

    // Check if user is already logged in (for demonstration purposes)
    if (localStorage.getItem('userLoggedIn') === 'true') {
        const userName = localStorage.getItem('userName') || 'User';
        const welcomeBackMessages = [
            `${userName} is back! Hide the good snacks!`,
            `Look who it is! ${userName} has returned!`,
            `${userName} is back! We already have your usual ready.`,
            `Welcome back, ${userName}! We promise we didn't touch anything.`
        ];

        showSuccessBanner(welcomeBackMessages[Math.floor(Math.random() * welcomeBackMessages.length)]);
    }
});