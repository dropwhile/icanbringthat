// Start registration when the user clicks a button
async function registerPasskey(csrfToken) {
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

    const { value: keyname } = await Swal.fire({
        title: 'Enter a name for this key',
        input: 'text',
        inputLabel: 'Key Name',
        inputValue: inputValue,
        showCancelButton: false,
        inputValidator: (value) => {
            if (!value) {
                return 'You need to write something!'
            }
        }
    });

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