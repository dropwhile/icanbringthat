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
    const resp = await fetch('/webauthn/register');
    const respJ = await resp.json();

    let attResp;
    try {
        // Pass the options to the authenticator and wait for a response
        attResp = await SimpleWebAuthnBrowser.startRegistration(respJ.publicKey);
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
                text: 'Error: '+error
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

async function authPasskey(csrfToken) {
    if (!SimpleWebAuthnBrowser.browserSupportsWebAuthn()) {
        Swal.fire({
            icon: 'error',
            title: 'This browser does not support passkeys',
            showCancelButton: false,
        }).then((result) => {
            location.href = '/login';
        });
        return;
    }

    // GET authentication options from the endpoint that calls
    // @simplewebauthn/server -> generateAuthenticationOptions()
    const resp = await fetch('/webauthn/login');
    const respJ = await resp.json();

    let asseResp;
    try {
      // Pass the options to the authenticator and wait for a response
      asseResp = await SimpleWebAuthnBrowser.startAuthentication(respJ.publicKey);
    } catch (error) {
      Swal.fire({
          icon: 'error', title: 'Oops...',
          text: 'Error: '+error
      }).then((result) => {
          location.href = '/login';
      });
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
      body: JSON.stringify(asseResp),
    });

    // Wait for the results of verification
    const verificationJSON = await verificationResp.json();

    // Show UI appropriate for the `verified` status
    if (verificationJSON && verificationJSON.verified) {
        location.href = "/dashboard";
    } else {
      Swal.fire({
          icon: 'error', title: 'Oops...',
          text: 'Error: '+error
      }).then((result) => {
          location.href = '/login';
      });
      return;
    }
}

window.authPasskey = authPasskey