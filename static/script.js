function editField(event, fieldId, userID) {
  event.preventDefault();
  // Get the span element containing the data
  const spanElement = document.getElementById(fieldId);
  const currentValue = spanElement.innerText;

  // Create an input element
  const inputElement = document.createElement("input");
  inputElement.type = "text";
  inputElement.value = currentValue;

  // Replace the span with the input element
  spanElement.replaceWith(inputElement);

  // Focus on the input field
  inputElement.focus();

  // Flag to track if the input field has been interacted with
  let hasInteracted = false;

  // Handle saving the new value when the user presses Enter
  inputElement.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
      saveField(fieldId, inputElement, userID);
    }
  });

  
  // Set the flag to true when the user interacts with the input field
  inputElement.addEventListener("input", () => {
    hasInteracted = true;
  });
}

function saveField(fieldId, inputElement, userID) {
  const newValue = inputElement.value;

  // Send the updated value to the server
  fetch(`/update/${userID}/${fieldId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ value: newValue }),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.success) {
        // Replace the input element with the updated span
        const spanElement = document.createElement("span");
        spanElement.id = fieldId;
        spanElement.innerText = newValue;
        inputElement.replaceWith(spanElement);
      } else {
        alert("Failed to update the field.");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      alert("An error occurred while updating the field.");
    });
}

// Override fetch so it sends token with it
const originalFetch = window.fetch; // save original fetch so we can still use it as original fetch
window.fetch = async function(url, options={}){
  const token = localStorage.getItem('token');
 
  if (!options.headers) {
      options.headers = new Headers();
  } else if (!(options.headers instanceof Headers)) {
      options.headers = new Headers(options.headers);
  }

  if (token) {
    console.log("Token inside of :", token)
    options.headers.set('Authorization', `Bearer ${token}`);
  }
  
  return originalFetch(url, options);
};


document.addEventListener("DOMContentLoaded", function () {
    const authLink = document.getElementById("auth-link");
    const token = localStorage.getItem("token");

    if (token) {
        // If logged in, show "Logout" link
        authLink.innerHTML = '<a href="#" id="logout-btn">Logout</a>';
        // Event listener for logout
        document.getElementById("logout-btn").addEventListener("click", function (event) {
            event.preventDefault();
            localStorage.removeItem("token");  // Clear token
            location.reload();  // Refresh to reflect changes
        });
    }
});

