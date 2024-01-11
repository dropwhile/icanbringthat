// Start registration when the user clicks a button
async function registerPasskey(csrfToken) {
    const { value: keyname } = await Swal.fire({
        title: 'Enter a name for this key',
        input: 'text',
        inputLabel: 'Key Name',
        showCancelButton: false,
        inputValidator: (value) => {
            if (!value) {
                return 'You need to write something!'
            }
        }
    });

    if (!keyname) {
        return;
    }

    // GET registration options from the endpoint that calls
    // @simplewebauthn/server -> generateRegistrationOptions()
    const registerResp = await fetch('/webauthn/register');
    const registerJSON = await registerResp.json();

    let attResp;
    try {
        // Pass the options to the authenticator and wait for a response
        attResp = await SimpleWebAuthnBrowser.startRegistration(registerJSON.publicKey);
    } catch (error) {
        // Some basic error handling
        if (error.name === 'InvalidStateError') {
            Swal.fire({
                icon: 'error', title: 'Oops...',
                text: 'Error: Authenticator was probably already registered by user'
            });
        } else {
            Swal.fire({
                icon: 'error', title: 'Oops...',
                text: 'Error: ' + error
            });
        }
        return;
    }

    // POST the response to the endpoint that calls
    // @simplewebauthn/server -> verifyRegistrationResponse()
    const verificationResp = await fetch(
        '/webauthn/register?' + new URLSearchParams({key_name: keyname}).toString(),
        {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify(attResp),
        }
    );

    // Wait for the results of verification
    const verificationJSON = await verificationResp.json();

    // Show UI appropriate for the `verified` status
    if (verificationJSON && verificationJSON.verified) {
        Swal.fire('Passkey is now registered!').then((result) => {
            location.reload();
        });
    } else {
        Swal.fire({
            icon: 'error', title: 'Oops...',
            html: 'Unexpected error response: <br>' + 
                `<pre>${JSON.stringify(verificationJSON)}</pre>`
        });
    }
};

window.registerPasskey = registerPasskey

async function authPasskey(csrfToken, autofill=false) {
    const notyf = new Notyf({
        ripple: false,
        dismissible: true,
        duration: 2500,
        position: {
          x: 'center',
          y: 'top',
        }
    });

    if (!SimpleWebAuthnBrowser.browserSupportsWebAuthn()) {
        notyf.error('This browser does not support passkeys.');
        return;
    }

    // GET authentication options from the endpoint that calls
    // @simplewebauthn/server -> generateAuthenticationOptions()
    const loginResp = await fetch('/webauthn/login');
    const loginJSON = await loginResp.json();

    if (!loginJSON) {
        notyf.error('Oops. Something went wrong.');
        return;
    } else if (loginJSON.error) {
        notyf.error(verificationJSON.error);
        return;
    }

    let startAuthResp;
    try {
      // Pass the options to the authenticator and wait for a response
      startAuthResp = await SimpleWebAuthnBrowser.startAuthentication(
        loginJSON.publicKey, autofill);
    } catch (error) {
        notyf.error('Error: ' + error);
        return;
    }

    // POST the response to the endpoint that calls
    // @simplewebauthn/server -> verifyAuthenticationResponse()
    const verificationResp = await fetch('/webauthn/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify(startAuthResp),
    });

    // Wait for the results of verification
    const verificationJSON = await verificationResp.json();

    // Show UI appropriate for the `verified` status
    if (verificationJSON && verificationJSON.verified) {
        location.href = "/dashboard";
    } else if (verificationJSON && verificationJSON.error) {
        notyf.error(verificationJSON.error);
        return;
    } else {
        notyf.error('Oops. Something went wrong.');
        return;
    }
}

window.authPasskey = authPasskey